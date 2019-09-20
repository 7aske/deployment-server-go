package controllers

import (
	"../app"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const HELP_FORMAT = "%-10s\t%-20s\t%s\n"

type Cli struct {
	deployer *Deployer
}

func NewCli(d *Deployer) *Cli {
	return &Cli{d}
}
func (c *Cli) ParseCommand(args ...string) {
	if len(args) > 0 {
		switch args[0] {
		case "help", "?":
			printHelp()
		case "quit", "exit":
			os.Exit(0)
		case "deploy", "dep":
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
				printHelp()
			} else {
				c.Run(args[1])
			}

		case "find", "ls", "list":
			if len(args) == 1 {
				c.Find("", "")
			} else if len(args) == 2 {
				if args[1] == "dep" {
					c.Find("", "deployed")
				} else if args[1] == "run" {
					c.Find("", "running")
				} else {
					c.Find(args[1], "")
				}
			} else if len(args) == 3 {
				if args[1] == "dep" || args[1] == "run" {
					c.Find(args[2], args[1])
				} else if args[2] == "dep" || args[2] == "run" {
					c.Find(args[1], args[2])
				} else {
					printHelp()
				}
			} else {
				printHelp()
			}
		case "remove", "rm":
			if len(args) == 2 {
				c.Remove(args[1])
			} else {
				printHelp()
			}
		case "kill":
			if len(args) == 2 {
				c.Kill(args[1])
			} else {
				printHelp()
			}
		case "settings":
			if len(args) == 3 {
				c.Settings(args[1], args[2])
			} else {
				printHelp()
			}
		case "config":
			if len(args) == 3 {
				c.Config(args[1], args[2])
			} else {
				printHelp()
			}
		default:
			fmt.Printf("unrecognized command \"%s\"\n", args[0])
			fmt.Println("type \"help\" or \"?\" from more information")
		}
	}
}
func (c *Cli) Deploy(repo string, runner string, hostname string, port int) {
	a, err := c.deployer.Deploy(repo, runner, hostname, port)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.deployer.Install(a)
	if err != nil {
		fmt.Println(err)
		return
	}
	a.Print()
}
func (c *Cli) Run(query string) {
	if appJson, ok := c.deployer.GetAppD(query); ok {
		a := app.NewAppFromJson(appJson)
		if c.deployer.IsAppRunning(a) {
			fmt.Println("already running")
		} else {
			err := c.deployer.Run(a)
			if err != nil {
				fmt.Println(err)
				return
			}
			a.Print()
		}
	} else {
		fmt.Println("not found")
	}
}
func (c *Cli) Kill(query string) {
	if a, ok := c.deployer.GetApp(query); ok {
		_ = c.deployer.Kill(a)
		fmt.Println("killed app with pid " + strconv.Itoa(a.GetPid()))
	} else {
		fmt.Println("not found")
	}
}
func (c *Cli) Remove(query string) {
	if appJson, ok := c.deployer.GetAppD(query); ok {
		if a, ok := c.deployer.GetApp(appJson.Id); ok {
			_ = c.deployer.Kill(a)
		}
		err := c.deployer.Remove(appJson)
		if err != nil {
			fmt.Println("failed to remove app with id of " + appJson.Id)
			fmt.Println(err)
		} else {
			fmt.Println("removed app with id of " + appJson.Id)
		}
	} else {
		fmt.Println("not found")
	}
}

func (c *Cli) Find(query string, typ string) {
	apps := c.deployer.GetApps()
	appsD := c.deployer.GetDeployedApps()
	if typ == "deployed" {
		if query == "" {
			for _, a := range appsD {
				a.Print()
			}
		} else {
			for _, a := range appsD {
				if a.Id == query || a.Name == query || strconv.Itoa(a.Pid) == query {
					a.Print()
					return
				}
			}
		}
	} else if typ == "running" {
		if query == "" {
			for _, a := range *apps {
				a.Print()
			}
		} else {
			for _, a := range *apps {
				if a.Id == query || a.Name == query {
					a.Print()
					return
				}
			}
		}
	} else {
		if query == "" {
			for _, a := range *apps {
				a.Print()
			}
			for _, a := range appsD {
				a.Print()
			}
		} else {
			for _, a := range *apps {
				if a.Id == query || a.Name == query || strconv.Itoa(a.Pid) == query {
					a.Print()
					return
				}
			}
			for _, a := range appsD {
				if a.Id == query || a.Name == query {
					a.Print()
					return
				}
			}
		}
	}

}
func (c *Cli) Settings(query string, setting string) {
	kv := strings.Split(setting, "=")
	if len(kv) != 2 {
		fmt.Printf("invalid settings key-value pair\n")
		return
	}
	settings := make(map[string]string, 1)
	settings[kv[0]] = kv[1]
	err := c.deployer.Settings(query, settings)
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Printf("updated " + query + "\n")
	}
}

func (c *Cli) Config(s string, s2 string) {

}

func printHelp() {
	// deploy run find kill remove settings
	fmt.Print("deployment-server 0.0.1 == Nikola Tasic == github.com/7aske\n\n")
	fmt.Printf(HELP_FORMAT, "deploy", "<repo-url>|<usr> <repo>", "deploy app from specified")
	fmt.Printf(HELP_FORMAT, "", "", "github repository")

	fmt.Printf(HELP_FORMAT, "run", "<app|id>", "run the deployed app")
	fmt.Printf(HELP_FORMAT, "", "", "with specified name or id")

	fmt.Printf(HELP_FORMAT, "find", "[dep|run] [app|id]", "list apps based on search")
	fmt.Printf(HELP_FORMAT, "", "", "terms")

	fmt.Printf(HELP_FORMAT, "kill", "<app|id>", "kill app with specified")
	fmt.Printf(HELP_FORMAT, "", "", "name or id")

	fmt.Printf(HELP_FORMAT, "remove", "<app|id>", "remove app with specified")
	fmt.Printf(HELP_FORMAT, "", "", "name or id")

	fmt.Printf(HELP_FORMAT, "settings", "<app|id> <key=value>", "change the settings of a deployed")
	fmt.Printf(HELP_FORMAT, "", "", "app based on name or id")

	fmt.Printf(HELP_FORMAT, "quit", "", "exits the application")
}
