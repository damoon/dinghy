//go:generate protoc -I ../../notify --go_opt=paths=source_relative --go_out=plugins=grpc:../pkg/pb service.proto

package dinghy
