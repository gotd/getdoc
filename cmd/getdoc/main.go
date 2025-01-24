// Binary getdoc extracts Telegram documentation to json file.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-faster/sdk/zctx"
	"go.uber.org/zap"

	"github.com/gotd/getdoc"
	"github.com/gotd/getdoc/dl"
)

func main() {
	dir := flag.String("dir", filepath.Join(os.TempDir(), "getdoc"), "working directory")
	outDir := flag.String("out-dir", "", "path to write schema")
	host := flag.String("host", "core.telegram.org", "host")
	outFile := flag.String("out-file", "", "filename of schema")
	readonly := flag.Bool("readonly", false, "read-only mode")
	pretty := flag.Bool("pretty", false, "pretty json output")
	flag.Parse()

	lg, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	ctx := zctx.Base(context.Background(), lg)

	cacheDir := filepath.Join(*dir, "cache")
	lg.Info("Cache dir", zap.String("path", cacheDir))

	client, err := dl.NewClient(dl.Options{
		Path:     cacheDir,
		Host:     *host,
		Readonly: *readonly,
	})
	if err != nil {
		panic(err)
	}

	doc, err := getdoc.Extract(ctx, client)
	if err != nil {
		panic(err)
	}

	out := new(bytes.Buffer)
	enc := json.NewEncoder(out)
	if *pretty {
		enc.SetIndent("", "  ")
	}
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}

	outFileName := fmt.Sprintf("%d.json", doc.Index.Layer)
	if out := *outFile; out != "" {
		outFileName = out
	}

	outFilePath := filepath.Join(*dir, outFileName)
	if out := *outDir; out != "" {
		if err := os.MkdirAll(out, 0o600); err != nil {
			panic(err)
		}
		outFilePath = filepath.Join(out, outFileName)
	}

	if err := os.WriteFile(outFilePath, out.Bytes(), 0o600); err != nil {
		panic(err)
	}

	fmt.Println("Wrote layer", doc.Index.Layer, "to", outFilePath)
}
