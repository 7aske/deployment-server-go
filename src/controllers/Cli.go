package controllers

import (
	"fmt"
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
			} else {
				c.Deploy(args[1], args[2])
			}
		case "run":
			if len(args) < 2 {
				fmt.Println("run <query>")
			} else {
				c.Run(args[1])
			}
		}
	}
}
func (c *Cli) Deploy(repo string, runner string) {
	app, _ := c.deployer.Deploy(repo, runner)
	_ = c.deployer.Install(app)
	app.Print()
}
func (c *Cli) Run(query string) {
	if appJson, ok := c.deployer.GetAppD(query); ok {
		app := NewAppFromJson(appJson)
		if c.deployer.IsAppRunning(app) {
			fmt.Println("already running")
		} else {
			err := c.deployer.Run(app)
			app.Print()
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("not found")
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
