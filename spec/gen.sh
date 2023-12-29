#!/bin/sh
cd spec
mkdir ../src/api_user/pb
protoc --go_out=../src/api_user/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=../src/api_user/pb \
  --go-grpc_opt=paths=source_relative \
  user.proto
cd ..