#!bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -extldflags "-static"' -o bin/deployment-server cmd/deployment-server/main.go
echo "Compiled deployment-server"
npm -C ./client install
echo "Installed client modules"
