# DEPRECATED

Now the performance is no better than a `map` or `sync.Map` except that the extra memory allocations are better. So without better performance, maybe `map` or `sync.Map` are good built-in alternatives

Because of that, I decided to archive this library and end support.

Thanks so much 😉!

# dictpool

[![Test status](https://github.com/savsgio/dictpool/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/savsgio/dictpool/actions?workflow=test)
[![Go Report Card](https://goreportcard.com/badge/github.com/savsgio/dictpool)](https://goreportcard.com/report/github.com/savsgio/dictpool)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/savsgio/dictpool)

Memory store like `map[string]interface{}` without extra memory allocations.


## Benchmarks:

```
BenchmarkDict-12                           41162             28199 ns/op               0 B/op          0 allocs/op
BenchmarkStdMap-12                        117976              9195 ns/op           10506 B/op         11 allocs/op
BenchmarkSyncMap-12                        73750             15524 ns/op            3200 B/op        200 allocs/op
```

_Benchmark with Go 1.19_

## Example:

```go
d := dictpool.AcquireDict()
// d.BinarySearch = true  // Useful on big heaps

key := "foo"

d.Set(key, "Hello DictPool")

if d.Has(key){
    fmt.Println(d.Get(key))  // Output: Hello DictPool
}

d.Del(key)

dictpool.ReleaseDict(d)
```
