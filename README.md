# scaffold

[![Build Status](https://travis-ci.org/catchplay/scaffold.svg)](https://travis-ci.org/catchplay/scaffold)
[![codecov](https://codecov.io/gh/catchplay/scaffold/branch/master/graph/badge.svg)](https://codecov.io/gh/catchplay/scaffold)
[![Go Report Card](https://goreportcard.com/badge/github.com/catchplay/scaffold)](https://goreportcard.com/report/github.com/catchplay/scaffold)
[![GoDoc](https://godoc.org/github.com/catchplay/scaffold?status.svg)](https://godoc.org/github.com/catchplay/scaffold)

Scaffold generates starter Go project layout. Let you can focus on  buesiness logic implemeted. 

The following is Go project layout scaffold generated:

```
├── cmd
│   └── main.go
├── config
│   ├── config.go
│   ├── config.yml
│   ├── database.go
│   ├── http.go
│   └── release.go
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
$ cd $GOPATH/src/$USER_NAME/$YOUR_PROJECT
```

2. Run `scaffold init` in the new project folder:

```sh
$ scaffold init
```

3. That will generate a whole new starter project files:

```

```

4. And you can the new project by busing:
```
sh
$ make run 
```｀