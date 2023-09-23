VERSION := $(shell git describe --tags 2>/dev/null)
LDFLAGS = -X main.version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o dbanon main.go

bench:
	$$GOPATH/bin/go-bindata -pkg bindata -o bindata/bindata.go etc/*
	go test -run=XXX -bench=. -benchtime=20s $$GOPATH/src/github.com/mdshack/dbanon/src
	rm -rf bindata