package main

import (
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"sync"

	"github.com/inconshreveable/muxado"
)

func main() {
	if len(os.Args) < 2 {
		panic("Socket path is required")
	}
	if err := os.MkdirAll(path.Dir(os.Args[1]), 0666); err != nil {
		panic(err)
	}
	l, err := net.Listen("unix", os.Args[1])
	if err != nil {
		panic(err)
	}
	sess := muxado.Client(struct {
		io.ReadCloser
		io.Writer
	}{ioutil.NopCloser(os.Stdin), os.Stdout}, nil)
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go func(conn net.Conn) {
			defer conn.Close()
			stream, err := sess.OpenStream()
			if err != nil {
				panic(err)
			}
			defer stream.Close()
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				io.Copy(conn, stream)
				conn.(*net.UnixConn).CloseWrite()
				wg.Done()
			}()
			go func() {
				io.Copy(stream, conn)
				stream.CloseWrite()
				wg.Done()
			}()
			wg.Wait()
		}(conn)
	}
}
