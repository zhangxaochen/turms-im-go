package main

import (
	"fmt"
	"github.com/pires/go-proxyproto"
	"io/ioutil"
	"net"
	"time"
)

func main() {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	pl := &proxyproto.Listener{Listener: l}
	go func() {
		conn, err := pl.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			return
		}
		b, _ := ioutil.ReadAll(conn)
		fmt.Printf("Received: %s, remote: %s\n", string(b), conn.RemoteAddr())
		conn.Close()
	}()

	time.Sleep(100 * time.Millisecond)
	conn, _ := net.Dial("tcp", l.Addr().String())
	conn.Write([]byte("hello world"))
	conn.Close()

	time.Sleep(time.Second)
}
