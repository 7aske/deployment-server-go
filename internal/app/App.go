package app

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type App struct {
	Id          string
	Repo        string
	Name        string
	Root        string
	Port        int
	Hostname    string
	Deployed    time.Time
	LastUpdated time.Time
	LastRun     time.Time
	Uptime      string
	Runner      string
	Pid         int
	Process     *os.Process
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
		Id:          a.Id,
		Repo:        a.Repo,
		Name:        a.Name,
		Root:        a.Root,
		Port:        a.Port,
		Hostname:    a.Hostname,
		Deployed:    a.Deployed,
		LastUpdated: a.LastUpdated,
		LastRun:     a.LastRun,
		Uptime:      a.Uptime,
		Runner:      a.Runner,
	})
}

func NewApp(repo string, name string, runner string) *App {
	return &App{Repo: repo, Name: name, Runner: runner}
}
func NewAppFromJson(a *AppJSON) *App {
	return &App{
		Id:          a.Id,
		Repo:        a.Repo,
		Name:        a.Name,
		Root:        a.Root,
		Port:        a.Port,
		Hostname:    a.Hostname,
		Deployed:    a.Deployed,
		LastUpdated: a.LastUpdated,
		LastRun:     a.LastRun,
		Uptime:      a.Uptime,
		Runner:      a.Runner,
	}
}
func (a *App) GetId() string {
	return a.Id
}
func (a *App) SetId(id string) {
	a.Id = id
}
func (a *App) GetRepo() string {
	return a.Repo
}
func (a *App) SetRepo(repo string) {
	a.Repo = repo
}
func (a *App) GetName() string {
	return a.Name
}
func (a *App) SetName(name string) {
	a.Name = name
}
func (a *App) GetRoot() string {
	return a.Root
}
func (a *App) SetRoot(root string) {
	a.Root = root
}
func (a *App) GetPort() int {
	return a.Port
}
func (a *App) SetPort(port int) {
	a.Port = port
}
func (a *App) GetHostname() string {
	return a.Hostname
}
func (a *App) SetHostname(hostname string) {
	a.Hostname = hostname
}
func (a *App) GetDeployed() time.Time {
	return a.Deployed
}
func (a *App) SetDeployed(t time.Time) {
	a.Deployed = t
}
func (a *App) GetLastUpdated() time.Time {
	return a.LastUpdated
}
func (a *App) SetLastUpdated(t time.Time) {
	a.LastUpdated = t
}
func (a *App) GetLastRun() time.Time {
	return a.LastRun
}
func (a *App) SetLastRun(t time.Time) {
	a.LastRun = t
}
func (a *App) GetUptime() string {
	return a.Uptime
}
func (a *App) SetUptime(t string) {
	a.Uptime = t
}
func (a *App) GetPid() int {
	return a.Pid
}
func (a *App) SetPid(p int) {
	a.Pid = p
}
func (a *App) GetProcess() *os.Process {
	return a.Process
}
func (a *App) SetProcess(p *os.Process) {
	a.Process = p
}
func (a *App) GetRunner() string {
	return a.Runner
}
func (a *App) SetRunner(r string) {
	a.Runner = r
}
func (a *App) Print() {
	fmt.Println("-running-")
	fmt.Printf("Id:      \t%s\n", a.Id)
	fmt.Printf("Repo:    \t%s\n", a.Repo)
	fmt.Printf("Name:    \t%s\n", a.Name)
	fmt.Printf("Root:    \t%s\n", a.Root)
	fmt.Printf("Port:    \t%d\n", a.Port)
	//fmt.Printf("Hostname:\t%s\n",a.Hostname)
	//fmt.Printf("Deployed:\t%s\n",a.Deployed)
	//fmt.Printf("LastRun: \t%s\n",a.LastRun)
	//fmt.Printf("Uptime:  \t%s\n",a.Uptime)
	fmt.Printf("Runner:  \t%s\n", a.Runner)
	fmt.Printf("Pid:     \t%d\n", a.Pid)
}
func (a *AppJSON) Print() {
	fmt.Println("-Deployed-")
	fmt.Printf("Id:      \t%s\n", a.Id)
	fmt.Printf("Repo:    \t%s\n", a.Repo)
	fmt.Printf("Name:    \t%s\n", a.Name)
	fmt.Printf("Root:    \t%s\n", a.Root)
	fmt.Printf("Port:    \t%d\n", a.Port)
	//fmt.Printf("Hostname:\t%s\n",a.Hostname)
	//fmt.Printf("Deployed:\t%s\n",a.Deployed)
	//fmt.Printf("LastRun: \t%s\n",a.LastRun)
	//fmt.Printf("Uptime:  \t%s\n",a.Uptime)
	fmt.Printf("Runner:  \t%s\n", a.Runner)
	//fmt.Printf("Pid:     \t%d\n",a.Pid)
}
