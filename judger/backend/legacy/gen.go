//go:generate protoc -I. -I../../../app/model model.proto --go_out=plugins=grpc:$GOPATH/src
package legacy
