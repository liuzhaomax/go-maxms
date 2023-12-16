#!/bin/sh

protoc -I . data.proto --go_out=plugins=grpc:../src