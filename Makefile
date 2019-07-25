OUT=bin
CLIENT_DIR=client
FLAGS=-a -v -ldflags '-w -extldflags "-static"'

default: build

install:
	go get github.com/dgrijalva/jwt-go
	go get github.com/go-ini/ini
	go get github.com/pkg/errors

build: cmd/deployment-server/main.go ./client/package.json
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(FLAGS) -o $(OUT)/deployment-server cmd/deployment-server/main.go && echo "Compiled deployment-server"
	npm -C $(CLIENT_DIR) install && echo "Installed client modules"

build-arm: cmd/deployment-server/main.go ./client/package.json
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build $(FLAGS) -o $(OUT)/deployment-server cmd/deployment-server/main.go && echo "Compiled deployment-server"
	npm -C $(CLIENT_DIR) install && echo "Installed client modules"

run: $(OUT)/deployment-server
	$(OUT)/deployment-server

