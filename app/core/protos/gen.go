//go:generate protoc -I. -I../../model judge.proto --go_out=plugins=grpc:$GOPATH/src
package protos

// Dependency: github.com/amsokol/protoc-gen-gotag
