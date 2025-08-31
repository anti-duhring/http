package main

import (
	"bytes"
	"fmt"
	"net"
)

func getLinesChannel(conn net.Conn) <-chan string {
	ch := make(chan string)

	str := ""
	go func() {
		defer close(ch)
		defer conn.Close()

		for {
			b := make([]byte, 8)
			_, err := conn.Read(b)
			if err != nil {
				break
			}

			if i := bytes.IndexByte(b, '\n'); i != -1 {
				str += string(b[:i])
				b = b[i+1:]
				ch <- str
				str = ""
			}

			str += string(b)
		}
	}()

	return ch
}

func main() {
	lis, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}

		ch := getLinesChannel(conn)
		for v := range ch {
			fmt.Printf("%s\n", v)
		}
	}
}
