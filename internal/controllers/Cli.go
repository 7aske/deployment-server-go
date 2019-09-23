package controllers

import (
	"../app"
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strconv"
	"strings"
)

var cmdList []string
var cmdListIdx = 0
var prompt = "\r-> "

const HELP_FORMAT = "%-10s\t%-20s\t%s\r\n"

type Cli struct {
	deployer *Deployer
	running  bool
	state    *terminal.State
}

func NewCli(d *Deployer) *Cli {
	return &Cli{deployer: d, running: true}
}

func (c *Cli) Release() {
	c.running = false
	if err := terminal.Restore(0, c.state); err != nil {
		log.Println("warning: failed to restore terminal:", err)
	}
}

func (c *Cli) Start() {
	state, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalln("setting stdin to raw: ", err)
	}
	c.state = state
	defer c.Release()

	in := bufio.NewReader(os.Stdin)

	for c.running {
		line := ""
		fmt.Print(prompt)
		for c.running {
			r, _, err := in.ReadRune()
			if err != nil {
				log.Println("stdin: ", err)
				break
			}
			if r == 13 || r == 10 {
				fmt.Print("\r\n")
				break
			} else {
				line += string(r)
			}
			if len(line) >= 2 && strings.Index(line, "\033[") != -1 {
				switch line[len(line)-1] {
				case 'A':
					if len(cmdList) > 0 && cmdListIdx >= 0 {
						fmt.Print("\r\033[K")
						line = cmdList[cmdListIdx]
						if cmdListIdx > 0 {
							cmdListIdx--
						}
						strings.Trim(line, "\033[A")
						fmt.Print(line)
					}
				case 'B':
					if cmdListIdx < len(cmdList) {
						fmt.Print("\r\033[K")
						line = cmdList[cmdListIdx]
						if cmdListIdx < len(cmdList)-1 {
							cmdListIdx++
						} else {
							line = ""
						}
						strings.Trim(line, "\033[B")
						fmt.Print(line)
					}
				}

			} else {
				fmt.Printf("%c", r)
			}

			switch r {
			case 3:
				return
			case 12:
				line = line[:len(line)-1]
				fmt.Print("\033[2J\033[2H\033[1A")
				fmt.Print(prompt)
			case 127:
				if len(line) >= 2 {
					line = line[:len(line)-2]
					fmt.Print("\r\033[K")
					fmt.Print(prompt)
					fmt.Print(line)
				} else {
					line = ""
				}
			}
		}
		if len(line) > 0 {
			cmdList = append(cmdList, line)
			cmdListIdx = len(cmdList) - 1
		}
		c.ParseCommand(strings.Split(line, " ")...)
		fmt.Print(prompt)
		line = ""
	}
}

func (c *Cli) ParseCommand(args ...string) {
	if len(args) > 0 {
		switch args[0] {
		case "help", "?":
			printHelp()
		case "pid":
			fmt.Println(strconv.Itoa(os.Getpid()) + "\r\n")
		case "clear", "cls":
			fmt.Printf("\033[2J\033[2H\033[1A")
		case "quit", "exit", "q":
			c.running = false
		case "deploy", "dep":
			if len(args) < 3 {
				fmt.Print("deploy <repo> <runner>\r\n")
			} else if len(args) == 6 {
				port, _ := strconv.Atoi(args[5])
				c.Deploy(fmt.Sprintf("https://github.com/%s/%s", args[1], args[2]), args[3], args[4], port)
			} else {
				c.Deploy(args[1], args[2], "", 0)
			}
		case "run", "start":
			if len(args) < 2 {
				printHelp()
			} else {
				c.Run(args[1])
			}
		case "update":
			if len(args) < 2 {
				printHelp()
			} else {
				c.Update(args[1])
			}
		case "find", "ls", "list":
			if len(args) == 1 {
				c.Find("", "")
			} else if len(args) == 2 {
				if args[1] == "dep" || args[1] == "deployed" {
					c.Find("", "deployed")
				} else if args[1] == "run" || args[1] == "running" {
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
			fmt.Printf("unrecognized command \"%s\"\r\n", args[0])
			fmt.Print("type \"help\" or \"?\" from more information\r\n")
		}
	}
}
func (c *Cli) Deploy(repo string, runner string, hostname string, port int) {
	a, err := c.deployer.Deploy(repo, runner, hostname, port)
	if err != nil {
		fmt.Print("\r")
		fmt.Println(err)
		return
	}
	err = c.deployer.Install(a)
	if err != nil {
		fmt.Print("\r")
		fmt.Println(err)
		return
	}
	a.Print()
}
func (c *Cli) Run(query string) {
	if appJson, ok := c.deployer.GetAppD(query); ok {
		a := app.NewAppFromJson(appJson)
		if c.deployer.IsAppRunning(a) {
			fmt.Println("already running\r")
		} else {
			err := c.deployer.Run(a)
			if err != nil {
				fmt.Print("\r")
				fmt.Println(err)
				return
			}
			a.Print()
		}
	} else {
		fmt.Println("not found\r")
	}
}
func (c *Cli) Update(query string) {
	if appJson, ok := c.deployer.GetAppD(query); ok {
		a := app.NewAppFromJson(appJson)
		if c.deployer.IsAppRunning(a) {
			fmt.Println("already running\r")
		} else {
			err := c.deployer.Update(a)
			if err != nil {
				fmt.Print("\r")
				fmt.Println(err)
				return
			}
		}
		fmt.Printf("updated - %s from %s\r\n", a.GetName(), a.GetRepo())
	} else {
		fmt.Println("not found\r")
	}
}
func (c *Cli) Kill(query string) {
	if a, ok := c.deployer.GetApp(query); ok {
		_ = c.deployer.Kill(a)
		fmt.Println("killed app with pid " + strconv.Itoa(a.GetPid()) + "\r")
	} else {
		fmt.Println("not found\r")
	}
}
func (c *Cli) Remove(query string) {
	if appJson, ok := c.deployer.GetAppD(query); ok {
		if a, ok := c.deployer.GetApp(appJson.Id); ok {
			_ = c.deployer.Kill(a)
		}
		err := c.deployer.Remove(appJson)
		if err != nil {
			fmt.Println("failed to remove app with id of " + appJson.Id + "\r")
			fmt.Println(err)
		} else {
			fmt.Println("removed app with id of " + appJson.Id + "\r")
		}
	} else {
		fmt.Println("not found\r")
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
		fmt.Printf("invalid settings key-value pair\r\n")
		return
	}
	settings := make(map[string]string, 1)
	settings[kv[0]] = kv[1]
	err := c.deployer.Settings(query, settings)
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Printf("updated " + query + "\r\n")
	}
}
func (c *Cli) Config(key string, value string) {

}
func printHelp() {
	fmt.Print("deployment-server 0.0.1 == Nikola Tasic == github.com/7aske\r\n\r\n")
	fmt.Printf(HELP_FORMAT, "deploy, dep", "<repo-url>", "deploy app from specified")
	fmt.Printf(HELP_FORMAT, "", "", "github repository")
	fmt.Printf(HELP_FORMAT, "run", "<app|id>", "run the deployed app")
	fmt.Printf(HELP_FORMAT, "", "", "with specified name or id")
	fmt.Printf(HELP_FORMAT, "start", "", "run alias")
	fmt.Printf(HELP_FORMAT, "find", "[dep|run] [app|id]", "list apps based on search")
	fmt.Printf(HELP_FORMAT, "", "", "terms")
	fmt.Printf(HELP_FORMAT, "update", "<app|id>", "update app with specified")
	fmt.Printf(HELP_FORMAT, "", "", "name or id")
	fmt.Printf(HELP_FORMAT, "kill", "<app|id>", "kill app with specified")
	fmt.Printf(HELP_FORMAT, "", "", "name or id")
	fmt.Printf(HELP_FORMAT, "remove, rm", "<app|id>", "remove app with specified")
	fmt.Printf(HELP_FORMAT, "", "", "name or id")
	fmt.Printf(HELP_FORMAT, "settings", "<app|id> <key=value>", "change the settings of a deployed")
	fmt.Printf(HELP_FORMAT, "", "", "app based on name or id")
	fmt.Printf(HELP_FORMAT, "quit, q", "", "exits the application")
	fmt.Printf(HELP_FORMAT, "exit", "", "quit alias")
	fmt.Printf(HELP_FORMAT, "clear, cls", "", "clears the screen")
	fmt.Printf(HELP_FORMAT, "pid", "", "returns the deployer pid")
}
