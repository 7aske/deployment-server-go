# deployment-server

## Description

Custom web-server tooled to easily deploy other web based projects (apps) using Git, inspired by Heroku. Deployment-server is designed to be a lightweight and simple to use.

## Installation

Project is being managed by `make` since Golang is, to me, pretty weird when it comes to the whole GOROOT thing. In order to install dependencies (NOTE: GOROOT must be properly set up) run:

```
make install
```

After installing dependencies build the server using:

```
make build
```

And lastly run the server with:

```
make run
```

## Usage

```
bin/deployment-server
```

This will generate the default config file and test run the server. Default port is port `8080`.

### Config file

`config/server.cfg`

```ini
[dev]
; port for dev server and apps
appsPort = 30001
; should not be the same
port     = 30000

[router]
; router port
port = 8080

[auth]
; hash secret
secret = hashingsecret
; username and password required to log to dev server
user   = admin
pass   = admin

[deployer]
; location of where the app repos are stored relative to app root
root     = apps
; location of nodejs server script that runs 'web' apps
server   = server/server.js
; default hostname used for parsing subdomains, ignored if on local
hostname = example.com
```

### Additional options

## Package

Basic structure of package is

## Requirements

None

## Installation

## Conclusion

