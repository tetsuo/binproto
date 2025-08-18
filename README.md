# binproto

binproto provides binary message framing using a length-prefixed format and supports multiplexed streams.

## Install

```sh
go get github.com/tetsuo/binproto
```

## Usage

To start exchanging messages with binproto, simply wrap a `net.Conn` with a `binproto.Conn` instance.

Then, once connected, use `ReadMessage()` to read incoming messages and `Send()` to send new ones over the network.

#### Server

```go
package main

import (
    "fmt"
    "log"
    "net"

    "github.com/tetsuo/binproto"
)

func main() {
    s := &server{}
    if err := s.serve("tcp", ":4242"); err != nil {
        log.Fatal(err)
    }
}

type server struct {
    listener net.Listener
}

// handle manages an individual client connection.
// It reads each message, prints its content, and responds with a simple reply.
func (s *server) handle(conn net.Conn) {
    defer conn.Close()
    c := binproto.NewConn(conn)

    for {
        msg, err := c.ReadMessage()
        if err != nil {
            fmt.Printf("error: %v", err)
            return
        }
        fmt.Printf("%d %d %s\n", msg.ID, msg.Channel, msg.Data)

        _, err = c.Send(binproto.NewMessage(112, 5, []byte("hello")))
        if err != nil {
            log.Fatal(err)
        }
    }
}

// serve starts listening for incoming TCP connections.
// For each new connection, it launches handle in a separate goroutine.
func (s *server) serve(network, address string) error {
    l, err := net.Listen(network, address)
    if err != nil {
        return err
    }
    s.listener = l

    for {
        if s.listener == nil {
            break
        }
        c, err := l.Accept()
        if err != nil {
            continue
        }
        go s.handle(c)
    }
    return nil
}

// close cleanly shuts down the server listener.
func (s *server) close() error {
    err := s.listener.Close()
    s.listener = nil
    return err
}
```

#### Client

```go
package main

import (
    "fmt"
    "log"

    "github.com/tetsuo/binproto"
)

func main() {
    // Connect to the server on localhost:4242 using binproto
    c, err := binproto.Dial("tcp", ":4242")
    if err != nil {
        log.Fatal(err)
    }

    // Start a goroutine to read and print incoming messages from the server
    go func() {
        for {
            msg, err := c.ReadMessage()
            if err != nil {
                log.Fatal(err)
                return
            }

            fmt.Printf("%d %d %s\n", msg.ID, msg.Channel, msg.Data)
        }
    }()

    // Send a message to the server: ID=42, Channel=3, Data="hi"
    _, err = c.Send(binproto.NewMessage(42, 3, []byte("hi")))
    if err != nil {
        log.Fatal(err)
    }

    select {} // Keeps main from exiting immediately
}
```

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

## Encryption

binproto itself doesn't provide encryption. However, encryption can be added by using Go libraries that act as drop-in replacements for `net`, for example by wrapping your connections with an implementation of the [NOISE protocol](http://www.noiseprotocol.org/).

## Configure buffer size

binproto operates with a default internal buffer size of 4096 bytes, meaning data is processed as long as it meets or exceeds this buffer size, which is an effective default for many applications. This value can be adjusted to better suit protocols that use larger or smaller data chunks, optimizing performance as needed.

> 📄 **For more details, see the API documentation at [pkg.go.dev](https://pkg.go.dev/github.com/tetsuo/binproto).**
