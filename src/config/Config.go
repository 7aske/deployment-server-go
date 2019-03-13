package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path"
	"strconv"
)

type Config struct {
	port     int
	appsPort int
	appsPath string
	secret   []byte
	pass     string
	user     string
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
func (c *Config) SetAppsPath(pth string) {
	c.appsPath = pth
}
func (c *Config) GetAppsPath() string {
	return c.appsPath
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

func (c *Config) Write() {
	cwd, _ := os.Getwd()
	cFilePath := path.Join(cwd, "config", "server.cfg")
	cFile, err := ini.Load(cFilePath)
	if err != nil {
		log.Fatal("unable to open ", cFilePath)
	}
	cFile.Section("server").Key("path").SetValue(c.appsPath)
	cFile.Section("server").Key("port").SetValue(string(c.port))
	cFile.Section("auth").Key("secret").SetValue(string(c.secret))
	cFile.Section("auth").Key("pass").SetValue(string(c.pass))

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
		c.SetPort(8080)
	} else {
		c.SetPort(port)
	}
	c.SetAppsPort(c.GetPort() + 1)
	pth := cFile.Section("server").Key("path").Value()
	if pth == "" {
		c.SetAppsPath(path.Join(cwd, "apps"))
	} else {
		c.SetAppsPath(path.Join(cwd, pth))
	}
	secret := []byte(cFile.Section("auth").Key("secret").Value())
	pass := cFile.Section("auth").Key("pass").Value()
	user := cFile.Section("auth").Key("user").Value()
	c.SetUser(user)
	c.SetSecret(secret)
	c.SetPass(pass)
}

func LoadConfig() Config {
	config := Config{}
	config.Read()
	return config
}
