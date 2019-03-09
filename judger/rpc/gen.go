//go:generate protoc -I. -I../../app/model judge_rpc.proto --go_out=plugins=grpc:$GOPATH/src
package rpc
