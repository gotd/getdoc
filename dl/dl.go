// Package dl implements documentation downloading facility.
package dl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-faster/sdk/zctx"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"

	"github.com/gotd/getdoc/dl/internal/cache"
)

type Client struct {
	rate     ratelimit.Limiter
	http     HTTPClient
	cache    cache.Cacher
	host     string
	readonly bool
}

type Options struct {
	Client   HTTPClient
	Host     string
	Path     string
	Readonly bool
	FromZip  bool
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(opt Options) (*Client, error) {
	if opt.Client == nil {
		opt.Client = http.DefaultClient
	}
	if opt.Host == "" {
		opt.Host = "core.telegram.org"
	}

	c := &Client{
		http:     opt.Client,
		readonly: opt.Readonly,
		host:     opt.Host,
		rate:     ratelimit.New(1000),
	}

	if opt.FromZip {
		zipCache, err := cache.NewFromZip(opt.Path)
		if err != nil {
			return nil, err
		}

		c.cache = zipCache
	} else {
		c.cache = cache.NewFromDirectory(opt.Path)
	}

	return c, nil
}

var ErrReadOnly = errors.New("write operation in read only mode")

// NoLayer can be passed as "layer" argument.
const NoLayer = 0

func (c *Client) download(ctx context.Context, layer int, key string) ([]byte, error) {
	if c.readonly {
		return nil, ErrReadOnly
	}

	// `stel_dev_layer=117`
	u := &url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   key,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	if layer != NoLayer {
		// Current layer is determined by cookie value.
		req.AddCookie(&http.Cookie{
			Name:  "stel_dev_layer",
			Value: strconv.Itoa(layer),

			HttpOnly: true,
		})
	}

	c.rate.Take()

	start := time.Now()
	res, err := c.http.Do(req)
	duration := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()
	zctx.From(ctx).Debug("Request",
		zap.String("path", u.Path),
		zap.Duration("duration", duration),
		zap.Int("status", res.StatusCode),
	)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

// Get fetches documentation by key and layer.
//
// Pass "NoLayer" as layer to use default documentation.
// Examples for "key" value:
//   - "schema" for index of documentation
//   - "constructor/inputMediaGeoLive" for constructor "inputMediaGeoLive"
//   - "type/InputMedia" for class "InputMedia"
//   - "method/messages.addChatUser" for "messages.addChatUser" method
//
// Blank key is invalid.
func (c *Client) Get(ctx context.Context, layer int, key string) ([]byte, error) {
	cacheKey := key
	if layer != NoLayer {
		cacheKey = fmt.Sprintf("%d/%s", layer, key)
	}

	// Trying to get from cache.
	buf, err := c.cache.Get(cacheKey)
	if err == nil {
		zctx.From(ctx).Debug("Cache hit",
			zap.String("key", cacheKey),
		)
		return buf, nil
	}
	if !errors.Is(err, cache.ErrNotFound) {
		return nil, fmt.Errorf("cache: %w", err)
	}

	// Downloading with retry backoff.
	var data []byte
	if err := backoff.Retry(func() error {
		out, getErr := c.download(ctx, layer, key)
		for _, permanentErr := range []error{
			context.Canceled,
			context.DeadlineExceeded,
			ErrReadOnly,
		} {
			// Immediately stop on unrecoverable errors.
			if errors.Is(getErr, permanentErr) {
				return backoff.Permanent(getErr)
			}
		}
		if getErr != nil {
			return getErr
		}
		data = out
		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}

	// Adding value to cache.
	if err := c.cache.Set(cacheKey, data); err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return data, nil
}
