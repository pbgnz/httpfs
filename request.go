package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// Request struct represents an HTTP Request
type Request struct {
	Method   string
	Path     string
	Headers  map[string]string
	Protocol string
	*bufio.Reader
}

// NewRequest parses an HTTP Request and returns
// a Request struct
func NewRequest(conn net.Conn) (*Request, error) {
	req := &Request{
		Headers: map[string]string{},
		Reader:  bufio.NewReader(conn),
	}

	// parse the request-line
	requestLine, err := req.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error in the request line: %v", err)
	}

	if *verbose {
		fmt.Print(requestLine)
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
		headerLine, err := req.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("Error in the header line %v", err)
		}
		if *verbose {
			fmt.Print(headerLine)
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
