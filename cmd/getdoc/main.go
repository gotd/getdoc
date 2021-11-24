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

	"github.com/gotd/getdoc"
	"github.com/gotd/getdoc/dl"
)

func main() {
	dir := flag.String("dir", filepath.Join(os.TempDir(), "getdoc"), "working directory")
	readonly := flag.Bool("readonly", false, "read-only mode")
	pretty := flag.Bool("pretty", false, "pretty json output")
	flag.Parse()

	client, err := dl.NewClient(dl.Options{
		Path:     filepath.Join(*dir, "cache"),
		Readonly: *readonly,
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	fmt.Println("Extracting")
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
	outFilePath := filepath.Join(*dir, outFileName)
	if err := os.WriteFile(outFilePath, out.Bytes(), 0600); err != nil {
		panic(err)
	}

	fmt.Println("Wrote layer", doc.Index.Layer, "to", outFilePath)
}
