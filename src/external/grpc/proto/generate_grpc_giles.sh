#!/bin/bash

echo $1

if [[ $1 == "" ]]; then
  echo "Error: proto name is not set"
  exit 1
else
  DETACH_FLAG=$1
fi

protoc -I=./$DETACH_FLAG/ \
  --go_out=./$DETACH_FLAG/ --go_opt=paths=source_relative \
  --go-grpc_out=./$DETACH_FLAG/ --go-grpc_opt=paths=source_relative \
  ./$DETACH_FLAG/*.proto