package controllers

// TODO: in-memory deployed apps
import (
	"../app"
	"../config"
	"../logger"
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
	"syscall"
	"time"
)

type Deployer struct {
	config  *config.Config
	port    int
	apps    []*app.App
	appsD   []app.AppJSON
	runners []string
	logger  *logger.Logger
}
type AppsJSON struct {
	Apps []app.AppJSON `json:"apps"`
}
type PackageJSON struct {
	Main string `json:"main"`
}

func New(cfg *config.Config) Deployer {
	d := Deployer{}
	d.config = cfg
	d.port = cfg.GetAppsPort()
	d.apps = []*app.App{}
	utils.MakeDirIfNotExist(path.Join(cfg.GetCwd(), cfg.GetAppsRoot()))
	d.initAppsJson()
	d.GetDeployedApps()
	d.runners = []string{"node", "web", "python", "flask"}
	d.logger = logger.NewLogger(logger.LOG_DEPLOYER)
	return d
}
func (d *Deployer) SetLogger(l *logger.Logger) {
	d.logger = l
}

func (d *Deployer) GetLogger() *logger.Logger {
	return d.logger
}
func (d *Deployer) GetApps() *[]*app.App {
	return &d.apps
}
func (d *Deployer) GetAppsD() *[]app.AppJSON {
	return &d.appsD
}
func (d *Deployer) GetAppD(search string) (*app.AppJSON, bool) {
	for _, a := range d.GetDeployedApps() {
		if a.Name == search || a.Repo == search || a.Id == search {
			return &a, true
		}
	}
	return &app.AppJSON{}, false
}
func (d *Deployer) GetApp(search string) (*app.App, bool) {
	for _, a := range d.apps {
		if a.GetName() == search || a.GetRepo() == search || a.GetId() == search {
			return a, true
		}
	}
	return &app.App{}, false
}
func (d *Deployer) GetAppAsJSON(a *app.App) app.AppJSON {
	return app.AppJSON{
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
		Pid:         a.Pid,
	}
}

// output current running apps as JSON array
func (d *Deployer) GetAppsAsJSON() []app.AppJSON {
	var arr []app.AppJSON
	for _, a := range d.apps {
		a := app.AppJSON{
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
			Pid:         a.Pid,
		}
		a.Uptime = time.Now().Sub(a.LastRun).String()
		arr = append(arr, a)
	}
	return arr
}

func (d *Deployer) AddApp(a *app.App) bool {
	if !d.IsAppRunning(a) {
		d.apps = append(d.apps, a)
		return true
	}
	return false
}
func (d *Deployer) RemoveApp(appToRemove *app.App) {
	var newApps []*app.App
	if d.IsAppRunning(appToRemove) {
		for _, a := range d.apps {
			if !(a.Name == appToRemove.Name || a.Repo == appToRemove.Repo || a.Id == appToRemove.Id) {
				newApps = append(newApps, a)
			}
		}
	}
	d.apps = newApps
}
func (d *Deployer) RemoveAppD(appDToRemove *app.AppJSON) {
	var newApps []app.AppJSON
	for _, a := range d.appsD {
		if !(a.Name == appDToRemove.Name || a.Repo == appDToRemove.Repo || a.Id == appDToRemove.Id) {
			newApps = append(newApps, a)
		}
	}
	d.appsD = newApps
}

func (d *Deployer) Deploy(repo string, runner string, hostname string, port int) (*app.App, error) {
	if utils.Contains(runner, &d.runners) == -1 {
		err := errors.New("deploy - unknown runner " + runner)
		d.logger.Log(err.Error())
		return &app.App{}, err
	}
	name := utils.GetNameFromRepo(repo)
	a := app.NewApp(repo, name, runner)
	a.SetRoot(path.Join(d.GetConfig().GetCwd(), d.GetConfig().GetAppsRoot(), a.GetName()))
	if !d.isPortUsed(port) {
		a.SetPort(port)
	}
	a.SetHostname(hostname)
	if _, ok := d.GetAppD(repo); !ok {
		git := exec.Command("git", "-C", path.Join(d.GetConfig().GetCwd(), d.GetConfig().GetAppsRoot()), "clone", repo)
		git.Stdin = nil
		git.Stdout = os.Stdout
		git.Stderr = os.Stderr
		err := git.Run()
		if err != nil {
			d.logger.Log(err.Error())
			return a, err
		}
		a.SetId(shortid.MustGenerate())
		a.SetDeployed(time.Now())
		d.SaveAppToJson(d.GetAppAsJSON(a))
		d.logger.Log("deploy - deployed app " + name)
		return a, nil
	} else {
		err := errors.New("deploy - app already deployed " + name)
		d.logger.Log(err.Error())
		return &app.App{Name: name, Repo: repo, Runner: runner}, err
	}

}
func (d *Deployer) Update(appToUpdate *app.App) error {
	if d.IsAppRunning(appToUpdate) {
		err := d.Kill(appToUpdate)
		if err != nil {
			d.logger.Log(err.Error())
		}
		d.logger.Log("update - killing running app " + appToUpdate.GetName())
	}
	git := exec.Command("git", "-C", appToUpdate.GetRoot(), "pull")
	git.Stdin = nil
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	err := git.Run()
	if err != nil {
		d.logger.Log(err.Error())
		return err
	}
	appToUpdate.SetLastUpdated(time.Now())
	d.SaveAppToJson(d.GetAppAsJSON(appToUpdate))
	d.logger.Log("update - updated app " + appToUpdate.GetName())
	return nil
}
func (d *Deployer) Install(appToInstall *app.App) error {
	if utils.Contains(appToInstall.GetRunner(), &d.runners) == -1 {
		err := errors.New("install - unknown runner " + appToInstall.GetRunner())
		d.logger.Log(err.Error())
		return err
	}
	switch appToInstall.GetRunner() {
	case "node":
		npm := exec.Command("npm", "install")
		npm.Dir = appToInstall.GetRoot()
		npm.Stdout = os.Stdout
		npm.Stderr = os.Stderr
		err := npm.Run()
		if err != nil {
			d.logger.Log(err.Error())
			return err
		} else {
			d.logger.Log("install - npm install finished " + appToInstall.GetName())
			return nil
		}
	case "web":
		npm := exec.Command("npm", "install")
		npm.Dir = appToInstall.GetRoot()
		npm.Stdout = os.Stdout
		npm.Stderr = os.Stderr
		err := npm.Run()
		if err != nil {
			d.logger.Log(err.Error())
			return err
		} else {
			d.logger.Log("install - npm install finished " + appToInstall.GetName())
			return nil
		}
	case "python":
		return nil
	case "flask":
		return nil

	default:
		err := errors.New("install - unknown runner " + appToInstall.GetRunner())
		d.logger.Log(err.Error())
		return err
	}

}

func (d *Deployer) Run(appToRun *app.App) error {
	if utils.Contains(appToRun.GetRunner(), &d.runners) == -1 {
		err := errors.New("run - unknown runner " + appToRun.GetRunner())
		d.logger.Log(err.Error())
		return err
	}
	switch appToRun.GetRunner() {
	case "node":
		return d.runNode(appToRun)
	case "web":
		return d.runWeb(appToRun)
	case "python":
		return d.runPython(appToRun)
	case "flask":
		return d.runPythonFlask(appToRun)
	default:
		err := errors.New("run - unknown runner " + appToRun.GetRunner())
		d.logger.Log(err.Error())
		return err
	}

}
func (d *Deployer) Kill(appToKill *app.App) error {
	if d.config.GetContainer() {
		mountPoint := path.Join(appToKill.GetRoot(), "ccont_server")
		err := syscall.Unmount(mountPoint, 0)
		if err != nil {
			fmt.Println(err.Error())
		}
		_ = syscall.Rmdir(mountPoint)

	}
	err := appToKill.GetProcess().Signal(syscall.SIGINT)
	if err != nil {
		d.logger.Log(err.Error())
		return err
	}
	_, err = appToKill.GetProcess().Wait()
	if err != nil {
		d.logger.Log(err.Error())
		return err
	}

	jApp := d.GetAppAsJSON(appToKill)
	jApp.Pid = -1
	d.SaveAppToJson(jApp)
	d.RemoveApp(appToKill)
	d.logger.Log("kill - killed app " + appToKill.GetName())
	return nil
}
func (d *Deployer) Remove(appToRemove *app.AppJSON) error {
	a := app.NewAppFromJson(appToRemove)
	if !d.IsAppRunning(a) {
		if strings.HasPrefix(appToRemove.Root, path.Join(d.GetConfig().GetCwd(), d.GetConfig().GetAppsRoot())) {
			err := os.RemoveAll(appToRemove.Root)
			if err != nil {
				d.logger.Log(err.Error())
				return err
			}
			d.RemoveAppFromJson(*appToRemove)
			d.logger.Log("remove - removed app " + appToRemove.Name)
			return nil
		} else {
			err := errors.New("remove - invalid path " + appToRemove.Root)
			d.logger.Log(err.Error())
			return err
		}
	} else {
		err := errors.New("remove - app is running " + appToRemove.Name)
		d.logger.Log(err.Error())
		return err
	}
}
func (d *Deployer) Settings(id string, settings map[string]string) error {
	if appToChange, ok := d.GetApp(id); ok {
		if d.IsAppRunning(appToChange) {
			err := errors.New("settings - app is running")
			d.logger.Log(err.Error())
			return err
		}
		err := errors.New("settings - app in memory")
		d.logger.Log(err.Error())
		return err
	}
	if appJson, ok := d.GetAppD(id); ok {
		changed := false
		for key, value := range settings {
			switch key {
			case "port":
				port, err := strconv.Atoi(value)
				if err != nil {
					d.logger.Log(err.Error())
					return err
				}
				if d.isPortUsed(port) {
					err := errors.New("port in use")
					d.logger.Log(err.Error())
					return err
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
			d.logger.Log("settings - updated app settings " + appJson.Name)
			d.SaveAppToJson(*appJson)
		}
		return nil
	} else {
		err := errors.New("app not found")
		d.logger.Log(err.Error())
		return err
	}
}

func (d *Deployer) IsAppRunning(search *app.App) bool {
	for _, a := range d.apps {
		if a.Name == search.Name || a.Id == search.Id || a.Repo == search.Repo {
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
func (d *Deployer) runNode(appToRun *app.App) error {
	packageJSONFile, _ := ioutil.ReadFile(path.Join(appToRun.GetRoot(), "package.json"))
	packageJson := PackageJSON{}
	err := json.Unmarshal(packageJSONFile, &packageJson)
	if err != nil {
		d.logger.Log(err.Error())
		return err
	}
	port := appToRun.GetPort()
	if port == 0 {
		appToRun.SetPort(d.generatePort())
		port = appToRun.GetPort()
	}

	if d.config.GetContainer() {
		ccont := exec.Command("ccont", "--copy=node-cont", appToRun.GetId(), "-c", "node", packageJson.Main)
		ccont.Dir = appToRun.GetRoot()
		ccont.Env = os.Environ()
		ccont.Env = append(ccont.Env, fmt.Sprintf("CONT_PORT=%d", port))
		ccont.Stdout = os.Stdout
		ccont.Stderr = os.Stderr
		err = ccont.Start()
		if err != nil {
			d.logger.Log(err.Error())
			return err
		}
		appToRun.SetPid(ccont.Process.Pid)
		appToRun.SetProcess(ccont.Process)
	} else {
		node := exec.Command("node", packageJson.Main)
		node.Dir = appToRun.GetRoot()
		node.Env = append(node.Env, fmt.Sprintf("CONT_PORT=%d", port))
		node.Stdout = os.Stdout
		node.Stderr = os.Stderr
		err = node.Start()
		if err != nil {
			d.logger.Log(err.Error())
			return err
		}
		appToRun.SetPid(node.Process.Pid)
		appToRun.SetProcess(node.Process)
	}

	appToRun.SetLastRun(time.Now())
	d.AddApp(appToRun)
	d.SaveAppToJson(d.GetAppAsJSON(appToRun))
	d.logger.Log(fmt.Sprintf("run - starting %s server with pid - %d on port %d", appToRun.GetRunner(), appToRun.GetPid(), appToRun.GetPort()))
	return nil
}
func (d *Deployer) runWeb(appToRun *app.App) error {
	port := appToRun.GetPort()
	if port == 0 {
		appToRun.SetPort(d.generatePort())
		port = appToRun.GetPort()
	}
	serverPath := path.Join(d.GetConfig().GetCwd(), d.GetConfig().GetBasicServer())
	if d.config.GetContainer() {
		mountPoint := path.Join(appToRun.GetRoot(), "ccont_server")
		_ = os.Mkdir(mountPoint, 0755)
		err := syscall.Mount(path.Dir(serverPath), mountPoint, "tmpfs", syscall.MS_BIND, "")
		if err != nil {
			fmt.Println(err.Error())
		}
		ccont := exec.Command("ccont", "--rbind", "--copy=node-cont", appToRun.GetId(), "-c", "node", "ccont_server/server.js")
		ccont.Dir = appToRun.GetRoot()
		ccont.Env = os.Environ()
		ccont.Env = append(ccont.Env, fmt.Sprintf("CONT_PORT=%d", port))
		ccont.Stdout = os.Stdout
		ccont.Stderr = os.Stderr
		err = ccont.Start()
		if err != nil {
			d.logger.Log(err.Error())
			return err
		}
		appToRun.SetPid(ccont.Process.Pid)
		appToRun.SetProcess(ccont.Process)
	} else {
		node := exec.Command("node", path.Join(d.GetConfig().GetCwd(), d.GetConfig().GetBasicServer()))
		node.Dir = appToRun.GetRoot()
		node.Env = append(node.Env, fmt.Sprintf("PORT=%d", port))
		node.Stdout = os.Stdout
		node.Stderr = os.Stderr
		err := node.Start()
		if err != nil {
			fmt.Println(err)
			return err
		}
		appToRun.SetPid(node.Process.Pid)
		appToRun.SetProcess(node.Process)
	}
	appToRun.SetLastRun(time.Now())
	d.AddApp(appToRun)
	d.SaveAppToJson(d.GetAppAsJSON(appToRun))
	d.logger.Log(fmt.Sprintf("run - starting %s server with pid - %d on port %d", appToRun.GetRunner(), appToRun.GetPid(), appToRun.GetPort()))
	return nil
}
func (d *Deployer) runPython(appToRun *app.App) error {
	port := appToRun.GetPort()
	if port == 0 {
		appToRun.SetPort(d.generatePort())
		port = appToRun.GetPort()
	}
	if d.config.GetContainer() {
		ccont := exec.Command("ccont", "--copy=python-cont", appToRun.GetId(), "-c", "python3", "__main__.py")
		ccont.Dir = appToRun.GetRoot()
		ccont.Env = os.Environ()
		ccont.Env = append(ccont.Env, fmt.Sprintf("CONT_PORT=%d", port))
		ccont.Stdout = os.Stdout
		ccont.Stderr = os.Stderr
		err := ccont.Start()
		if err != nil {
			d.logger.Log(err.Error())
			return err
		}
		appToRun.SetPid(ccont.Process.Pid)
		appToRun.SetProcess(ccont.Process)
	} else {
		python := exec.Command("python3", "__main__.py")
		python.Dir = appToRun.GetRoot()
		python.Env = append(python.Env, fmt.Sprintf("PORT=%d", port))
		python.Stdout = os.Stdout
		python.Stderr = os.Stderr
		err := python.Start()
		if err != nil {
			d.logger.Log(err.Error())
			return err
		}
		appToRun.SetPid(python.Process.Pid)
		appToRun.SetProcess(python.Process)
	}
	appToRun.SetLastRun(time.Now())
	d.AddApp(appToRun)
	d.SaveAppToJson(d.GetAppAsJSON(appToRun))
	d.logger.Log(fmt.Sprintf("run - starting %s server with pid - %d on port %d", appToRun.GetRunner(), appToRun.GetPid(), appToRun.GetPort()))
	return nil
}
func (d *Deployer) runPythonFlask(appToRun *app.App) error {
	port := appToRun.GetPort()
	if port == 0 {
		appToRun.SetPort(d.generatePort())
		port = appToRun.GetPort()
	}
	if d.config.GetContainer() {
		pip := exec.Command("ccont", "--copy=flask-cont", appToRun.GetId(), "-c", "pip3", "install", "-r", "requirements.txt")
		pip.Dir = appToRun.GetRoot()
		pip.Stdout = os.Stdout
		pip.Stderr = os.Stderr
		err := pip.Run()
		ccont := exec.Command("ccont", appToRun.GetId(), "-c", "flask", "run", "--host=\"0.0.0.0\"")
		ccont.Dir = appToRun.GetRoot()
		ccont.Env = os.Environ()
		ccont.Env = append(ccont.Env, fmt.Sprintf("CONT_FLASK_RUN_PORT=%d", port))
		ccont.Env = append(ccont.Env, "CONT_LC_ALL=en_US.utf-8")
		ccont.Env = append(ccont.Env, "CONT_LANG=en_US.utf-8")
		ccont.Stdout = os.Stdout
		ccont.Stderr = os.Stderr
		err = ccont.Start()
		if err != nil {
			d.logger.Log(err.Error())
			return err
		}
		appToRun.SetPid(ccont.Process.Pid)
		appToRun.SetProcess(ccont.Process)
	} else {
		pip := exec.Command("pip3", "install", "-r", "requirements.txt")
		pip.Dir = appToRun.GetRoot()
		pip.Stdout = os.Stdout
		pip.Stderr = os.Stderr
		err := pip.Run()
		if err != nil {
			d.logger.Log(err.Error())
			return err
		}
		python := exec.Command("flask", "run", "--host=\"0.0.0.0\"")
		python.Dir = appToRun.GetRoot()
		python.Env = append(python.Env, fmt.Sprintf("FLASK_RUN_PORT=%d", port))
		python.Stdout = os.Stdout
		python.Stderr = os.Stderr
		err = python.Start()
		if err != nil {
			d.logger.Log(err.Error())
			return err
		}
		appToRun.SetPid(python.Process.Pid)
		appToRun.SetProcess(python.Process)
	}
	appToRun.SetLastRun(time.Now())
	d.AddApp(appToRun)
	d.SaveAppToJson(d.GetAppAsJSON(appToRun))
	d.logger.Log(fmt.Sprintf("run - starting %s server with pid - %d on port %d", appToRun.GetRunner(), appToRun.GetPid(), appToRun.GetPort()))
	return nil
}

// load apps from json file into appsD array
func (d *Deployer) GetDeployedApps() []app.AppJSON {
	pth := path.Join(d.GetConfig().GetCwd(), d.GetConfig().GetAppsRoot(), "apps.json")
	jsonFile, _ := ioutil.ReadFile(pth)
	appsD := AppsJSON{}
	err := json.Unmarshal(jsonFile, &appsD)
	if err != nil {
		fmt.Println("json - " + err.Error())
	}
	return appsD.Apps
}

// save currently running apps to json
func (d *Deployer) RemoveAppFromJson(a app.AppJSON) {
	pth := path.Join(d.GetConfig().GetCwd(), d.GetConfig().GetAppsRoot(), "apps.json")
	appsJson := d.GetDeployedApps()
	if pos := d.ContainsAppJSON(&appsJson, &a); pos != -1 {
		appsJson = append(appsJson[:pos], appsJson[pos+1:]...)
	}
	apps, _ := json.Marshal(AppsJSON{Apps: appsJson})
	err := ioutil.WriteFile(pth, apps, 0775)
	if err != nil {
		fmt.Println(err)
	}
}
func (d *Deployer) SaveAppToJson(appToSave app.AppJSON) {
	pth := path.Join(d.GetConfig().GetCwd(), d.GetConfig().GetAppsRoot(), "apps.json")
	appsJson := d.GetDeployedApps()
	appToSave.Pid = -1
	appToSave.Uptime = ""
	if pos := d.ContainsAppJSON(&appsJson, &appToSave); pos == -1 {
		appsJson = append(appsJson, appToSave)
	} else {
		appsJson = append(appsJson[:pos], appsJson[pos+1:]...)
		appsJson = append(appsJson, appToSave)
	}
	apps, _ := json.Marshal(AppsJSON{Apps: appsJson})
	err := ioutil.WriteFile(pth, apps, 0775)
	if err != nil {
		fmt.Println(err)
	}
}
func (d *Deployer) initAppsJson() {
	pth := path.Join(d.GetConfig().GetCwd(), d.GetConfig().GetAppsRoot(), "apps.json")
	if _, err := os.Stat(pth); err != nil {
		fmt.Println("initializing apps folder")
		emptyArr, _ := json.Marshal(&AppsJSON{Apps: []app.AppJSON{}})
		err := ioutil.WriteFile(pth, emptyArr, 0775)
		if err != nil {
			fmt.Println("apps json init failed - " + err.Error())
		}
	} else {
		folders, _ := ioutil.ReadDir(path.Join(d.GetConfig().GetCwd(), d.GetConfig().GetAppsRoot()))
		appsD := d.GetDeployedApps()
		for _, a := range appsD {
			if !utils.ContainsFile(a.Name, &folders) {
				fmt.Println("not found ", a.Name)
				d.RemoveAppFromJson(a)
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
	} else if port == d.port {
		return true
	}
	for _, a := range d.GetDeployedApps() {
		if port == a.Port {
			return true
		}
	}
	return false
}
func (d *Deployer) ContainsAppJSON(arr *[]app.AppJSON, ac *app.AppJSON) int {
	for i, a := range *arr {
		if a.Id == ac.Id || a.Name == ac.Name || a.Repo == ac.Repo {
			return i
		}
	}
	return -1
}
