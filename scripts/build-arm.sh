#!/bin/bash

if [ "$1" == "server" ]
then
    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -ldflags '-w -extldflags "-static"' -o bin/deployment-server-arm7 cmd/deployment-server/main.go
    echo "Compiled deployment-server"
else
    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -ldflags '-w -extldflags "-static"' -o bin/deployment-server cmd/deployment-server/main.go
    echo "Compiled deployment-server"
    npm -C ./client install
    echo "Installed client modules"
fi


