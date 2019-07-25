package controllers

import (
	"../app"
	"../config"
	"../encryption/auth"
	"../http/requests"
	httpresp "../http/responses"
	"../logger"
	"../utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	config   *config.Config
	deployer *Deployer
	logger   *logger.Logger
}

func NewHandler(cfg *config.Config, d *Deployer) Handler {
	h := Handler{}
	h.config = cfg
	h.deployer = d
	h.logger = logger.NewLogger(logger.LOG_SERVER)
	return h
}

func (h *Handler) SetConfig(c *config.Config) {
	h.config = c
}

func (h *Handler) GetConfig() *config.Config {
	return h.config
}
func (h *Handler) SetDeployer(d *Deployer) {
	h.deployer = d
}

func (h *Handler) GetDeployer() *Deployer {
	return h.deployer
}
func (h *Handler) SetLogger(l *logger.Logger) {
	h.logger = l
}

func (h *Handler) GetLogger() *logger.Logger {
	return h.logger
}

func (h *Handler) HandleDeploy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("deploy - bad cookie " + r.RemoteAddr)
		httpresp.ResponseUnauthorized(w)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if auth.VerifyToken(token, h.GetConfig().GetSecret()) {
			if r.Method == http.MethodPost {
				jsonBody := getJsonMap(&r.Body)
				repo := jsonBody["repo"]
				runner := jsonBody["runner"]
				hostname := jsonBody["hostname"]
				port, err := strconv.Atoi(jsonBody["port"])
				if err != nil {
					port = 0
				}
				a, err := h.GetDeployer().Deploy(repo, runner, hostname, port)
				if err != nil {
					h.logger.Log("deploy - " + err.Error())
					if err.Error() == "exit status 128" {
						httpresp.ResponseBadRequest(w, requests.ErrorResponse{Id: utils.GetNameFromRepo(repo), Message: "invalid repo"})
					} else {
						httpresp.ResponseInternalServerError(w, requests.ErrorResponse{Id: utils.GetNameFromRepo(repo), Message: err.Error()})
					}
					return
				}
				err = h.GetDeployer().Install(a)
				a.Print()
				if err != nil {
					h.logger.Log(err.Error())
					httpresp.ResponseInternalServerError(w, requests.ErrorResponse{Message: err.Error(), Id: a.GetName()})
					return
				}
				err = h.GetDeployer().Run(a)
				if err != nil {
					h.logger.Log(err.Error())
					httpresp.ResponseInternalServerError(w, requests.ErrorResponse{Message: err.Error(), Id: a.GetName()})
					return
				}
				h.logger.Log("deploy - deployed " + repo)
				httpresp.ResponseCreated(w, requests.SuccessResponse{Message: "deployed", App: h.GetDeployer().GetAppAsJSON(a)})
			} else {
				h.logger.Log("deploy - method not allowed" + r.RemoteAddr)
				httpresp.ResponseMethodNotAllowed(w)
			}
		} else {
			h.logger.Log("deploy - bad token " + r.RemoteAddr)
			httpresp.ResponseUnauthorized(w)
		}
	}
}
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("deploy - bad token " + r.RemoteAddr)
		httpresp.ResponseUnauthorized(w)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if auth.VerifyToken(token, h.GetConfig().GetSecret()) {
			if r.Method == http.MethodPost {
				body := getJsonMap(&r.Body)
				name := body["app"]
				if aJson, ok := h.GetDeployer().GetAppD(name); ok {
					a := app.NewAppFromJson(aJson)
					err := h.GetDeployer().Update(a)
					if err != nil {
						h.logger.Log(err.Error())
						httpresp.ResponseInternalServerError(w, requests.ErrorResponse{Message: "update failed", Id: name})
						return
					}
					h.logger.Log("update - updated app " + name)
					httpresp.ResponseOK(w, requests.SuccessResponse{Message: "updated", App: *aJson})
				} else {
					h.logger.Log("update - app not found " + name)
					httpresp.ResponseBadRequest(w, requests.ErrorResponse{Message: "app not found", Id: name})
				}
			} else {
				h.logger.Log("update - method not allowed" + r.RemoteAddr)
				httpresp.ResponseMethodNotAllowed(w)
			}
		} else {
			h.logger.Log("update - bad token " + r.RemoteAddr)
			httpresp.ResponseUnauthorized(w)
		}
	}
}
func (h *Handler) HandleRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("run - bad token " + r.RemoteAddr)
		httpresp.ResponseUnauthorized(w)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if auth.VerifyToken(token, h.GetConfig().GetSecret()) {
			if r.Method == http.MethodPost {
				body := getJsonMap(&r.Body)
				name := body["app"]
				if aJson, ok := h.GetDeployer().GetAppD(name); ok {
					a := app.NewAppFromJson(aJson)
					if h.GetDeployer().IsAppRunning(a) {
						h.logger.Log("run - app already running " + name)
						httpresp.ResponseNoContent(w, requests.ErrorResponse{Message: "app already running", Id: a.GetName()})
					} else {
						err := h.GetDeployer().Run(a)
						if err != nil {
							h.logger.Log(err.Error())
							httpresp.ResponseInternalServerError(w, requests.ErrorResponse{Message: "unable to run", Id: name})
							return
						}
						h.logger.Log("run - running app " + name)
						httpresp.ResponseOK(w, requests.SuccessResponse{Message: "running", App: *aJson})
					}
				} else {
					h.logger.Log("run - app not found " + name)
					httpresp.ResponseBadRequest(w, requests.ErrorResponse{Message: "app not found", Id: name})
				}
			} else {
				h.logger.Log("run - method not allowed " + r.RemoteAddr)
				httpresp.ResponseMethodNotAllowed(w)
			}
		} else {
			h.logger.Log("run - bad token " + r.RemoteAddr)
			httpresp.ResponseUnauthorized(w)
		}
	}
}
func (h *Handler) HandleFind(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("find - bad token " + r.RemoteAddr)
		httpresp.ResponseUnauthorized(w)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if auth.VerifyToken(token, h.GetConfig().GetSecret()) {
			if r.Method == http.MethodGet {
				query := r.URL.Query().Get("app")
				as := h.GetDeployer().GetAppsAsJSON()
				asD := h.GetDeployer().GetDeployedApps()
				h.logger.Log("find - querying apps")
				if query == "" {
					httpresp.ResponseOK(w, requests.FindResponse{Running: &as, Deployed: &asD})
				} else {
					httpresp.ResponseOK(w, requests.FindResponse{Running: queryApps(query, &as), Deployed: queryApps(query, &asD)})
				}
			} else {
				h.logger.Log("find - method not allowed " + r.RemoteAddr)
				httpresp.ResponseMethodNotAllowed(w)
			}
		} else {
			h.logger.Log("find - bad token " + r.RemoteAddr)
			httpresp.ResponseUnauthorized(w)
		}
	}
}

// TODO: kill error handling
func (h *Handler) HandleKill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("kill - bad token " + r.RemoteAddr)
		httpresp.ResponseUnauthorized(w)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if auth.VerifyToken(token, h.GetConfig().GetSecret()) {
			if r.Method == http.MethodPost {
				body := getJsonMap(&r.Body)
				name := body["app"]
				if a, ok := h.GetDeployer().GetApp(name); ok {
					err := h.GetDeployer().Kill(a)
					if err != nil {
						h.logger.Log(err.Error())
						httpresp.ResponseInternalServerError(w, requests.ErrorResponse{Message: "unable to kill", Id: name})
						return
					}
					h.logger.Log("kill - killed app " + name)
					httpresp.ResponseOK(w, requests.SuccessResponse{Message: "killed", App: h.GetDeployer().GetAppAsJSON(a)})
				} else {
					h.logger.Log("kill - app not found " + name)
					httpresp.ResponseInternalServerError(w, requests.ErrorResponse{Message: "app not found", Id: name})
				}
			} else {
				h.logger.Log("kill - method not allowed " + r.RemoteAddr)
				httpresp.ResponseMethodNotAllowed(w)
			}
		} else {
			h.logger.Log("kill - bad token " + r.RemoteAddr)
			httpresp.ResponseUnauthorized(w)
		}
	}

}
func (h *Handler) HandleRemove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("remove - bad token " + r.RemoteAddr)
		httpresp.ResponseUnauthorized(w)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if auth.VerifyToken(token, h.GetConfig().GetSecret()) {
			if r.Method == http.MethodPost {
				body := getJsonMap(&r.Body)
				name := body["app"]
				if aJson, ok := h.GetDeployer().GetAppD(name); ok {
					if a, ok := h.GetDeployer().GetApp(aJson.Id); ok {
						_ = h.GetDeployer().Kill(a)
					}
					err := h.GetDeployer().Remove(aJson)
					if err != nil {
						h.logger.Log(err.Error())
						httpresp.ResponseInternalServerError(w, requests.ErrorResponse{Message: err.Error(), Id: name})
						return
					}
					h.logger.Log("remove - removed app " + name)
					httpresp.ResponseOK(w, requests.SuccessResponse{Message: "removed", App: *aJson})
				} else {
					h.logger.Log("remove - a not found " + name)
					httpresp.ResponseInternalServerError(w, requests.ErrorResponse{Message: "app not found", Id: name})
				}
			} else {
				h.logger.Log("remove - method not allowed")
				httpresp.ResponseMethodNotAllowed(w)
			}
		} else {
			h.logger.Log("remove - bad token " + r.RemoteAddr)
			httpresp.ResponseUnauthorized(w)
		}
	}

}
func (h *Handler) HandleSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("settings - bad token " + r.RemoteAddr)
		httpresp.ResponseUnauthorized(w)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if auth.VerifyToken(token, h.GetConfig().GetSecret()) {
			if r.Method == http.MethodPost {
				req := requests.SettingsRequest{}
				req.Read(&r.Body)
				err := h.GetDeployer().Settings(req.Id, req.Settings)
				if err != nil {
					h.logger.Log(err.Error())
					httpresp.ResponseInternalServerError(w, requests.ErrorResponse{Message: err.Error(), Id: req.Id})
				} else {
					h.logger.Log("settings - updated a settings " + req.Id)
					httpresp.ResponseOK(w, requests.SuccessResponse{Message: "updated", App: app.AppJSON{Id: req.Id}})
				}
			} else {
				h.logger.Log("settings - method not allowed")
				httpresp.ResponseMethodNotAllowed(w)
			}
		} else {
			h.logger.Log("settings - bad token " + r.RemoteAddr)
			httpresp.ResponseUnauthorized(w)
		}
	}
}
func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	h.logger.Log("root - " + r.URL.Path)
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("root - redirecting bad token " + r.RemoteAddr)
		http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if auth.VerifyToken(token, h.GetConfig().GetSecret()) {
			var absPth string
			pth := strings.Replace(r.URL.String(), "/", string(filepath.Separator), -1)
			root := "client/dist"
			if strings.HasPrefix(pth, "/node_modules") {
				root = "client"
			}
			if pth == "/" || pth == "\\" {
				absPth = root
			} else {
				absPth = path.Join(root, pth)
			}
			if fi, err := os.Stat(absPth); err == nil && fi.IsDir() {
				if dir, err := ioutil.ReadDir(absPth); err == nil {
					if utils.ContainsFile("index.html", &dir) {
						w.Header().Set("Content-Type", "text/html; charset=utf-8")
						http.ServeFile(w, r, path.Join(absPth, "index.html"))
					}
				} else {
					httpresp.ResponseInternalServerError(w, nil)
				}
			} else if err == nil {
				http.ServeFile(w, r, absPth)
			} else {
				httpresp.ResponseNotFound(w)
			}
		} else {
			h.logger.Log("root - redirecting bad token " + r.RemoteAddr)
			http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
		}
	}
}
func (h *Handler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		pass := r.FormValue("password")
		user := r.FormValue("username")
		auser := h.GetConfig().GetUser()
		apass := h.GetConfig().GetPass()
		if auth.VerifyCredentials(auser, apass, user, pass) {
			token := auth.GenerateToken(h.GetConfig().GetSecret())
			cookie := http.Cookie{Name: "Authorization", Value: fmt.Sprintf("Bearer %s", token), Path: "/", Expires: time.Now().Add(24 * time.Hour)}
			h.logger.Log("auth - authorized " + r.RemoteAddr)
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			length, _ := w.Write(utils.RenderLoginPage())
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	case http.MethodGet:
		if cookie, err := r.Cookie("Authorization"); err == nil {
			token := strings.Split(cookie.Value, "Bearer ")[1]
			if auth.VerifyToken(token, h.GetConfig().GetSecret()) {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				break
			}
		}
		w.WriteHeader(http.StatusOK)
		length, _ := w.Write(utils.RenderLoginPage())
		w.Header().Set("Content-Length", strconv.Itoa(length))
	}
}

func queryApps(search string, as *[]app.AppJSON) *[]app.AppJSON {
	var out []app.AppJSON
	for _, a := range *as {
		if strings.Contains(a.Name, strings.ToLower(search)) || strings.Contains(a.Id, search) {
			out = append(out, a)
		}
	}
	return &out
}
func getJsonMap(body *io.ReadCloser) map[string]string {
	output := make(map[string]string)
	jsonBytes, _ := ioutil.ReadAll(*body)
	err := json.Unmarshal(jsonBytes, &output)
	if err != nil {
		fmt.Println(err)
	}
	return output
}
