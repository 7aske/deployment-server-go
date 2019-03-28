package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type App struct {
	id          string
	repo        string
	name        string
	root        string
	port        int
	hostname    string
	deployed    time.Time
	lastUpdated time.Time
	lastRun     time.Time
	uptime      string
	runner      string
	pid         int
	process     *os.Process
}
type AppJSON struct {
	Id          string    `json:"id"`
	Repo        string    `json:"repo"`
	Name        string    `json:"name"`
	Root        string    `json:"root"`
	Port        int       `json:"port"`
	Hostname    string    `json:"hostname"`
	Deployed    time.Time `json:"deployed"`
	LastUpdated time.Time `json:"last_updated"`
	LastRun     time.Time `json:"last_run"`
	Uptime      string    `json:"uptime"`
	Runner      string    `json:"runner"`
	Pid         int       `json:"pid"`
}

func (a *App) GetJSON() ([]byte, error) {
	return json.Marshal(AppJSON{
		Id:          a.id,
		Repo:        a.repo,
		Name:        a.name,
		Root:        a.root,
		Port:        a.port,
		Hostname:    a.hostname,
		Deployed:    a.deployed,
		LastUpdated: a.lastUpdated,
		LastRun:     a.lastRun,
		Uptime:      a.uptime,
		Runner:      a.runner,
	})
}

//return json.Marshal(map[string]interface{}{
//"some_field": w.SomeField,
//})

func NewApp(repo string, name string, runner string) *App {
	return &App{repo: repo, name: name, runner: runner}
}
func NewAppFromJson(a *AppJSON) *App {
	return &App{
		id:          a.Id,
		repo:        a.Repo,
		name:        a.Name,
		root:        a.Root,
		port:        a.Port,
		hostname:    a.Hostname,
		deployed:    a.Deployed,
		lastUpdated: a.LastUpdated,
		lastRun:     a.LastRun,
		uptime:      a.Uptime,
		runner:      a.Runner,
	}
}
func (a *App) GetId() string {
	return a.id
}
func (a *App) SetId(id string) {
	a.id = id
}
func (a *App) GetRepo() string {
	return a.repo
}
func (a *App) SetRepo(repo string) {
	a.repo = repo
}
func (a *App) GetName() string {
	return a.name
}
func (a *App) SetName(name string) {
	a.name = name
}
func (a *App) GetRoot() string {
	return a.root
}
func (a *App) SetRoot(root string) {
	a.root = root
}
func (a *App) GetPort() int {
	return a.port
}
func (a *App) SetPort(port int) {
	a.port = port
}
func (a *App) GetHostname() string {
	return a.hostname
}
func (a *App) SetHostname(hostname string) {
	a.hostname = hostname
}
func (a *App) GetDeployed() time.Time {
	return a.deployed
}
func (a *App) SetDeployed(t time.Time) {
	a.deployed = t
}
func (a *App) GetLastUpdated() time.Time {
	return a.lastUpdated
}
func (a *App) SetLastUpdated(t time.Time) {
	a.lastUpdated = t
}
func (a *App) GetLastRun() time.Time {
	return a.lastRun
}
func (a *App) SetLastRun(t time.Time) {
	a.lastRun = t
}
func (a *App) GetUptime() string {
	return a.uptime
}
func (a *App) SetUptime(t string) {
	a.uptime = t
}
func (a *App) GetPid() int {
	return a.pid
}
func (a *App) SetPid(p int) {
	a.pid = p
}
func (a *App) GetProcess() *os.Process {
	return a.process
}
func (a *App) SetProcess(p *os.Process) {
	a.process = p
}
func (a *App) GetRunner() string {
	return a.runner
}
func (a *App) SetRunner(r string) {
	a.runner = r
}
func (a *App) Print() {
	fmt.Println("-running-")
	fmt.Printf("id:      \t%s\n", a.id)
	fmt.Printf("repo:    \t%s\n", a.repo)
	fmt.Printf("name:    \t%s\n", a.name)
	fmt.Printf("root:    \t%s\n", a.root)
	fmt.Printf("port:    \t%d\n", a.port)
	//fmt.Printf("hostname:\t%s\n",a.hostname)
	//fmt.Printf("deployed:\t%s\n",a.deployed)
	//fmt.Printf("lastRun: \t%s\n",a.lastRun)
	//fmt.Printf("uptime:  \t%s\n",a.uptime)
	fmt.Printf("runner:  \t%s\n", a.runner)
	fmt.Printf("pid:     \t%d\n", a.pid)
}
func (a *AppJSON) Print() {
	fmt.Println("-deployed-")
	fmt.Printf("id:      \t%s\n", a.Id)
	fmt.Printf("repo:    \t%s\n", a.Repo)
	fmt.Printf("name:    \t%s\n", a.Name)
	fmt.Printf("root:    \t%s\n", a.Root)
	fmt.Printf("port:    \t%d\n", a.Port)
	//fmt.Printf("hostname:\t%s\n",a.Hostname)
	//fmt.Printf("deployed:\t%s\n",a.Deployed)
	//fmt.Printf("lastRun: \t%s\n",a.LastRun)
	//fmt.Printf("uptime:  \t%s\n",a.Uptime)
	fmt.Printf("runner:  \t%s\n", a.Runner)
	//fmt.Printf("pid:     \t%d\n",a.Pid)
}
