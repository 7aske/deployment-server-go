package controllers

import (
	"../config"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

type Deployer struct {
	config   config.Config
	port     int
	children []Child
}

type Child struct {
	Repo string
}

func (d *Deployer) LoadConfig() {
	d.config = config.LoadConfig()
}

func (d *Deployer) GetConfig() *config.Config {
	return &d.config
}

func New(startingPort int) *Deployer {
	dep := &Deployer{}
	dep.SetPort(startingPort)
	dep.children = []Child{}
	return dep
}
func (d *Deployer) GetChildren() *[]Child {
	return &d.children
}

func (d *Deployer) SetPort(port int) {
	d.port = port
}

func (d *Deployer) GetPort() int {
	return d.port
}

func (d *Deployer) Deploy(repo string) (*Child, error) {
	cwd, _ := os.Getwd()
	clone := exec.Command("git", "-C", cwd, "clone", repo)
	err := clone.Run()
	child := Child{Repo: repo}
	if err != nil {
		return &Child{}, errors.New(fmt.Sprintf("unable to clone repo %s", repo))
	} else {
		d.children = append(d.children, child)
		return &child, nil
	}
}
