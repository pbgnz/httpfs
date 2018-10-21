# httpfs
httpfs is a simple file server

## Requirements
1. Go 1.7 or later

## Detailed Usage

### General

#### server

``` bash
httpfs is a simple file server.

usage: httpfs [-v] [-p PORT] [-d PATH-TO-DIR]
    -v Prints debugging messages.
    -p Specifies the port number that the server will listen and serve at. 
    Default is 8080.
    -d Specifies the directory that the server will use to read/write requested files. 
    Default is the current directory when launching the application.
```

#### client

GET / returns a list of the current files in the data directory.

``` bash
curl -get localhost:8080/

httpc -p 8080 get 'http://localhost:8080/'
```

GET /foo.txt returns the content of the file named foo.txt in the data directory.
``` bash
curl -get localhost:8080/foo.txt

httpc -v -p 8080 get 'http://localhost:8080/foo.txt'
```

POST /bar should create or overwrite the file named bar in the data directory with
the content of the body of the request.
``` bash
curl -post -d "foo" localhost:8080/bar.txt

httpc -p 8080 -h Content-Type:application/text -d 'foo' post 'http://localhost:8080/bar.txt'
```