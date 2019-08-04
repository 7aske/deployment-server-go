package config

import (
	"../utils"
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path"
	"strconv"
)

const (
	APPS_ROOT    = "apps"
	BASIC_SERVER = "server/server.js"
	CONFIG_PATH = "config/config.cfg"
)

type Config struct {
	cwd         string
	clientRoot  string
	port        int
	appsPort    int
	routerPort  int
	secret      []byte
	pass        string
	user        string
	basicServer string
	appsRoot    string
	hostname    string
	container   bool
}

func New() *Config {
	config := Config{}
	absPath := utils.GetAbsDir(os.Args[0])
	config.SetCwd(path.Dir(path.Dir(absPath)))
	setDefaultConfig(&config)
	config.Read()

	if config.container {
		config.SetContainer(utils.VerifyCcont())
	}
	fmt.Println(config.container)
	return &config
}

func (c *Config) Write() {
	cFilePath := path.Join(c.GetCwd(), CONFIG_PATH)
	cFile, err := ini.Load(cFilePath)
	if err != nil {
		_ = os.MkdirAll(path.Dir(cFilePath), 0775)
		fp, _ := os.Create(cFilePath)
		_, _ = fp.Write([]byte{})
		_ = fp.Close()
		if os.Getuid() == 0 {
			uid, _ := strconv.Atoi(os.Getenv("SUDO_UID"))
			gid, _ := strconv.Atoi(os.Getenv("SUDO_GID"))
			_ = os.Chown(cFilePath, uid, gid)
		}

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
		//cFile.Section("deployer").Key("root").Comment = "location of where the app repos are stored relative to app root"
		//cFile.Section("deployer").Key("server").Comment = "location of nodejs server script that runs 'web' apps"
	}
	cFile.Section("dev").Key("appsPort").SetValue(strconv.Itoa(c.appsPort))
	cFile.Section("dev").Key("port").SetValue(strconv.Itoa(c.port))

	cFile.Section("router").Key("port").SetValue(strconv.Itoa(c.routerPort))

	cFile.Section("auth").Key("secret").SetValue(string(c.secret))
	cFile.Section("auth").Key("user").SetValue(string(c.user))
	cFile.Section("auth").Key("pass").SetValue(string(c.pass))

	cFile.Section("deployer").Key("hostname").SetValue(c.hostname)
	cFile.Section("deployer").Key("container").SetValue(strconv.FormatBool(c.container))
	//cFile.Section("deployer").Key("root").SetValue(c.appsRoot)
	//cFile.Section("deployer").Key("server").SetValue(c.basicServer)

	err = cFile.SaveTo(cFilePath)
	if err != nil {
		fmt.Println("error saving config", err.Error())
	}
}

func (c *Config) Read() {
	cwd, _ := os.Getwd()
	cFilePath := path.Join(cwd, CONFIG_PATH)
	cFile, err := ini.Load(cFilePath)
	if err != nil {
		fmt.Println("no config file found - generating default config")
		setDefaultConfig(c)
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

		hostname := cFile.Section("deployer").Key("hostname").Value()
		c.hostname = hostname

		container, _ := strconv.ParseBool(cFile.Section("deployer").Key("container").Value())
		c.container = container
	}

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
func (c *Config) setAppsRoot(pth string) {
	c.appsRoot = pth
}
func (c *Config) GetAppsRoot() string {
	return c.appsRoot
}
func (c *Config) setBasicServer(server string) {
	c.basicServer = server
}
func (c *Config) GetBasicServer() string {
	return c.basicServer
}
func (c *Config) SetContainer(cont bool) {
	c.container = cont
}
func (c *Config) GetContainer() bool {
	return c.container
}

func setDefaultConfig(c *Config) {
	c.port = 30000
	c.appsPort = 30001

	c.routerPort = 8080

	c.user = "admin"
	c.pass = "admin"
	c.secret = []byte("secret")

	c.appsRoot = APPS_ROOT
	c.basicServer = BASIC_SERVER
	c.hostname = "127.0.0.1"
	c.container = false
}
