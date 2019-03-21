//go:generate protoc -I . model.proto --go_out=grpc=$GOPATH/src:$GOPATH/src
//go:generate protoc -I . model.proto --gotype_out=.
// OLD go:generate protoc -I . model.proto "--gotag_out=xxx=bson+\"-\",output_path=.:."
//go:generate protoc -I . api.proto --go_out=grpc=$GOPATH/src:$GOPATH/src
//go:generate protoc -I . api.proto --gotype_out=.
package model

// Dependency: github.com/syzoj/protoc-gen-gotype
