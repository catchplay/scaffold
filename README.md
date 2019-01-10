# scaffold

[![Build Status](https://travis-ci.org/catchplay/scaffold.svg)](https://travis-ci.org/catchplay/scaffold)
[![codecov](https://codecov.io/gh/catchplay/scaffold/branch/master/graph/badge.svg)](https://codecov.io/gh/catchplay/scaffold)
[![Go Report Card](https://goreportcard.com/badge/github.com/catchplay/scaffold)](https://goreportcard.com/report/github.com/catchplay/scaffold)
[![GoDoc](https://godoc.org/github.com/catchplay/scaffold?status.svg)](https://godoc.org/github.com/catchplay/scaffold)

Scaffold generates starter Go project layout. Let you can focus on  buesiness logic implemeted. 

[![asciicast](https://asciinema.org/a/MA0ppdKfZSEl64cskUnqfsSiH.svg)](https://asciinema.org/a/MA0ppdKfZSEl64cskUnqfsSiH?autoplay=1&speed=2)

The following is Go project layout scaffold generated:

```
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── main.go
├── config
│   ├── config.go
│   ├── config.yml
│   ├── database.go
│   ├── http.go
│   └── release.go
├── docker-compose.yml
├── model
│   └── model.go
└── web
    ├── routes.go
    ├── server.go
    └── version.go
```

## Installation

 Download scaffold by using:
```sh
$ go get -u github.com/catchplay/scaffold
```

## Create a new project

1. Going to your new project folder:
```sh
# change to project directory
$ cd $GOPATH/src/path/to/project
```

2. Run `scaffold init` in the new project folder:

```sh
$ scaffold init
```

3. That will generate a whole new starter project files like:

```
Create Dockerfile
Create README.md
Create cmd/main.go
Create config/config.go
Create config/database.go
Create config/http.go
Create config/release.go
Create docker-compose.yml
Create model/model.go
Create web/routes.go
Create web/server.go
Create web/version.go
Create Makefile
Create config/config.yml
Success Created. Please excute `make up` to start service.

```

4. And you can run the new project by using:
```sh
$ make run 
```