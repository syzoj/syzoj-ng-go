//go:generate protoc -I . model.proto --go_out=grpc=.:.
//go:generate protoc -I . model.proto "--gotag_out=xxx=bson+\"-\",output_path=.:."
package model

// Dependency: github.com/amsokol/protoc-gen-gotag
