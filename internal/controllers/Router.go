package controllers

import (
	"../config"
	"../logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type RouterHandler struct {
	deployer                  *Deployer
	config                    *config.Config
	statusInternalServerError []byte
	statusOK                  []byte
	statusUnauthorized        []byte
	statusNotFound            []byte
	statusMethodNotAllowed    []byte
	hosts                     map[string]string
	logger                    *logger.Logger
}

func NewRouterHandler(d *Deployer, c *config.Config) *RouterHandler {
	rh := RouterHandler{}
	rh.deployer = d
	rh.config = c
	rh.hosts = map[string]string{}
	rh.logger = logger.NewLogger(logger.LOG_SERVER)
	rh.statusInternalServerError = []byte("( ͠° ͟ʖ ͡°) 500 INTERNAL SERVER ERROR")
	rh.statusNotFound = []byte("( ͡° ʖ̯ ͡°) 404 NOT FOUND")
	rh.statusUnauthorized = []byte("( ͠° ͟ʖ ͡°) 401 UNAUTHORIZED")
	rh.statusMethodNotAllowed = []byte("( ͠° ͟ʖ ͡°) 405 METHOD NOT ALLOWED")
	rh.statusOK = []byte("( ͡ᵔ ͜ʖ ͡ᵔ ) 200 OK")
	rh.UpdateHosts()
	return &rh
}

func (rh *RouterHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	rh.logger.Log(fmt.Sprintf("router - %s %s", r.URL.Path, r.RemoteAddr, ))
	protocol := "http://"
	if r.URL.Scheme == "https" {
		protocol = "https://"
	}
	host := r.Host
	url := ""
	if host == rh.config.GetHostname() || strings.HasPrefix(host, "dev.") {
		url = protocol + host + ":" + strconv.Itoa(rh.config.GetPort())
	} else {
		url = protocol + host + ":" + rh.hosts[host]
	}
	if url == protocol+host+":" {
		w.WriteHeader(http.StatusNotFound)
		length, _ := w.Write(rh.statusNotFound)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}
func (rh *RouterHandler) GetHosts() *map[string]string {
	return &rh.hosts
}
func (rh *RouterHandler) UpdateHosts() {
	pth := path.Join(rh.config.GetAppsRoot(), "apps.json")
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
