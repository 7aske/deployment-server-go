package config

import (
	"../utils"
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	APPS_ROOT    = "apps"
	BASIC_SERVER = "server/server.js"
	PATH         = "config/config.cfg"
)

// Config fields enum
var FIELDS = []string{
	"APPSROOT",
	"PORT",
	"APPSPORT",
	"ROUTERPORT",
	"SECRET",
	"PASS",
	"USER",
	"HOSTNAME",
	"CONTAINER",
}

type Config struct {
	cwd         string
	basicServer string
	appsRoot    string
	port        int
	appsPort    int
	routerPort  int
	secret      string
	pass        string
	user        string
	hostname    string
	container   bool
}

func New() *Config {
	config := Config{}
	absPath := utils.GetAbsDir(os.Args[0])
	config.cwd = path.Dir(path.Dir(absPath))
	config.setDefault()
	config.Read()

	if config.container {
		config.container = utils.VerifyCcont()
	}
	fmt.Printf("%v\r\n", config.container)
	return &config
}

func (c *Config) Write() {
	cFilePath := path.Join(c.GetCwd(), PATH)
	cFile, err := ini.Load(cFilePath)
	if err != nil {
		utils.MakeFileIfNotExist(cFilePath)
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
		cFile.Section("deployer").Key("hostname").Comment = "default hostname used for parsing subdomains, ignored if on local"
		cFile.Section("deployer").Key("container").Comment = "toggle whether to use ccont containers"
	}
	cFile.Section("dev").Key("appsPort").SetValue(strconv.Itoa(c.appsPort))
	cFile.Section("dev").Key("port").SetValue(strconv.Itoa(c.port))
	cFile.Section("router").Key("port").SetValue(strconv.Itoa(c.routerPort))
	cFile.Section("auth").Key("secret").SetValue(c.secret)
	cFile.Section("auth").Key("user").SetValue(c.user)
	cFile.Section("auth").Key("pass").SetValue(c.pass)
	cFile.Section("deployer").Key("hostname").SetValue(c.hostname)
	cFile.Section("deployer").Key("container").SetValue(strconv.FormatBool(c.container))

	err = cFile.SaveTo(cFilePath)
	if err != nil {
		fmt.Printf("error saving config %v\r\n", err.Error())
	}
}

func (c *Config) Read() {
	cwd, _ := os.Getwd()
	cFilePath := path.Join(cwd, PATH)
	cFile, err := ini.Load(cFilePath)
	if err != nil {
		fmt.Print("no config file found - generating default config\r\n")
		c.setDefault()
		c.Write()
	} else {
		port, err := strconv.Atoi(cFile.Section("dev").Key("port").Value())
		if err != nil {
			c.port = 30000
		} else {
			c.port = port
		}
		appsPort, err := strconv.Atoi(cFile.Section("dev").Key("appsPort").Value())
		if err != nil {
			c.appsPort = 30001
		} else {
			c.appsPort = appsPort
		}
		routerPort, err := strconv.Atoi(cFile.Section("router").Key("port").Value())
		if err != nil {
			c.routerPort = 8080
		} else {
			c.routerPort = routerPort
		}
		secret := cFile.Section("auth").Key("secret").Value()
		pass := cFile.Section("auth").Key("pass").Value()
		user := cFile.Section("auth").Key("user").Value()

		if user == "" {
			user = "admin"
		}
		if len(secret) == 0 {
			secret = "secret"
		}
		if pass == "" {
			pass = "admin"
		}
		c.user = user
		c.secret = secret
		c.pass = pass

		hostname := cFile.Section("deployer").Key("hostname").Value()
		c.hostname = hostname

		container, _ := strconv.ParseBool(cFile.Section("deployer").Key("container").Value())
		c.container = container
		c.Write()
	}

}

func (c *Config) SetCwd(cwd string) {
	c.cwd = cwd
	c.Write()
}
func (c *Config) GetCwd() string {
	return c.cwd
}
func (c *Config) SetPort(port int) {
	c.port = port
	c.Write()
}
func (c *Config) GetPort() int {
	return c.port
}
func (c *Config) SetHostname(hostname string) {
	c.hostname = hostname
	c.Write()
}
func (c *Config) GetHostname() string {
	return c.hostname
}
func (c *Config) SetRouterPort(port int) {
	c.port = port
	c.Write()
}
func (c *Config) GetRouterPort() int {
	return c.routerPort
}
func (c *Config) SetAppsPort(port int) {
	c.appsPort = port
	c.Write()
}
func (c *Config) GetAppsPort() int {
	return c.appsPort
}
func (c *Config) SetSecret(secret string) {
	c.secret = secret
	c.Write()
}
func (c *Config) GetSecret() string {
	return c.secret
}

func (c *Config) SetPass(pass string) {
	c.pass = pass
	c.Write()
}
func (c *Config) GetPass() string {
	return c.pass
}
func (c *Config) SetUser(user string) {
	c.user = user
	c.Write()
}
func (c *Config) GetUser() string {
	return c.user
}
func (c *Config) setAppsRoot(pth string) {
	c.appsRoot = pth
	c.Write()
}
func (c *Config) GetAppsRoot() string {
	return c.appsRoot
}
func (c *Config) setBasicServer(server string) {
	c.basicServer = server
	c.Write()
}
func (c *Config) GetBasicServer() string {
	return c.basicServer
}
func (c *Config) SetContainer(cont bool) {
	c.container = cont
	c.Write()
}
func (c *Config) GetContainer() bool {
	return c.container
}

func (c *Config) setDefault() {
	c.port = 30000
	c.appsPort = 30001
	c.routerPort = 8080
	c.user = "admin"
	c.pass = "admin"
	c.secret = "secret"
	c.appsRoot = APPS_ROOT
	c.basicServer = BASIC_SERVER
	c.hostname = "127.0.0.1"
	c.container = false
}

func (c *Config) Set(key string, value string) bool {
	switch strings.ToUpper(key) {
	case FIELDS[0]:
		c.setAppsRoot(value)
	case FIELDS[1]:
		port, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		c.SetPort(port)
	case FIELDS[2]:
		port, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		c.SetAppsPort(port)
	case FIELDS[3]:
		port, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		c.SetRouterPort(port)
	case FIELDS[4]:
		c.SetSecret(value)
	case FIELDS[5]:
		c.SetPass(value)
	case FIELDS[6]:
		c.SetUser(value)
	case FIELDS[7]:
		c.SetHostname(value)
	case FIELDS[8]:
		if b, err := strconv.ParseBool(value); err == nil {
			c.SetContainer(b)
		}
	default:
		return false
	}
	return true
}

func (c *Config) PrintFields() {
	fmt.Printf("\r\n")
	for _, field := range FIELDS {
		fmt.Printf("%s\r\n", field)
	}
	fmt.Printf("\r\n")
}
