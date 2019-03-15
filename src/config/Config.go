package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type Config struct {
	port        int
	appsPort    int
	appsRoot    string
	secret      []byte
	pass        string
	user        string
	clientRoot  string
	basicServer string
}

func (c *Config) Write() {
	cwd, _ := os.Getwd()
	cFilePath := path.Join(cwd, "config", "server.cfg")
	cFile, err := ini.Load(cFilePath)
	if err != nil {
		log.Fatal("unable to open ", cFilePath)
	}
	cFile.Section("server").Key("port").SetValue(strconv.Itoa(c.port))
	cFile.Section("server").Key("appsPort").SetValue(strconv.Itoa(c.appsPort))


	cFile.Section("auth").Key("secret").SetValue(string(c.secret))
	cFile.Section("auth").Key("user").SetValue(string(c.user))
	cFile.Section("auth").Key("pass").SetValue(string(c.pass))

	cFile.Section("deployer").Key("root").SetValue(c.appsRoot)
	cFile.Section("deployer").Key("server").SetValue(c.basicServer)

	err = cFile.SaveTo(cFilePath)
	if err != nil {
		fmt.Println("error saving config")
	}
}

func (c *Config) Read() {
	cwd, _ := os.Getwd()
	cFilePath := path.Join(cwd, "config", "server.cfg")
	cFile, err := ini.Load(cFilePath)
	if err != nil {
		log.Fatal("unable to open ", cFilePath)
	}
	port, err := strconv.Atoi(cFile.Section("server").Key("port").Value())
	if err != nil {
		c.port = 8080
	} else {
		c.port = port
	}
	appsPort, err := strconv.Atoi(cFile.Section("server").Key("appsPort").Value())
	if err != nil {
		c.appsPort = 8081
	} else {
		c.appsPort = appsPort
	}

	secret := []byte(cFile.Section("auth").Key("secret").Value())
	pass := cFile.Section("auth").Key("pass").Value()
	user := cFile.Section("auth").Key("user").Value()

	c.user = user
	c.secret = secret
	c.pass = pass

	pth := cFile.Section("deployer").Key("root").Value()
	if pth == "" {
		c.appsRoot = path.Join(cwd, "apps")
	} else {
		if filepath.IsAbs(pth){
			c.appsRoot = pth

		} else {
			c.appsRoot = path.Join(cwd, pth)
		}
	}
	server := cFile.Section("deployer").Key("server").Value()
	if server == "" {
		c.basicServer = path.Join(cwd, "server", "server.js")
	} else {
		if filepath.IsAbs(server){
			c.basicServer = server

		} else {
			c.basicServer = path.Join(cwd, server)
		}
	}
}

func LoadConfig() *Config {
	config := Config{}
	config.Read()
	return &config
}

func (c *Config) SetPort(port int) {
	c.port = port
}
func (c *Config) GetPort() int {
	return c.port
}
func (c *Config) SetAppsPort(port int) {
	c.appsPort = port
}
func (c *Config) GetAppsPort() int {
	return c.appsPort
}
func (c *Config) SetSecret(secret []byte) {
	c.secret = secret
}
func (c *Config) GetSecret() []byte {
	return c.secret
}
func (c *Config) SetAppsRoot(pth string) {
	c.appsRoot = pth
}
func (c *Config) GetAppsRoot() string {
	return c.appsRoot
}
func (c *Config) SetPass(pass string) {
	c.pass = pass
}
func (c *Config) GetPass() string {
	return c.pass
}
func (c *Config) SetUser(user string) {
	c.user = user
}
func (c *Config) GetUser() string {
	return c.user
}
func (c *Config) SetBasicServer(server string) {
	c.basicServer = server
}
func (c *Config) GetBasicServer() string {
	return c.basicServer
}

//func (c *Config) SetClientRoot(clientRoot string) {
//	c.clientRoot = clientRoot
//}
//func (c *Config) GetClientRoot() string {
//	return c.clientRoot
//}
