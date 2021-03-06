OUT=bin
NAME=deployment_server
MAIN=cmd/deployment-server/main.go
CLIENT_DIR=client
FLAGS=-a -v -ldflags '-w -extldflags "-static"'

default: build

install:
	sudo ln -sf $(shell pwd)/$(OUT)/$(NAME) /usr/bin/$(NAME)

dep:
	go get github.com/dgrijalva/jwt-go
	go get github.com/go-ini/ini
	go get github.com/pkg/errors
	go get github.com/teris-io/shortid
	go get golang.org/x/crypto/ssh/terminal

build: cmd/deployment-server/main.go ./client/package.json
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(FLAGS) -o $(OUT)/$(NAME) $(MAIN) && \
	npm -C $(CLIENT_DIR) install

build-arm: $(MAIN) ./$(CLIENT_DIR)/package.json
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(FLAGS) -o $(OUT)/$(NAME) $(MAIN) && \
	npm -C $(CLIENT_DIR) install

client:
	npm -C $(CLIENT_DIR) install


run: $(OUT)/deployment_server
	go build -o $(OUT)/$(NAME) $(MAIN) && \
	./$(OUT)/$(NAME) -i

