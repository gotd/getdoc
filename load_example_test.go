package getdoc_test

import (
	"fmt"

	"github.com/gotd/getdoc"
)

func ExampleLoad() {
	layer := 133
	if !getdoc.LayerExists(121) {
		panic("not exists")
	}
	if !getdoc.LayerExists(133) {
		panic("not exists")
	}
	doc, err := getdoc.Load(layer)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Layer %d, constructors: %d\n", doc.Index.Layer, len(doc.Constructors))
	// Output:
	// Layer 133, constructors: 926
}
