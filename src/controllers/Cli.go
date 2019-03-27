package controllers

import (
	"fmt"
	"strconv"
)

type Cli struct {
	deployer    *Deployer
	lastCommand []byte
}

func NewCli(d *Deployer) *Cli {
	return &Cli{d, []byte{}}
}
func (c *Cli) ParseCommand(args ...string) {
	if len(args) > 0 {
		switch args[0] {
		case "help":
			printHelp()
		case "deploy":
			if len(args) < 3 {
				fmt.Println("deploy <repo> <runner>")
			} else if len(args) == 6 {
				port, _ := strconv.Atoi(args[5])
				c.Deploy(fmt.Sprintf("https://github.com/%s/%s", args[1], args[2]), args[3], args[4], port)
			} else {
				c.Deploy(args[1], args[2], "", 0)
			}
		case "run":
			if len(args) < 2 {
				fmt.Println("run <query>")
			} else {
				c.Run(args[1])
			}
		case "find":
			if len(args) == 1 {
				c.Find("")
			} else if len(args) == 2 {
				c.Find(args[1])
			} else {
				fmt.Println("find <query>")
			}
		case "remove":
			if len(args) == 2 {
				c.Remove(args[1])
			} else {
				fmt.Println("remove <query>")
			}
		case "kill":
			if len(args) == 2 {
				c.Kill(args[1])
			} else {
				fmt.Println("kill <query>")
			}
		}

	}
}
func (c *Cli) Deploy(repo string, runner string, hostname string, port int) {
	app, err := c.deployer.Deploy(repo, runner, hostname, port)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.deployer.Install(app)
	if err != nil {
		fmt.Println(err)
		return
	}
	app.Print()
}
func (c *Cli) Run(query string) {
	if appJson, ok := c.deployer.GetAppD(query); ok {
		app := NewAppFromJson(appJson)
		if c.deployer.IsAppRunning(app) {
			fmt.Println("already running")
		} else {
			err := c.deployer.Run(app)
			if err != nil {
				fmt.Println(err)
				return
			}
			app.Print()
		}
	} else {
		fmt.Println("not found")
	}
}
func (c *Cli) Kill(query string) {
	if app, ok := c.deployer.GetApp(query); ok {
		_ = c.deployer.Kill(app)
		fmt.Println("killed app with pid " + strconv.Itoa(app.GetPid()))
	} else {
		fmt.Println("not found")
	}
}
func (c *Cli) Remove(query string) {
	if appJson, ok := c.deployer.GetAppD(query); ok {
		if app, ok := c.deployer.GetApp(appJson.Id); ok {
			_ = c.deployer.Kill(app)
		}
		c.deployer.Remove(appJson)
		fmt.Println("removed app with id of " + appJson.Id)
	} else {
		fmt.Println("not found")
	}
}

func (c *Cli) Find(query string) {
	apps := c.deployer.GetApps()
	appsD := c.deployer.GetDeployedApps()
	if query == "" {
		for _, a := range *apps {
			a.Print()
		}
		for _, a := range appsD {
			a.Print()
		}
	} else {
		for _, a := range *apps {
			if a.id == query || a.name == query || strconv.Itoa(a.pid) == query {
				a.Print()
			}
		}
		for _, a := range appsD {
			if a.Id == query || a.Name == query {
				a.Print()
			}
		}
	}
}

//func (c *Cli) PutLastCommand() {
//fmt.Println(string(c.lastCommand))
//	buffer := bytes.Buffer{}
//	buffer.Write(c.lastCommand)
//	os.Stdin = buffer
//}
func printHelp() {
	fmt.Println("printing help")
}
