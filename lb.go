package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var (
	counter int
	//TODO configurable
	listenAddr = "localhost:8080"

	//TODO configuralbe
	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}
)

func main() {
	listener, err := net.Listen("tcp", listenAddr)

	if err != nil {
		log.Fatal("failed to listen %s", err)

	}
	defer listener.Close()

	for {

		conn, err := listener.Accept()

		if err != nil {
			log.Printf("failed to accept connection: %s", err)
		}

		beckend := chooseBeckend()
		fmt.Printf("counter=%d beckend=%s\n", counter, beckend)

		go func() {
			err := proxy(beckend, conn)
			if err != nil {
				log.Printf("WARNING: proxying failed: %v", err)

			}
		}()
	}
}

func proxy(beckend string, c net.Conn) error {
	bc, err := net.Dial("tcp", beckend)
	if err != nil {
		return fmt.Errorf("failed to connect to beckend %s : %v ", beckend, err)

	}

	go io.Copy(bc, c)

	go io.Copy(c, bc)

	return nil

}

func chooseBeckend() string {
	s := server[counter%len(server)]
	counter++
	return s

}
