package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	verbose   = flag.Bool("v", false, "Prints debugging messages.")
	port      = flag.Int("p", 8080, "Specifies the port number that the server will listen and serve at.")
	directory = flag.String("d", ".", "Specifies the directory that the server will use to read/write requested files.")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Error listening on port %d: %v", *port, err)
	}
	defer listener.Close()

	if *verbose {
		fmt.Println("httpfs is running in verbose mode")
	}

	// connection-loop:
	// handle incomming requests
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v", err)
			continue
		}
		if *verbose {
			fmt.Print("\n")
			fmt.Println("Connected to ", conn.RemoteAddr())
			fmt.Print("\n")
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := NewRequest(conn)
	if err != nil {
		log.Printf("Error reading the request: %v", err)
		return
	}

	res := NewResponse(conn)
	filepath := *directory + req.Path

	if req.Method == "GET" {
		if req.Path == "/" {
			files, err := readDirectory(*directory)
			if err != nil {
				log.Printf("Error could not read directory: %v", err)
				return
			}
			if err = res.Send(200, strings.Join(files, "\r\n")+"\r\n", ""); err != nil {
				log.Printf("Error could not send response: %v", err)
			}
		} else {
			if err = res.Send(200, "", filepath); err != nil {
				log.Printf("Error could not send response: %v", err)
			}
		}
	} else if req.Method == "POST" {
		if req.Path == "/" {
			if err = res.Send(400, "BAD REQUEST: need to pick filename\r\n", ""); err != nil {
				log.Printf("Error could not send response: %v", err)
			}
			return
		}

		if _, val := req.Headers["Content-Length"]; !val {
			if err = res.Send(400, "Content-Length header is required", ""); err != nil {
				log.Printf("Error could not send response: %v", err)
			}
			return
		}

		l, err := strconv.Atoi(req.Headers["Content-Length"])
		if err != nil {
			log.Printf("Error could not read content-length: %v. value: %v", err, req.Headers["Content-Length"])
			return
		}

		f, err := os.Create(filepath)
		if err != nil {
			log.Printf("Error could not open file %s for writing: %v", req.Path[1:], err)
			return
		}
		defer f.Close()

		var r io.Reader
		if *verbose {
			r = io.TeeReader(req, os.Stdout)
		} else {
			r = req
		}

		if _, err = io.CopyN(f, r, int64(l)); err != nil {
			log.Printf("Error writing to file: %v", err)
			return
		}

		if err = res.SendStatus(200); err != nil {
			log.Printf("could not send response: %v", err)
		}

	}
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
