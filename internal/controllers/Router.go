package controllers

import (
	"../config"
	httpresp "../http/responses"
	"../logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type RouterHandler struct {
	deployer *Deployer
	config   *config.Config
	hosts    map[string]string
	logger   *logger.Logger
}

func NewRouterHandler(d *Deployer, c *config.Config) *RouterHandler {
	rh := RouterHandler{}
	rh.deployer = d
	rh.config = c
	rh.hosts = map[string]string{}
	rh.logger = logger.NewLogger(logger.LOG_SERVER)
	rh.UpdateHosts()
	return &rh
}

func (rh *RouterHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	rh.logger.Log(fmt.Sprintf("router - %s %s", r.URL.Path, r.RemoteAddr, ))
	protocol := "http://"
	if r.URL.Scheme == "https" {
		protocol = "https://"
	}
	host := strings.Split(r.Host, ":")[0]
	newurl := ""
	if host == rh.config.GetHostname() || strings.HasPrefix(host, "dev.") {
		newurl = protocol + host + ":" + strconv.Itoa(rh.config.GetPort())
	} else {
		newurl = protocol + host + ":" + rh.hosts[host]
	}
	if newurl == protocol+host+":" {
		httpresp.ResponseNotFound(w)
	} else {
		http.Redirect(w, r, newurl, http.StatusTemporaryRedirect)
	}
}
func (rh *RouterHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	rh.logger.Log(fmt.Sprintf("router - %s %s", r.URL.Path, r.RemoteAddr, ))
	protocol := "http://"
	if r.URL.Scheme == "https" {
		protocol = "https://"
	}
	host := strings.Split(r.Host, ":")[0]
	var newurl string
	if strings.HasPrefix(r.Host, "dev.") || host == rh.config.GetHostname() {
		newurl = protocol + host + ":" + strconv.Itoa(rh.config.GetPort())
	} else if proxiedPort, ok := rh.hosts[host]; ok {
		newurl = protocol + host + ":" + proxiedPort
	} else {
		httpresp.ResponseNotFound(w)
		return
	}
	u, err := url.Parse(newurl)
	if err != nil {
		httpresp.ResponseNotFound(w)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
func (rh *RouterHandler) GetHosts() *map[string]string {
	return &rh.hosts
}
func (rh *RouterHandler) UpdateHosts() {
	pth := path.Join(rh.config.GetCwd(), rh.config.GetAppsRoot(), "apps.json")
	jsonFile, _ := ioutil.ReadFile(pth)
	appsJson := AppsJSON{}
	err := json.Unmarshal(jsonFile, &appsJson)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, a := range appsJson.Apps {
		fmt.Println(a.Hostname + " " + strconv.Itoa(a.Port))
		rh.hosts[a.Hostname] = strconv.Itoa(a.Port)
	}
}
