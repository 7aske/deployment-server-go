package controllers

import (
	"encoding/json"
	"os"
	"time"
)

type App struct {
	id       string
	repo     string
	name     string
	root     string
	port     int
	hostname string
	deployed time.Time
	lastRun  time.Time
	uptime   time.Time
	runner   string
	pid      int
	process  *os.Process
}
type AppJSON struct {
	Id       string    `json:"id"`
	Repo     string    `json:"repo"`
	Name     string    `json:"name"`
	Root     string    `json:"root"`
	Port     int       `json:"port"`
	Hostname string    `json:"hostname"`
	Deployed time.Time `json:"deployed"`
	LastRun  time.Time `json:"last_run"`
	Uptime   time.Time `json:"uptime"`
	Runner   string    `json:"runner"`
}

func (a *App) GetJSON() ([]byte, error) {
	return json.Marshal(AppJSON{
		Id:       a.id,
		Repo:     a.repo,
		Name:     a.name,
		Root:     a.root,
		Port:     a.port,
		Hostname: a.hostname,
		Deployed: a.deployed,
		LastRun:  a.lastRun,
		Uptime:   a.uptime,
		Runner:   a.runner,
	})
}

//return json.Marshal(map[string]interface{}{
//"some_field": w.SomeField,
//})

func NewApp(repo string, name string, runner string) *App {
	a := App{repo: repo, name: name, runner: runner}
	return &a
}
func NewAppFromJson(a *AppJSON) *App {
	return &App{
		id:       a.Id,
		repo:     a.Repo,
		name:     a.Name,
		root:     a.Root,
		port:     a.Port,
		hostname: a.Hostname,
		deployed: a.Deployed,
		lastRun:  a.LastRun,
		uptime:   a.Uptime,
		runner:   a.Runner,
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
func (a *App) GetLastRun() time.Time {
	return a.lastRun
}
func (a *App) SetLastRun(t time.Time) {
	a.lastRun = t
}
func (a *App) GetUptime() time.Time {
	return a.uptime
}
func (a *App) SetUptime(t time.Time) {
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
