// Package dl implements documentation downloading facility.
package dl

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/cenkalti/backoff/v4"
	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
	"go.uber.org/ratelimit"
	"golang.org/x/xerrors"
)

type Client struct {
	rate     ratelimit.Limiter
	http     HTTPClient
	db       *pebble.DB
	readonly bool
}

type Options struct {
	Path     string
	FS       vfs.FS
	Client   HTTPClient
	Readonly bool
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(opt Options) (*Client, error) {
	db, err := pebble.Open(opt.Path, &pebble.Options{
		FS:               opt.FS,
		ReadOnly:         opt.Readonly,
		ErrorIfNotExists: opt.Readonly,
	})
	if err != nil {
		return nil, err
	}
	if opt.Client == nil {
		opt.Client = http.DefaultClient
	}

	return &Client{
		http:     opt.Client,
		readonly: opt.Readonly,
		rate:     ratelimit.New(10),
		db:       db,
	}, nil
}

var ErrReadOnly = errors.New("write operation in read only mode")

func (c *Client) Close() error {
	return c.db.Close()
}

// NoLayer can be passed as "layer" argument.
const NoLayer = 0

func (c *Client) download(ctx context.Context, layer int, key string) ([]byte, error) {
	if c.readonly {
		return nil, ErrReadOnly
	}

	// `stel_dev_layer=117`
	u := &url.URL{
		Scheme: "https",
		Host:   "core.telegram.org",
		Path:   key,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	res, err := c.http.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to do request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, xerrors.Errorf("unexpected http code %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

// Get fetches documentation by key and layer.
//
// Pass "NoLayer" as layer to use default documentation.
// Examples for "key" value:
//	* "schema" for index of documentation
//	* "constructor/inputMediaGeoLive" for constructor "inputMediaGeoLive"
//	* "type/InputMedia" for class "InputMedia"
//	* "method/messages.addChatUser" for "messages.addChatUser" method
//
// Blank key is invalid.
func (c *Client) Get(ctx context.Context, layer int, key string) ([]byte, error) {
	k := []byte(key)
	if layer != NoLayer {
		k = []byte(fmt.Sprintf("%d:%s", layer, k))
	}

	// Trying to get from cache.
	buf, closer, err := c.db.Get(k)
	if err == nil {
		// Copy buf because slice is not valid after close.
		data := append([]byte(nil), buf...)
		if err := closer.Close(); err != nil {
			return nil, xerrors.Errorf("cache: %w", err)
		}
		return data, nil
	}
	if !errors.Is(err, pebble.ErrNotFound) {
		return nil, xerrors.Errorf("cache: %w", err)
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
		return nil, xerrors.Errorf("failed to fetch: %w", err)
	}

	// Adding value to cache.
	if err := c.db.Set(k, data, &pebble.WriteOptions{}); err != nil {
		return nil, xerrors.Errorf("cache: %w", err)
	}

	return data, nil
}
