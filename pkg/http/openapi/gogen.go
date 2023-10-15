package api

//go:generate protoc --plugin=../../../bin/protoc-gen-openapi --openapi_out=yaml=true,strict_proto3_optional=true,use_ref=true,per_file=true:. --proto_path=../../proto/ autoscaler.proto
//go:generate oapi-codegen -config gen-config.yaml api.yaml
