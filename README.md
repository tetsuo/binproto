# binproto

Simple binary protocol with multiplexing support.

[![Go Reference](https://pkg.go.dev/badge/github.com/tetsuo/binproto.svg)](https://pkg.go.dev/github.com/tetsuo/binproto)

## Message structure

Every message is encoded with a 64-bit header: a length prefix followed by a channel ID and type.

```
╔──────────────────────────────────────────────╗
│ length | channel ID × channel type │ payload │
╚──────────────────────────────────────────────╝
           └─ 60-bits   └─ 4-bits
```

* **Channel ID (first 60 bits)**: Identifies the specific channel for the message.
* **Channel Type (last 4 bits)**: Specifies the type of data in the message.

## Benchmarks

The benchmarks measure the performance for different payload sizes (64 bytes, 4 KB, and 1 MB).

With the default 4kB buffers, binproto achieves a throughput of approximately 50 GB/s for both writing and reading operations.

```
goos: darwin
goarch: arm64
pkg: github.com/tetsuo/binproto
cpu: Apple M4 Pro
BenchmarkWriter_64B-12          100000000               12.15 ns/op     5266.44 MB/s           0 B/op          0 allocs/op
BenchmarkWriter_4KB-12          14390127                81.01 ns/op     50564.20 MB/s          0 B/op          0 allocs/op
BenchmarkWriter_1MB-12             75183             15905 ns/op        65925.85 MB/s         14 B/op          0 allocs/op
BenchmarkReader_64B-12          70587026                16.55 ns/op     3868.22 MB/s           0 B/op          0 allocs/op
BenchmarkReader_4KB-12          15426055                78.05 ns/op     52477.75 MB/s          0 B/op          0 allocs/op
BenchmarkReader_1MB-12             38109             31472 ns/op        33317.35 MB/s         27 B/op          0 allocs/op
```
