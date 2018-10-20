# httpfs
httpfs is a simple file server

## Requirements
1. Go 1.7 or later

## Detailed Usage

### General

``` bash
httpfs is a simple file server.

usage: httpfs [-v] [-p PORT] [-d PATH-TO-DIR]
    -v Prints debugging messages.
    -p Specifies the port number that the server will listen and serve at. Default is 8080.
    -d Specifies the directory that the server will use to read/write requested files. Default is the current directory when launching the application.
```