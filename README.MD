deployment-server
==========================


## Description
Deployment server

## Usage

`bin/deployment-server`

Config file
 
 `config/server.cfg`
```ini
[dev]
; port for dev server and apps
appsPort = 30001
; should not be the same
port     = 30000

[router]
; router port
port = 80

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

####Additional options

## Package

Basic structure of package is

## Requirements

None

## Installation

## Conclusion

