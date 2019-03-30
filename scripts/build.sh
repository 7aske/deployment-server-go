#!bin/bash

go build -o bin/deployment-server cmd/deployment-server/main.go
echo "Compiled deployment-server"
npm -C ./client install
echo "Installed client modules"
