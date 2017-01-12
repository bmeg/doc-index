
GOPATH := $(shell pwd)/buildtools:$(shell pwd)
export GOPATH
PATH := ${PATH}:$(shell pwd)/bin:$(shell pwd)/buildtools/bin/
export PATH

build:
	go install github.com/bmeg/doc-index/doc-loader
	go install github.com/bmeg/doc-index/doc-searcher

depends:
	go get -d github.com/bmeg/doc-index/doc-loader

PROTO_INC= -I ./ 

proto:
	go get github.com/golang/protobuf/protoc-gen-go/
	cd src/github.com/bmeg/doc-index/proto && protoc \
			$(PROTO_INC) \
			--go_out=../ \
			index.proto

test:
	go test -v github.com/bmeg/doc-index