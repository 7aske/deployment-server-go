package controllers

import (
	"../config"
	"../utils"
	"encoding/json"
	"fmt"
	"github.com/teris-io/shortid"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type Deployer struct {
	config *config.Config
	port   int
	apps   []*App
	appsD  []AppJSON
}
type AppsJSON struct {
	Apps []AppJSON `json:"apps"`
}

//func (d *Deployer) LoadConfig() {
//	d.config = &config.LoadConfig()
//}
func NewDeployer(cfg *config.Config) Deployer {
	d := Deployer{}
	d.config = cfg
	d.port = cfg.GetAppsPort()
	d.apps = []*App{}
	utils.MakeDirIfNotExist(cfg.GetAppsRoot())
	d.initAppsJson()
	d.appsD = d.LoadAppsFromJson().Apps
	return d
}

func (d *Deployer) GetApps() *[]*App {
	return &d.apps
}
func (d *Deployer) GetAppsD() *[]AppJSON {
	d.appsD = d.LoadAppsFromJson().Apps
	return &d.appsD
}
func (d *Deployer) GetAppD(search string) (*AppJSON, bool) {
	for _, a := range d.appsD {
		if a.Name == search || a.Repo == search {
			return &a, true
		}
	}
	return &AppJSON{}, false
}
func (d *Deployer) GetApp(search string) (*App, bool) {
	for _, a := range d.apps {
		if a.GetName() == search || a.GetRepo() == search || a.GetId() == search {
			return a, true
		}
	}
	return &App{}, false
}

// output current running apps as JSON array
func (d *Deployer) GetAppsJSON() []AppJSON {
	apps := d.GetApps()
	var arr []AppJSON
	for _, a := range *apps {
		arr = append(arr, AppJSON{
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
	return arr
}

func (d *Deployer) AddApp(a *App) bool {
	if !d.IsAppRunning(a) {
		d.apps = append(d.apps, a)
		return true
	}
	return false
}
func (d *Deployer) RemoveApp(a *App) {
	newApps := []*App{}
	if d.IsAppRunning(a) {
		for _, app := range d.apps {
			if !(app.name == a.name || app.repo == a.repo || app.id == a.id) {
				newApps = append(newApps, app)
			}
		}
	}
	d.apps = newApps
}
func (d *Deployer) RemoveAppD(a *AppJSON) {
	newApps := []AppJSON{}
	if d.IsAppRunning(NewAppFromJson(a)) {
		for _, app := range d.appsD {
			if !(app.Name == a.Name || app.Repo == a.Repo || app.Id == a.Id) {
				newApps = append(newApps, app)
			}
		}
	}
	d.appsD = newApps
}

func (d *Deployer) Deploy(repo string, runner string) (*App, error) {
	git := exec.Command("git", "-C", d.GetConfig().GetAppsRoot(), "clone", repo)

	_ = git.Run()
	name := utils.GetNameFromRepo(repo)

	app := NewApp(repo, name, runner)
	app.SetId(shortid.MustGenerate())
	app.SetRoot(path.Join(d.GetConfig().GetAppsRoot(), app.name))
	app.SetDeployed(time.Now())
	return app, nil
}

func (d *Deployer) Install(a *App) error {
	npm := exec.Command("npm", "install")
	npm.Dir = a.GetRoot()
	npm.Stdout = os.Stdout
	npm.Stderr = os.Stderr
	err := npm.Run()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func (d *Deployer) Run(a *App) error {
	node := exec.Command("node", d.GetConfig().GetBasicServer())
	node.Dir = a.GetRoot()
	port := a.GetPort()
	if port == 0 {
		a.SetPort(d.generatePort())
	}
	node.Env = append(node.Env, fmt.Sprintf("PORT=%d", port))
	node.Stdout = os.Stdout
	node.Stderr = os.Stderr
	err := node.Start()
	if err != nil {
		fmt.Println(err)
		return err
	}
	a.SetLastRun(time.Now())
	a.SetPid(node.Process.Pid)
	a.SetProcess(node.Process)
	d.AddApp(a)
	d.SaveAppsToJson()
	d.appsD = d.LoadAppsFromJson().Apps
	fmt.Printf("starting server with pid - %d on port %d\n", a.GetPid(), a.GetPort())
	return nil
}
func (d *Deployer) Kill(app *App) {
	err := app.process.Kill()
	if err != nil {
		fmt.Println(err)
	}
	d.RemoveApp(app)
}
func (d *Deployer) Remove(app *AppJSON) {
	a := NewAppFromJson(app)
	if !d.IsAppRunning(a) {
		if strings.HasPrefix(app.Root, path.Join(d.GetConfig().GetAppsRoot())) {
			err := os.RemoveAll(app.Root)
			if err != nil {
				fmt.Println(err)
			}
			d.RemoveAppD(app)
			d.SaveAppsToJson()
		}
	}
}

//check whether the current app is running
func (d *Deployer) IsAppRunning(app *App) bool {
	for _, a := range d.apps {
		if a.name == app.name || a.id == app.id || a.repo == app.repo {
			return true
		}
	}
	return false
}

func (d *Deployer) GetRunningApp(search string) (*App, bool) {
	for _, a := range d.apps {
		if a.name == search || a.id == search || a.repo == search {
			return a, true
		}
	}
	return &App{}, false
}
func (d *Deployer) SetPort(port int) {
	d.port = port
}

func (d *Deployer) GetPort() int {
	return d.port
}
func (d *Deployer) GetConfig() *config.Config {
	return d.config
}
func (d *Deployer) SetConfig(c *config.Config) {
	d.config = c
}

func (d *Deployer) runNode() {

}
func (d *Deployer) runPython() {

}
func (d *Deployer) runPythonFlask() {

}
func (d *Deployer) runWeb() {

}

// load apps from json file into appsD array
func (d *Deployer) LoadAppsFromJson() AppsJSON {
	pth := path.Join(d.GetConfig().GetAppsRoot(), "apps.json")
	jsonFile, _ := ioutil.ReadFile(pth)
	appsD := AppsJSON{}
	err := json.Unmarshal(jsonFile, &appsD)
	if err != nil {
		fmt.Println("load " + err.Error())
	}
	return appsD
}

// save currently running apps to json
func (d *Deployer) SaveAppsToJson() {
	fmt.Println(fmt.Sprintf("saving %d apps", len(d.GetAppsJSON())))
	pth := path.Join(d.GetConfig().GetAppsRoot(), "apps.json")
	apps, _ := json.Marshal(AppsJSON{Apps: d.GetAppsJSON()})
	err := ioutil.WriteFile(pth, apps, 0775)
	if err != nil {
		fmt.Println(err)
	}
}
func (d *Deployer) initAppsJson() {
	pth := path.Join(d.GetConfig().GetAppsRoot(), "apps.json")
	if _, err := os.Stat(pth); err != nil {
		emptyArr, _ := json.Marshal(&AppsJSON{Apps: []AppJSON{}})
		err := ioutil.WriteFile(pth, emptyArr, 0775)
		if err != nil {

			fmt.Println("init " + err.Error())
		}
	}
}
func (d *Deployer) generatePort() int {
	port := d.GetConfig().GetAppsPort()
	for d.isPortUsed(port) {
		port++
	}
	return port
}
func (d *Deployer) isPortUsed(port int) bool{
	for _, a := range d.appsD {
		if port == a.Port {
			return true
		}
	}
	return false
}
