# Go benchmarks

A few benchmarks we did to better know go and make some decisions in
our code.

Feel free to comment, disagree, suggest!

## str_byte_cmp_test.go

Find the fastest way to compare a string and a []byte.

### Run it

```bash
go test -bench . str_byte_cmp_test.go
```

### Result

Given:

```go
var (
    s string
    b []byte
)

_ = s == string(b) // case 1
_ = bytes.Equal([]byte(s), b) // case 2
```

case 1 is ~65% faster than case 2

## reverse_map_test.go

In case we need the key corresponding to a given string value of a map,
which is faster between a lookup in a reversed version of the map (*not*
including the reverse map build) and a simple iteration in the map until
the right value is found?

### Run

```bash
go test -bench reverse_map_test.go
```

### Result

Even with only 1 item with a 1-byte-long string value, the map lookup is
faster than a 1 loop iteration...

## map_vs_switch.go

Measure how a map[string]xxx lookup compares to a switch/case.
Works by generating a \_test.go file with different sizes of map and strings.
The generated file contains the actual benchmarks.

This bench was done to decide how to generate the UnmarshalText function in
https://github.com/orus-io/enum-marshaler.

### Run it

```bash
go run map_vs_switch.go
go test -bench . map_vs_switch.go

```

### Result

A switch/case is faster than a map lookup up to 90 (+/-30) items.
