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
	cwd         string
	port        int
	appsPort    int
	hostname    string
	routerPort  int
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
		fp, _ := os.Create(cFilePath)
		_, _ = fp.Write([]byte{})
		fp.Close()
		cFile, err = ini.Load(cFilePath)
		if err != nil {
			log.Fatal("unable to open ", cFilePath)
		}
		_ = cFile.NewSections("dev", "router", "auth", "deployer")
		cFile.Section("dev").Key("appsPort").Comment = "port for dev server and apps"
		cFile.Section("dev").Key("port").Comment = "should not be the same"
		cFile.Section("router").Key("port").Comment = "router port"
		cFile.Section("auth").Key("secret").Comment = "hash secret"
		cFile.Section("auth").Key("user").Comment = "username and password required to log to dev server"
		cFile.Section("deployer").Key("root").Comment = "location of where the app repos are stored relative to app root"
		cFile.Section("deployer").Key("server").Comment = "location of nodejs server script that runs 'web' apps"
		cFile.Section("deployer").Key("hostname").Comment = "default hostname used for parsing subdomains, ignored if on local"
	}
	cFile.Section("dev").Key("appsPort").SetValue(strconv.Itoa(c.appsPort))
	cFile.Section("dev").Key("port").SetValue(strconv.Itoa(c.port))

	cFile.Section("router").Key("port").SetValue(strconv.Itoa(c.routerPort))

	cFile.Section("auth").Key("secret").SetValue(string(c.secret))
	cFile.Section("auth").Key("user").SetValue(string(c.user))
	cFile.Section("auth").Key("pass").SetValue(string(c.pass))

	cFile.Section("deployer").Key("root").SetValue(c.appsRoot)
	cFile.Section("deployer").Key("server").SetValue(c.basicServer)
	cFile.Section("deployer").Key("hostname").SetValue(c.hostname)

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
		fmt.Println(err)
		c.port = 30000
		c.appsPort = 30001

		c.routerPort = 8080

		c.user = "admin"
		c.pass = "admin"
		c.secret = []byte("secret")

		c.appsRoot = "apps"
		c.basicServer = "server/server.js"
		c.hostname = "127.0.0.1"
		c.Write()
	} else {
		port, err := strconv.Atoi(cFile.Section("dev").Key("port").Value())
		if err != nil {
			c.port = 30000
			c.Write()
		} else {
			c.port = port
		}
		appsPort, err := strconv.Atoi(cFile.Section("dev").Key("appsPort").Value())
		if err != nil {
			c.appsPort = 30001
			c.Write()
		} else {
			c.appsPort = appsPort
		}
		routerPort, err := strconv.Atoi(cFile.Section("router").Key("port").Value())
		if err != nil {
			c.routerPort = 8080
			c.Write()
		} else {
			c.routerPort = routerPort
		}
		secret := []byte(cFile.Section("auth").Key("secret").Value())
		pass := cFile.Section("auth").Key("pass").Value()
		user := cFile.Section("auth").Key("user").Value()

		if user == "" {
			user = "admin"
			c.Write()
		}
		if len(secret) == 0 {
			secret = []byte("secret")
			c.Write()
		}
		if pass == "" {
			pass = "admin"
			c.Write()
		}
		c.user = user
		c.secret = secret
		c.pass = pass

		pth := cFile.Section("deployer").Key("root").Value()
		if pth == "" {
			c.appsRoot = "apps"
			c.Write()
		} else {
			if filepath.IsAbs(pth) {
				c.appsRoot = path.Dir(pth)
			} else {
				c.appsRoot = pth
			}
		}
		server := cFile.Section("deployer").Key("server").Value()
		if server == "" {
			c.basicServer = "server/server.js"
			c.Write()
		} else {
			// TODO: this is probably wrong
			if filepath.IsAbs(server) {
				c.basicServer = path.Base(server) + "/server.js"
			} else {
				c.basicServer = server
			}
		}
		hostname := cFile.Section("deployer").Key("hostname").Value()
		c.hostname = hostname
	}

}

func LoadConfig() *Config {
	config := Config{}
	cwd, _ := os.Getwd()
	config.SetCwd(cwd)
	config.Read()
	return &config
}
func (c *Config) SetCwd(cwd string) {
	c.cwd = cwd
}
func (c *Config) GetCwd() string {
	return c.cwd
}

func (c *Config) SetPort(port int) {
	c.port = port
}
func (c *Config) GetPort() int {
	return c.port
}
func (c *Config) SetHostname(hostname string) {
	c.hostname = hostname
}
func (c *Config) GetHostname() string {
	return c.hostname
}
func (c *Config) SetRouterPort(port int) {
	c.port = port
}
func (c *Config) GetRouterPort() int {
	return c.routerPort
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
