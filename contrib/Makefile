
GOPATH := $(shell pwd)
export GOPATH
PATH := ${PATH}:$(shell pwd)/bin
export PATH

build:
	go install github.com/bmeg/doc-index/doc-loader
	go install github.com/bmeg/doc-index/doc-searcher

PROTO_INC= -I ./ 

proto:
	cd src/github.com/bmeg/doc-index/proto && protoc \
			$(PROTO_INC) \
			--go_out=../ \
			index.proto

download:
	go get github.com/bmeg/doc-index
	go get -d github.com/bmeg/doc-index/doc-loader
	go get github.com/golang/protobuf/protoc-gen-go
	
test:
	go test -v github.com/bmeg/doc-index