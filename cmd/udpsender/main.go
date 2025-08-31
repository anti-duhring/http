package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")

		s, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}

		_, err = conn.Write([]byte(s))
		if err != nil {
			panic(err)
		}
	}
}
