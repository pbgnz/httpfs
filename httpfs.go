package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
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
			fmt.Println("Connected to ", conn.RemoteAddr())
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
	log.Print(req)
}

type Request struct {
	Method   string
	Path     string
	Headers  map[string]string
	Protocol string
}

func handleRequest(conn net.Conn) (*Request, error) {
	defer conn.Close()

	req := &Request{
		Headers: map[string]string{},
	}

	request := bufio.NewReader(conn)
	requestLine, err := request.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error in the request line: %v", err)
	}

	if *verbose {
		fmt.Println(requestLine)
	}

	rl := strings.Fields(requestLine)
	if len(rl) != 3 {
		return nil, fmt.Errorf("Error in the request line %v", err)
	}

	if rl[0] == "GET" || rl[0] == "POST" {
		req.Method = rl[0]
	} else {
		return nil, fmt.Errorf("Error in the request method %v", err)
	}
	req.Path = rl[1]
	req.Protocol = rl[2]

	for {
		headerLine, err := request.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("Error in the header line %v", err)
		}
		if *verbose {
			fmt.Println(headerLine)
		}
		if headerLine == "\r\n" {
			break
		}
		parts := regexp.MustCompile(`^([\w-]+): (.+)\r\n$`).FindStringSubmatch(headerLine)
		if len(parts) != 3 {
			return nil, fmt.Errorf("Error in the header lines %v", err)
		}
		req.Headers[parts[1]] = parts[2]
	}
	return req, nil
}
