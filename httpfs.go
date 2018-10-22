package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	verbose   = flag.Bool("v", false, "Prints debugging messages.")
	port      = flag.Int("p", 8080, "Specifies the port number that the server will listen and serve at.")
	directory = flag.String("d", ".", "Specifies the directory that the server will use to read/write requested files.")
)

type Request struct {
	Method   string
	Path     string
	Headers  map[string]string
	Protocol string
}

type Response struct {
	headers map[string]string
	status  int
}

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

	// handle the client's request
	req, err := handleRequest(conn)
	if err != nil {
		log.Printf("Error reading the request: %v", err)
		return
	}

	// prepare the response
	res := &Response{
		headers: map[string]string{},
	}
	// default headers
	res.status = 200
	res.headers["Date"] = time.Now().Format(time.UnixDate)
	res.headers["Connection"] = "close"

	if req.Method == "GET" {
		if req.Path == "/" {
			files, err := readDirectory(*directory)
			if err != nil {
				return
			}
			fmt.Println(files)
		}
	}
}

func handleRequest(conn net.Conn) (*Request, error) {
	defer conn.Close()

	req := &Request{
		Headers: map[string]string{},
	}

	request := bufio.NewReader(conn)

	// parse the request line
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

	// parse the headers
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

func readDirectory(d string) ([]string, error) {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, fmt.Errorf("Error reading the directory %v", err)
	}
	fileList := []string{}
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}
	return fileList, nil
}
