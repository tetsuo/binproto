# binproto

`binproto` is a lightweight Go library designed to facilitate low-level binary protocol communication. It enables efficient, multiplexed messaging, making it well-suited for network applications requiring reliable and fast binary data exchange.

## Installation

Install `binproto` by running:

```sh
go get github.com/tetsuo/binproto
```

## Usage

Here’s a basic example of how to use `binproto` for a simple client-server communication setup:

```go
package main

import (
    "log"
    "net"
    "github.com/tetsuo/binproto"
)

func main() {
    // Create a new server instance
    server := &Server{}

    // Start the server
    go func() {
        if err := server.Start(":4242"); err != nil {
            log.Fatal(err)
        }
    }()

    // Connect a client
    client, err := binproto.Dial("tcp", ":4242")
    if err != nil {
        log.Fatal(err)
    }

    // Send a message from the client
    client.Send(binproto.NewMessage(1, 1, []byte("Hello Server")))

    // Close the client
    client.Close()
}
```

### Server Code

To set up a server using `binproto`, create a handler to process incoming messages.

```go
type Server struct {
    listener net.Listener
}

func (s *Server) Start(addr string) error {
    l, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    s.listener = l
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            return err
        }
        go s.handleConnection(conn)
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    c := binproto.NewConn(conn)
    for {
        msg, err := c.ReadMessage()
        if err != nil {
            return
        }
        log.Printf("Received message: %s", msg.Data)
    }
}
```

## License

MIT License. See [LICENSE](LICENSE) for details.
