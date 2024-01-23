#!/bin/sh

CONTRACT_PATH=user.proto
OUT_PATH=../src/api_user_rpc/pb

cd spec
mkdir -p ${OUT_PATH}
protoc --go_out=${OUT_PATH} \
  --go_opt=paths=source_relative \
  --go-grpc_out=require_unimplemented_servers=false:${OUT_PATH} \
  --go-grpc_opt=paths=source_relative \
  ${CONTRACT_PATH}
cd ..
