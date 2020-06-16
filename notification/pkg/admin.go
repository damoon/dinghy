//go:generate protoc -I ../. --go_out=plugins=grpc:../pkg service.proto

package notify
