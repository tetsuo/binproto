# binproto

Simple binary messaging protocol with multiplexing support.

[![Go Reference](https://pkg.go.dev/badge/github.com/tetsuo/binproto.svg)](https://pkg.go.dev/github.com/tetsuo/binproto)

## Message structure

```
┌──────────────────────────────────────────────┐
│ length | channel ID × message type │ payload │
└──────────────────────────────────────────────┘
           └─ 60-bits   └─ 4-bits
```

## Benchmarks

```
goos: darwin
goarch: arm64
pkg: github.com/tetsuo/binproto
cpu: Apple M4 Pro
BenchmarkWriter_64B-12          98590449                12.05 ns/op     5309.66 MB/s           0 B/op          0 allocs/op
BenchmarkWriter_4KB-12          14346594                80.00 ns/op     51198.63 MB/s          0 B/op          0 allocs/op
BenchmarkWriter_1MB-12             74955             15931 ns/op        65819.47 MB/s         14 B/op          0 allocs/op
BenchmarkReader_64B-12          70582526                16.24 ns/op     3941.01 MB/s           0 B/op          0 allocs/op
BenchmarkReader_4KB-12          15629764                78.92 ns/op     51903.93 MB/s          0 B/op          0 allocs/op
BenchmarkReader_1MB-12             38365             31480 ns/op        33309.32 MB/s          0 B/op          0 allocs/op
```
