//go:generate protoc -I. -I../../../app/model model.proto --go_out=plugins=grpc:$GOPATH/src
//go:generate protoc -I. -I../../../app/model model.proto "--gotag_out=xxx=bson+\"-\",output_path=.:."
package legacy
