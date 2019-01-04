# scaffold

[![Build Status](https://travis-ci.org/catchplay/scaffold.svg)](https://travis-ci.org/catchplay/scaffold)
[![codecov](https://codecov.io/gh/catchplay/scaffold/branch/master/graph/badge.svg)](https://codecov.io/gh/catchplay/scaffold)
[![Go Report Card](https://goreportcard.com/badge/github.com/catchplay/scaffold)](https://goreportcard.com/report/github.com/catchplay/scaffold)
[![GoDoc](https://godoc.org/github.com/catchplay/scaffold?status.svg)](https://godoc.org/github.com/catchplay/scaffold)

scaffold generates standard starter Go project layout at CATCHPLAY.

Starter Go project layout as bellow:

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
├── subtitle
└── web
    ├── routes.go
    ├── server.go
    └── version.go
```