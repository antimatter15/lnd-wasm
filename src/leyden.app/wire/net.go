package main

import (
    // "errors"
    "net"
    "time"
    "fmt"
    "syscall/js"
)

// TODO: this interface and its implementations should ideally be moved
// elsewhere as they are not Tor-specific.

// Net is an interface housing a Dial function and several DNS functions that
// allows us to abstract the implementations of these functions over different
// networks, e.g. clearnet, Tor net, etc.
type Net interface {
    // Dial connects to the address on the named network.
    Dial(network, address string) (net.Conn, error)

    // LookupHost performs DNS resolution on a given host and returns its
    // addresses.
    LookupHost(host string) ([]string, error)

    // LookupSRV tries to resolve an SRV query of the given service,
    // protocol, and domain name.
    LookupSRV(service, proto, name string) (string, []*net.SRV, error)

    // ResolveTCPAddr resolves TCP addresses.
    ResolveTCPAddr(network, address string) (*net.TCPAddr, error)
}

// ClearNet is an implementation of the Net interface that defines behaviour
// for regular network connections.
type ClearNet struct{}


type Conn struct {
    // host string
    id int
    // conn net.Conn

    // noise *Machine

    // readBuf bytes.Buffer
}

var _ net.Conn = (*Conn)(nil)

// Read reads data from the connection.  Read can be made to time out and
// return an Error with Timeout() == true after a fixed time limit; see
// SetDeadline and SetReadDeadline.
//
// Part of the net.Conn interface.
func (c *Conn) Read(b []byte) (n int, err error) {
    fmt.Println("read (start)", b)

    done := make(chan int)

    callback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        fmt.Println("got response")
        // callback.Release() // free up memory from callback
        done <- args[0].Int()
        return nil
    })

    // func printMessage(this js.Value, args []js.Value) interface{} {
    //     message := args[0].String()
    //     fmt.Println(message)
    //     
    //     return nil
    // }
    defer callback.Release()
    
    ta := js.TypedArrayOf(b)
    defer ta.Release()

    js.Global().Get("readFromSocket").Invoke(c.id, ta, callback)

    // wait until we've got our response
    bytesRead := <-done


    fmt.Println("read (done)", b)

//     // In order to reconcile the differences between the record abstraction
//     // of our AEAD connection, and the stream abstraction of TCP, we
//     // maintain an intermediate read buffer. If this buffer becomes
//     // depleted, then we read the next record, and feed it into the
//     // buffer. Otherwise, we read directly from the buffer.
//     if c.readBuf.Len() == 0 {
//         plaintext, err := c.noise.ReadMessage(c.conn)
//         if err != nil {
//             return 0, err
//         }

//         if _, err := c.readBuf.Write(plaintext); err != nil {
//             return 0, err
//         }
//     }

//     return c.readBuf.Read(b)
    return bytesRead, nil
}

// Write writes data to the connection.  Write can be made to time out and
// return an Error with Timeout() == true after a fixed time limit; see
// SetDeadline and SetWriteDeadline.
//
// Part of the net.Conn interface.
func (c *Conn) Write(b []byte) (n int, err error) {

    fmt.Println("writing", b)

    done := make(chan int)

    callback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        fmt.Println("wrote stuff")
        // callback.Release() // free up memory from callback
        done <- 0
        return nil
    })

    // func printMessage(this js.Value, args []js.Value) interface{} {
    //     message := args[0].String()
    //     fmt.Println(message)
    //     
    //     return nil
    // }
    defer callback.Release()
    ta := js.TypedArrayOf(b)
    defer ta.Release()
    js.Global().Get("writeToSocket").Invoke(c.id, ta, callback)

    // wait until we've got our response
    bytesWritten := <-done


//     // If the message doesn't require any chunking, then we can go ahead
//     // with a single write.
//     if len(b) <= math.MaxUint16 {
//         err = c.noise.WriteMessage(b)
//         if err != nil {
//             return 0, err
//         }
//         return c.noise.Flush(c.conn)
//     }

//     // If we need to split the message into fragments, then we'll write
//     // chunks which maximize usage of the available payload.
//     chunkSize := math.MaxUint16

    // bytesToWrite := len(b)
//     for bytesWritten < bytesToWrite {
//         // If we're on the last chunk, then truncate the chunk size as
//         // necessary to avoid an out-of-bounds array memory access.
//         if bytesWritten+chunkSize > len(b) {
//             chunkSize = len(b) - bytesWritten
//         }

//         // Slice off the next chunk to be written based on our running
//         // counter and next chunk size.
//         chunk := b[bytesWritten : bytesWritten+chunkSize]
//         if err := c.noise.WriteMessage(chunk); err != nil {
//             return bytesWritten, err
//         }

//         n, err := c.noise.Flush(c.conn)
//         bytesWritten += n
//         if err != nil {
//             return bytesWritten, err
//         }
//     }

    return bytesWritten, nil
}

// // WriteMessage encrypts and buffers the next message p for the connection. The
// // ciphertext of the message is prepended with an encrypt+auth'd length which
// // must be used as the AD to the AEAD construction when being decrypted by the
// // other side.
// //
// // NOTE: This DOES NOT write the message to the wire, it should be followed by a
// // call to Flush to ensure the message is written.
// func (c *Conn) WriteMessage(b []byte) error {
//     return c.noise.WriteMessage(b)
// }

// // Flush attempts to write a message buffered using WriteMessage to the
// // underlying connection. If no buffered message exists, this will result in a
// // NOP. Otherwise, it will continue to write the remaining bytes, picking up
// // where the byte stream left off in the event of a partial write. The number of
// // bytes returned reflects the number of plaintext bytes in the payload, and
// // does not account for the overhead of the header or MACs.
// //
// // NOTE: It is safe to call this method again iff a timeout error is returned.
// func (c *Conn) Flush() (int, error) {
//     return c.noise.Flush(c.conn)
// }

// Close closes the connection.  Any blocked Read or Write operations will be
// unblocked and return errors.
//
// Part of the net.Conn interface.
func (c *Conn) Close() error {
    // TODO(roasbeef): reset brontide state?
    // return c.conn.Close()
    fmt.Println("closed conn")
    return nil
}

// LocalAddr returns the local network address.
//
// Part of the net.Conn interface.
func (c *Conn) LocalAddr() net.Addr {
    // return c.conn.LocalAddr()
    fmt.Println("getting lcoal addr")
    return nil
}

// RemoteAddr returns the remote network address.
//
// Part of the net.Conn interface.
func (c *Conn) RemoteAddr() net.Addr {
    // return c.conn.RemoteAddr()
    fmt.Println("getting remote addr")
    return nil
}

// SetDeadline sets the read and write deadlines associated with the
// connection. It is equivalent to calling both SetReadDeadline and
// SetWriteDeadline.
//
// Part of the net.Conn interface.
func (c *Conn) SetDeadline(t time.Time) error {
    // return c.conn.SetDeadline(t)
    fmt.Println("set deadline", t)
    return nil
}

// SetReadDeadline sets the deadline for future Read calls.  A zero value for t
// means Read will not time out.
//
// Part of the net.Conn interface.
func (c *Conn) SetReadDeadline(t time.Time) error {
    // return c.conn.SetReadDeadline(t)
    fmt.Println("set read deadline", t)
    return nil
}

// SetWriteDeadline sets the deadline for future Write calls.  Even if write
// times out, it may return n > 0, indicating that some of the data was
// successfully written.  A zero value for t means Write will not time out.
//
// Part of the net.Conn interface.
func (c *Conn) SetWriteDeadline(t time.Time) error {
    // return c.conn.SetWriteDeadline(t)
    fmt.Println("set write deadline", t)
    return nil
}


// Dial on the regular network uses net.Dial
func (r *ClearNet) Dial(network, address string) (net.Conn, error) {
    fmt.Println("dialing beeep boop", network, address)



    // fmt.Println("writing", b)

    done := make(chan int)

    callback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        fmt.Println("finished dialing stuff")
        // callback.Release() // free up memory from callback
        done <- args[0].Int()
        return nil
    })

    // func printMessage(this js.Value, args []js.Value) interface{} {
    //     message := args[0].String()
    //     fmt.Println(message)
    //     
    //     return nil
    // }
    defer callback.Release()
    
    rawHost, rawPort, _ := net.SplitHostPort(address)


    js.Global().Get("dialSocket").Invoke(rawHost, rawPort, callback)

    // wait until we've got our response
    id := <-done

    // TODO: error if id < 0

    // return net.Dial(network, address)
    return &Conn{
        id: id,
    }, nil
}

// LookupHost for regular network uses the net.LookupHost function
func (r *ClearNet) LookupHost(host string) ([]string, error) {
    return net.LookupHost(host)
}

// LookupSRV for regular network uses net.LookupSRV function
func (r *ClearNet) LookupSRV(service, proto, name string) (string, []*net.SRV, error) {
    return net.LookupSRV(service, proto, name)
}

// ResolveTCPAddr for regular network uses net.ResolveTCPAddr function
func (r *ClearNet) ResolveTCPAddr(network, address string) (*net.TCPAddr, error) {
    return net.ResolveTCPAddr(network, address)
}

// // ProxyNet is an implementation of the Net interface that defines behaviour
// // for Tor network connections.
// type ProxyNet struct {
//     // // SOCKS is the host:port which Tor's exposed SOCKS5 proxy is listening
//     // // on.
//     // SOCKS string

//     // // DNS is the host:port of the DNS server for Tor to use for SRV
//     // // queries.
//     // DNS string

//     // // StreamIsolation is a bool that determines if we should force the
//     // // creation of a new circuit for this connection. If true, then this
//     // // means that our traffic may be harder to correlate as each connection
//     // // will now use a distinct circuit.
//     // StreamIsolation bool
// }

// // Dial uses the Tor Dial function in order to establish connections through
// // Tor. Since Tor only supports TCP connections, only TCP networks are allowed.
// func (p *ProxyNet) Dial(network, address string) (net.Conn, error) {
//     switch network {
//     case "tcp", "tcp4", "tcp6":
//     default:
//         return nil, errors.New("cannot dial non-tcp network via Tor")
//     }
//     return Dial(address, p.SOCKS, p.StreamIsolation)
// }

// // LookupHost uses the Tor LookupHost function in order to resolve hosts over
// // Tor.
// func (p *ProxyNet) LookupHost(host string) ([]string, error) {
//     // return LookupHost(host, p.SOCKS)
// }

// // LookupSRV uses the Tor LookupSRV function in order to resolve SRV DNS queries
// // over Tor.
// func (p *ProxyNet) LookupSRV(service, proto, name string) (string, []*net.SRV, error) {
//     // return LookupSRV(service, proto, name, p.SOCKS, p.DNS, p.StreamIsolation)
    
// }

// // ResolveTCPAddr uses the Tor ResolveTCPAddr function in order to resolve TCP
// // addresses over Tor.
// func (p *ProxyNet) ResolveTCPAddr(network, address string) (*net.TCPAddr, error) {
//     switch network {
//     case "tcp", "tcp4", "tcp6":
//     default:
//         return nil, errors.New("cannot dial non-tcp network via Tor")
//     }
//     // return ResolveTCPAddr(address, p.SOCKS)
//     return nil, nil
// }
