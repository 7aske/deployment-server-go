package controllers
// TODO: in-memory deployed apps
import (
	"../config"
	"../utils"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type Deployer struct {
	config  *config.Config
	port    int
	apps    []*App
	appsD   []AppJSON
	runners []string
}
type AppsJSON struct {
	Apps []AppJSON `json:"apps"`
}
type PackageJSON struct {
	Main string `json:"main"`
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
	d.GetDeployedApps()
	d.runners = []string{"node", "web", "python"}
	//arr := []int{1,2,3}
	//arr = append(arr[:1], arr[2:]...)
	//fmt.Println(arr)
	return d
}

func (d *Deployer) GetApps() *[]*App {
	return &d.apps
}
func (d *Deployer) GetAppsD() *[]AppJSON {
	return &d.appsD
}
func (d *Deployer) GetAppD(search string) (*AppJSON, bool) {
	for _, a := range d.GetDeployedApps() {
		if a.Name == search || a.Repo == search || a.Id == search {
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
func (d *Deployer) GetAppAsJSON(a *App) AppJSON {
	return AppJSON{
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
		Pid:         a.pid,
	}
}

// output current running apps as JSON array
func (d *Deployer) GetAppsAsJSON() []AppJSON {
	var arr []AppJSON
	for _, a := range d.apps {
		app := AppJSON{
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
			Pid:         a.pid,
		}
		app.Uptime = time.Now().Sub(a.lastRun).String()
		arr = append(arr, app)
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
	var newApps []*App
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
	var newApps []AppJSON
	for _, app := range d.appsD {
		if !(app.Name == a.Name || app.Repo == a.Repo || app.Id == a.Id) {
			newApps = append(newApps, app)
		}
	}
	d.appsD = newApps
}

func (d *Deployer) Deploy(repo string, runner string, hostname string, port int) (*App, error) {
	if utils.Contains(runner, &d.runners) == -1 {
		return &App{}, errors.New("unknown runner")
	}
	name := utils.GetNameFromRepo(repo)
	app := NewApp(repo, name, runner)
	app.SetRoot(path.Join(d.GetConfig().GetAppsRoot(), app.GetName()))
	if !d.isPortUsed(port) {
		app.SetPort(port)
	}
	app.SetHostname(hostname)
	//fmt.Println(app.GetRoot())
	//fmt.Println(app.GetRunner())
	if _, ok := d.GetAppD(repo); !ok {
		git := exec.Command("git", "-C", d.GetConfig().GetAppsRoot(), "clone", repo)
		git.Stdout = os.Stdout
		git.Stderr = os.Stderr
		err := git.Run()
		if err != nil {
			fmt.Println(err)
			return app, err
		}
		app.SetId(shortid.MustGenerate())
		app.SetDeployed(time.Now())
		d.SaveAppToJson(d.GetAppAsJSON(app))
		return app, nil
	} else {
		return &App{name: name, repo: repo, runner: runner}, errors.New("app already deployed")
	}

}
func (d *Deployer) Update(a *App) error {
	if d.IsAppRunning(a) {
		_ = d.Kill(a)
	}
	git := exec.Command("git", "-C", a.GetRoot(), "pull")
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	err := git.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}
	a.SetLastUpdated(time.Now())
	//err = d.Run(a)
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	d.SaveAppToJson(d.GetAppAsJSON(a))
	return nil
}
func (d *Deployer) Install(a *App) error {
	if utils.Contains(a.GetRunner(), &d.runners) == -1 {
		return errors.New("unknown runner")
	}
	switch a.GetRunner() {
	case "node":
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
	case "web":
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
	default:
		return errors.New("unsupported runner")
	}

}

func (d *Deployer) Run(a *App) error {
	if utils.Contains(a.GetRunner(), &d.runners) == -1 {
		return errors.New("unknown runner")
	}
	switch a.GetRunner() {
	case "node":
		return d.runNode(a)
	case "web":
		return d.runWeb(a)
	case "python":
		return d.runPython(a)
	default:
		return errors.New("unsupported runner")
	}

}
func (d *Deployer) Kill(app *App) error {
	err := app.GetProcess().Kill()
	if err != nil {
		fmt.Println(err)
		return err
	}
	jApp := d.GetAppAsJSON(app)
	jApp.Pid = 0
	d.SaveAppToJson(jApp)
	d.RemoveApp(app)
	return nil
}
func (d *Deployer) Remove(app *AppJSON) {
	a := NewAppFromJson(app)
	if !d.IsAppRunning(a) {
		if strings.HasPrefix(app.Root, path.Join(d.GetConfig().GetAppsRoot())) {
			err := os.RemoveAll(app.Root)
			if err != nil {
				fmt.Println(err)
			}
			d.RemoveAppFromJson(*app)
		}
	}
}
func (d *Deployer) Settings(id string, settings map[string]string) error {
	if app, ok := d.GetApp(id); ok {
		if d.IsAppRunning(app) {
			return errors.New("app is running")
		}
		return errors.New("app in memory")
	}
	if appJson, ok := d.GetAppD(id); ok {
		changed := false
		for key, value := range settings {
			switch key {
			case "port":
				port, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				if d.isPortUsed(port) {
					return errors.New("port in use")
				}
				appJson.Port = port
				changed = true
			case "hostname":
				appJson.Hostname = value
				changed = true
			case "runner":
				if utils.Contains(value, &d.runners) != -1 {
					appJson.Runner = value
					changed = true
				}
			}
		}
		if changed {
			d.SaveAppToJson(*appJson)
		}
		return nil
	} else {
		return errors.New("app not found")
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
func (d *Deployer) runNode(a *App) error {
	packageJSONFile, _ := ioutil.ReadFile(path.Join(a.GetRoot(), "package.json"))
	packageJson := PackageJSON{}
	_ = json.Unmarshal(packageJSONFile, &packageJson)
	node := exec.Command("node", path.Join(a.GetRoot(), packageJson.Main))
	node.Dir = a.GetRoot()
	port := a.GetPort()
	if port == 0 {
		a.SetPort(d.generatePort())
		port = a.GetPort()
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
	d.SaveAppToJson(d.GetAppAsJSON(a))
	fmt.Printf("starting server with pid - %d on port %d\n", a.GetPid(), a.GetPort())
	return nil
}
func (d *Deployer) runWeb(a *App) error {
	node := exec.Command("node", d.GetConfig().GetBasicServer())
	node.Dir = a.GetRoot()
	port := a.GetPort()
	if port == 0 {
		a.SetPort(d.generatePort())
		port = a.GetPort()
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
	d.SaveAppToJson(d.GetAppAsJSON(a))
	fmt.Printf("starting server with pid - %d on port %d\n", a.GetPid(), a.GetPort())
	return nil
}
func (d *Deployer) runPython(a *App) error {
	python := exec.Command("python3", "__main__.py")
	python.Dir = a.GetRoot()
	port := a.GetPort()
	if port == 0 {
		a.SetPort(d.generatePort())
		port = a.GetPort()
	}
	python.Env = append(python.Env, fmt.Sprintf("PORT=%d", port))
	python.Stdout = os.Stdout
	python.Stderr = os.Stderr
	err := python.Start()
	if err != nil {
		fmt.Println(err)
		return err
	}
	a.SetLastRun(time.Now())
	a.SetPid(python.Process.Pid)
	a.SetProcess(python.Process)
	d.AddApp(a)
	d.SaveAppToJson(d.GetAppAsJSON(a))
	fmt.Printf("starting server with pid - %d on port %d\n", a.GetPid(), a.GetPort())
	return nil
}
func (d *Deployer) runPythonFlask() {
}

// load apps from json file into appsD array
func (d *Deployer) GetDeployedApps() []AppJSON {
	pth := path.Join(d.GetConfig().GetAppsRoot(), "apps.json")
	jsonFile, _ := ioutil.ReadFile(pth)
	appsD := AppsJSON{}
	err := json.Unmarshal(jsonFile, &appsD)
	if err != nil {
		fmt.Println("json - " + err.Error())
	}
	return appsD.Apps
}

// save currently running apps to json
func (d *Deployer) RemoveAppFromJson(app AppJSON) {
	pth := path.Join(d.GetConfig().GetAppsRoot(), "apps.json")
	appsJson := d.GetDeployedApps()
	if pos := d.ContainsAppJSON(&appsJson, &app); pos != -1 {
		appsJson = append(appsJson[:pos], appsJson[pos+1:]...)
	}
	apps, _ := json.Marshal(AppsJSON{Apps: appsJson})
	err := ioutil.WriteFile(pth, apps, 0775)
	if err != nil {
		fmt.Println(err)
	}
}
func (d *Deployer) SaveAppToJson(app AppJSON) {
	pth := path.Join(d.GetConfig().GetAppsRoot(), "apps.json")
	appsJson := d.GetDeployedApps()
	app.Pid = 0
	app.Uptime = ""
	if pos := d.ContainsAppJSON(&appsJson, &app); pos == -1 {
		appsJson = append(appsJson, app)
	} else {
		appsJson = append(appsJson[:pos], appsJson[pos+1:]...)
		appsJson = append(appsJson, app)
	}
	apps, _ := json.Marshal(AppsJSON{Apps: appsJson})
	err := ioutil.WriteFile(pth, apps, 0775)
	if err != nil {
		fmt.Println(err)
	}
}
func (d *Deployer) initAppsJson() {
	pth := path.Join(d.GetConfig().GetAppsRoot(), "apps.json")
	if _, err := os.Stat(pth); err != nil {
		fmt.Println("initializing apps folder")
		emptyArr, _ := json.Marshal(&AppsJSON{Apps: []AppJSON{}})
		err := ioutil.WriteFile(pth, emptyArr, 0775)
		if err != nil {
			fmt.Println("json init - " + err.Error())
		}
	} else {
		folders, _ := ioutil.ReadDir(d.GetConfig().GetAppsRoot())
		//for _, f := range folders {
		//	fmt.Println(f.Name())
		//}
		appsD := d.GetDeployedApps()
		for _, app := range appsD {
			if !utils.ContainsFile(app.Name, &folders) {
				fmt.Println(app.Name)
				d.RemoveAppFromJson(app)
			}
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
func (d *Deployer) isPortUsed(port int) bool {
	if port < 1024 {
		return true
	}
	for _, a := range d.GetDeployedApps() {
		if port == a.Port {
			return true
		}
	}
	return false
}
func (d *Deployer) ContainsAppJSON(arr *[]AppJSON, app *AppJSON) int {
	for i, a := range *arr {
		if a.Id == app.Id || a.Name == app.Name || a.Repo == app.Repo {
			return i
		}
	}
	return -1
}
