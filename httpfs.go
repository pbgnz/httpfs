package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	verbose = flag.Bool("v", false, "Prints debugging messages.")
	port    = flag.Int("p", 8080, "Specifies the port number that the server will listen and serve at.")
	// directory = flag.String("d", "", "Specifies the directory that the server will use to read/write requested files.")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen %s with %v\n", addr, err)
		return
	}

	defer listener.Close()
	fmt.Println("httpfs is listening at", listener.Addr())

	// connection-loop:
	// handle incomming requests
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("Error: failed to close listener:", err)
			}
			continue
		}

		if *verbose {
			log.Println("Connected to", conn.RemoteAddr())
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := handleRequest(conn)
	if err != nil {
		log.Printf("Error reading the request: %v", err)
		return
	}
	log.Println(req)
}

func handleRequest(conn net.Conn) (*Request, error) {
	defer conn.Close()

	requestLine, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error in the request line: %v", err)
	}

	if *verbose {
		fmt.Println(requestLine)
	}

	return nil, nil
}
