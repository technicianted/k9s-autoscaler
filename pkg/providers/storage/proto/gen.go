package proto

//go:generate protoc --go_out=. --go_opt=paths=source_relative --plugin=$GOPATH/bin/protoc-gen-go -I . -I ../../../../ inline.proto
