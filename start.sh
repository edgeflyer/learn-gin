#!/bin/bash

echo "启动grpc server"
go run cmd/server/ratelimit/main.go &

SERVER_PID=$!

sleep 2

echo "正在启动gin"
go run cmd/main.gp

trap "kill $SERVER_PID" EXIT