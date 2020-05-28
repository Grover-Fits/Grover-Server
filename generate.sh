#!/bin/bash

protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/ \
--go_out=plugins=grpc,paths=source_relative:. --grpc-gateway_out=logtostderr=true:. api/proto/api.proto