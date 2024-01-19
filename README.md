# getdoc

Documentation extract utilities for Telegram schema using [goquery](https://github.com/PuerkitoBio/goquery).

Used by [gotd/td](https://github.com/gotd/td) for embedding documentation to generated code.

## Parsed documentation

Parsed documentation for 133 layer is available as [133.json](./_schema/133.json) with [schema](./_schema/schema.json).

## Example
Latest schema is embedded to package, so you can just use it:
```go
doc, err := getdoc.Load(133)
if err != nil {
    panic(err)
}
fmt.Printf("Layer %d, constructors: %d\n", doc.Index.Layer, len(doc.Constructors))
// Output:
// Layer 133, constructors: 926
```

## Reference

Please use [official documentation](https://core.telegram.org/schema), it is for humans,
this package is not.
