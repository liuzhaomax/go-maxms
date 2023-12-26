#!/bin/sh
mkdir ../src/data_api/pb
protoc --go_out=../src/data_api/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=../src/data_api/pb \
  --go-grpc_opt=paths=source_relative \
  data.proto