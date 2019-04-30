package controllers

import (
	"../config"
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
	newurl := ""
	if host == rh.config.GetHostname() || strings.HasPrefix(host, "dev.") {
		newurl = protocol + host + ":" + strconv.Itoa(rh.config.GetPort())
	} else {
		newurl = protocol + host + ":" + rh.hosts[host]
	}
	if newurl == protocol+host+":" {
		w.WriteHeader(http.StatusNotFound)
		length, _ := w.Write(rh.statusNotFound)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		http.Redirect(w, r, newurl, http.StatusPermanentRedirect)
	}
}
func (rh *RouterHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	rh.logger.Log(fmt.Sprintf("router - %s %s", r.URL.Path, r.RemoteAddr, ))
	protocol := "http://"
	if r.URL.Scheme == "https" {
		protocol = "https://"
	}

	host := r.Host
	newurl := ""
	proxiedPort := ":" + rh.hosts[host]
	if host == rh.config.GetHostname() || strings.HasPrefix(host, "dev.") {
		newurl = protocol + host + ":" + strconv.Itoa(rh.config.GetPort())
	} else if proxiedPort == ":" {
		w.WriteHeader(http.StatusNotFound)
		length, _ := w.Write(rh.statusNotFound)
		w.Header().Set("Content-Length", strconv.Itoa(length))
		return
	} else {
		newurl = protocol + host + proxiedPort
	}
	u, err := url.Parse(newurl)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		length, _ := w.Write(rh.statusNotFound)
		w.Header().Set("Content-Length", strconv.Itoa(length))
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	//proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	length, _ := w.Write(rh.statusNotFound)
	//	w.Header().Set("Content-Length", strconv.Itoa(length))
	//}
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
