package controllers

import (
	"../config"
	"../logger"
	"../utils"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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
	config                    *config.Config
	deployer                  *Deployer
	statusInternalServerError []byte
	statusOK                  []byte
	statusUnauthorized        []byte
	statusNotFound            []byte
	statusMethodNotAllowed    []byte
	logger                    *logger.Logger
}

type DeployRequest struct {
	Token string `json:"token"`
	Repo  string `json:"repo"`
}
type FindResponse struct {
	Running  *[]AppJSON `json:"running"`
	Deployed *[]AppJSON `json:"deployed"`
}
type ErrorResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}
type SuccessResponse struct {
	Message string  `json:"message"`
	App     AppJSON `json:"app"`
}
type SettingsRequest struct {
	Id       string            `json:"id"`
	Settings map[string]string `json:"settings"`
}

func NewHandler(cfg *config.Config, d *Deployer) Handler {
	h := Handler{}
	h.config = cfg
	h.deployer = d
	h.statusInternalServerError = []byte("( ͠° ͟ʖ ͡°) 500 INTERNAL SERVER ERROR")
	h.statusNotFound = []byte("( ͡° ʖ̯ ͡°) 404 NOT FOUND")
	h.statusUnauthorized = []byte("( ͠° ͟ʖ ͡°) 401 UNAUTHORIZED")
	h.statusMethodNotAllowed = []byte("( ͠° ͟ʖ ͡°) 405 METHOD NOT ALLOWED")
	h.statusOK = []byte("( ͡ᵔ ͜ʖ ͡ᵔ ) 200 OK")
	h.logger = logger.NewLogger(logger.LOG_SERVER)
	return h
}
func (s *SettingsRequest) Read(body *io.ReadCloser) {
	jsonBytes, _ := ioutil.ReadAll(*body)
	err := json.Unmarshal(jsonBytes, s)
	if err != nil {
		fmt.Println(err)
	}
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
		w.WriteHeader(http.StatusUnauthorized)
		length, _ := w.Write(h.statusUnauthorized)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if h.verifyToken(token) {
			if r.Method == http.MethodPost {
				jsonBody := getJsonMap(&r.Body)
				repo := jsonBody["repo"]
				runner := jsonBody["runner"]
				hostname := jsonBody["hostname"]
				port, err := strconv.Atoi(jsonBody["port"])
				if err != nil {
					port = 0
				}
				app, err := h.GetDeployer().Deploy(repo, runner, hostname, port)
				if err != nil {
					h.logger.Log(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					response := ErrorResponse{Id: utils.GetNameFromRepo(repo)}
					if err.Error() == "exit status 128" {
						response.Message = "invalid repo"
					} else {
						response.Message = err.Error()
					}
					jsonResponse, _ := json.Marshal(&response)
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
					return
				}
				err = h.GetDeployer().Install(app)
				app.Print()
				if err != nil {
					h.logger.Log(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					jsonResponse, _ := json.Marshal(ErrorResponse{Message: err.Error(), Id: app.GetName()})
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
					return
				}
				err = h.GetDeployer().Run(app)
				if err != nil {
					h.logger.Log(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					jsonResponse, _ := json.Marshal(ErrorResponse{Message: err.Error(), Id: app.GetName()})
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
					return
				}
				h.logger.Log("deployed " + repo)
				w.WriteHeader(http.StatusOK)
				jsonResponse, _ := json.Marshal(SuccessResponse{Message: "deployed", App: h.GetDeployer().GetAppAsJSON(app)})
				length, _ := w.Write(jsonResponse)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
				length, _ := w.Write(h.statusMethodNotAllowed)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			}
		} else {
			h.logger.Log("deploy - bad token " + r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			length, _ := w.Write(h.statusUnauthorized)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	}
}
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("deploy - bad token " + r.RemoteAddr)
		w.WriteHeader(http.StatusUnauthorized)
		length, _ := w.Write(h.statusUnauthorized)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if h.verifyToken(token) {
			if r.Method == http.MethodPost {
				body := getJsonMap(&r.Body)
				name := body["app"]
				if appJson, ok := h.GetDeployer().GetAppD(name); ok {
					app := NewAppFromJson(appJson)
					err := h.GetDeployer().Update(app)
					if err != nil {
						h.logger.Log(err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						jsonResponse, _ := json.Marshal(ErrorResponse{Message: "update failed", Id: name})
						length, _ := w.Write(jsonResponse)
						w.Header().Set("Content-Length", strconv.Itoa(length))
						return
					}
					h.logger.Log("update - updated app " + name)
					w.WriteHeader(http.StatusOK)
					jsonResponse, _ := json.Marshal(SuccessResponse{Message: "updated", App: *appJson})
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
				} else {
					h.logger.Log("update - app not found " + name)
					w.WriteHeader(http.StatusForbidden)
					jsonResponse, _ := json.Marshal(ErrorResponse{Message: "app not found", Id: name})
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
				}
			} else {
				h.logger.Log("update - method not allowed" + r.RemoteAddr)
				w.WriteHeader(http.StatusMethodNotAllowed)
				length, _ := w.Write(h.statusMethodNotAllowed)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			}
		} else {
			h.logger.Log("update - bad token " + r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			length, _ := w.Write(h.statusUnauthorized)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	}
}
func (h *Handler) HandleRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("run - bad token " + r.RemoteAddr)
		w.WriteHeader(http.StatusUnauthorized)
		length, _ := w.Write(h.statusUnauthorized)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if h.verifyToken(token) {
			if r.Method == http.MethodPost {
				body := getJsonMap(&r.Body)
				name := body["app"]
				if appJson, ok := h.GetDeployer().GetAppD(name); ok {
					app := NewAppFromJson(appJson)
					if h.GetDeployer().IsAppRunning(app) {
						h.logger.Log("run - app already running " + name)
						w.WriteHeader(http.StatusNoContent)
						jsonResponse, _ := json.Marshal(ErrorResponse{Message: "app already running", Id: app.GetName()})
						length, _ := w.Write(jsonResponse)
						w.Header().Set("Content-Length", strconv.Itoa(length))
					} else {
						err := h.GetDeployer().Run(app)
						if err != nil {
							h.logger.Log(err.Error())
							w.WriteHeader(http.StatusInternalServerError)
							jsonResponse, _ := json.Marshal(ErrorResponse{Message: "unable to run", Id: name})
							length, _ := w.Write(jsonResponse)
							w.Header().Set("Content-Length", strconv.Itoa(length))
							return
						}
						h.logger.Log("run - running app " + name)
						w.WriteHeader(http.StatusOK)
						jsonResponse, _ := json.Marshal(SuccessResponse{Message: "running", App: *appJson})
						length, _ := w.Write(jsonResponse)
						w.Header().Set("Content-Length", strconv.Itoa(length))
					}
				} else {
					h.logger.Log("run - app not found " + name)
					w.WriteHeader(http.StatusForbidden)
					jsonResponse, _ := json.Marshal(ErrorResponse{Message: "app not found", Id: name})
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
				}
			} else {
				h.logger.Log("run - method not allowed " + r.RemoteAddr)
				w.WriteHeader(http.StatusMethodNotAllowed)
				length, _ := w.Write(h.statusMethodNotAllowed)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			}
		} else {
			h.logger.Log("run - bad token " + r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			length, _ := w.Write(h.statusUnauthorized)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	}
}
func (h *Handler) HandleFind(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("find - bad token " + r.RemoteAddr)
		w.WriteHeader(http.StatusUnauthorized)
		length, _ := w.Write(h.statusUnauthorized)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if h.verifyToken(token) {
			if r.Method == http.MethodGet {
				query := r.URL.Query().Get("app")
				apps := h.GetDeployer().GetAppsAsJSON()
				appsD := h.GetDeployer().GetDeployedApps()
				var jsonResponse []byte
				if query == "" {
					jsonResponse, _ = json.Marshal(&FindResponse{Running: &apps, Deployed: &appsD})
				} else {
					jsonResponse, _ = json.Marshal(&FindResponse{Running: queryApps(query, &apps), Deployed: queryApps(query, &appsD)})
				}
				h.logger.Log("find - querying apps")
				w.WriteHeader(http.StatusOK)
				length, _ := w.Write(jsonResponse)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			} else {
				h.logger.Log("find - method not allowed " + r.RemoteAddr)
				w.WriteHeader(http.StatusMethodNotAllowed)
				length, _ := w.Write(h.statusMethodNotAllowed)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			}
		} else {
			h.logger.Log("find - bad token " + r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			length, _ := w.Write(h.statusUnauthorized)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	}
}

// TODO: kill error handling
func (h *Handler) HandleKill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("kill - bad token " + r.RemoteAddr)
		w.WriteHeader(http.StatusUnauthorized)
		length, _ := w.Write(h.statusUnauthorized)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if h.verifyToken(token) {
			if r.Method == http.MethodPost {
				body := getJsonMap(&r.Body)
				name := body["app"]
				if app, ok := h.GetDeployer().GetApp(name); ok {
					err := h.GetDeployer().Kill(app)
					if err != nil {
						h.logger.Log(err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						jsonResponse, _ := json.Marshal(ErrorResponse{Message: "unable to kill", Id: name})
						length, _ := w.Write(jsonResponse)
						w.Header().Set("Content-Length", strconv.Itoa(length))
						return
					}
					h.logger.Log("kill - killed app " + name)
					jsonResponse, _ := json.Marshal(SuccessResponse{Message: "killed", App: h.GetDeployer().GetAppAsJSON(app)})
					w.WriteHeader(http.StatusOK)
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					jsonResponse, _ := json.Marshal(ErrorResponse{Message: "app not found", Id: name})
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
				}
			} else {
				h.logger.Log("kill - method not allowed " + r.RemoteAddr)
				w.WriteHeader(http.StatusMethodNotAllowed)
				length, _ := w.Write(h.statusMethodNotAllowed)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			}
		} else {
			h.logger.Log("kill - bad token " + r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			length, _ := w.Write(h.statusUnauthorized)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	}

}
func (h *Handler) HandleRemove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("remove - bad token " + r.RemoteAddr)
		w.WriteHeader(http.StatusUnauthorized)
		length, _ := w.Write(h.statusUnauthorized)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if h.verifyToken(token) {
			if r.Method == http.MethodPost {
				body := getJsonMap(&r.Body)
				name := body["app"]
				if appJson, ok := h.GetDeployer().GetAppD(name); ok {
					if app, ok := h.GetDeployer().GetApp(appJson.Id); ok {
						_ = h.GetDeployer().Kill(app)
					}
					err := h.GetDeployer().Remove(appJson)
					if err != nil {
						h.logger.Log(err.Error())
						jsonResponse, _ := json.Marshal(ErrorResponse{Message: "removed", Id: name})
						w.WriteHeader(http.StatusInternalServerError)
						length, _ := w.Write(jsonResponse)
						w.Header().Set("Content-Length", strconv.Itoa(length))
						return
					}
					h.logger.Log("remove - removed app " + name)
					jsonResponse, _ := json.Marshal(SuccessResponse{Message: "removed", App: *appJson})
					w.WriteHeader(http.StatusOK)
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
				} else {
					h.logger.Log("remove - app not found " + name)
					w.WriteHeader(http.StatusInternalServerError)
					jsonResponse, _ := json.Marshal(ErrorResponse{Message: "app not found", Id: name})
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
				}
			} else {
				h.logger.Log("remove - method not allowed")
				w.WriteHeader(http.StatusMethodNotAllowed)
				length, _ := w.Write(h.statusMethodNotAllowed)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			}
		} else {
			h.logger.Log("remove - bad token " + r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			length, _ := w.Write(h.statusUnauthorized)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	}

}
func (h *Handler) HandleSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("settings - bad token " + r.RemoteAddr)
		w.WriteHeader(http.StatusUnauthorized)
		length, _ := w.Write(h.statusUnauthorized)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if h.verifyToken(token) {
			if r.Method == http.MethodPost {
				req := SettingsRequest{}
				req.Read(&r.Body)
				err := h.GetDeployer().Settings(req.Id, req.Settings)
				if err != nil {
					h.logger.Log(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					jsonResponse, _ := json.Marshal(ErrorResponse{Message: err.Error(), Id: req.Id})
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
				} else {
					h.logger.Log("settings - updated app settings " + req.Id)
					jsonResponse, _ := json.Marshal(SuccessResponse{Message: "updated", App: AppJSON{Id: req.Id}})
					length, _ := w.Write(jsonResponse)
					w.Header().Set("Content-Length", strconv.Itoa(length))
				}

			} else {
				h.logger.Log("settings - method not allowed")
				w.WriteHeader(http.StatusMethodNotAllowed)
				length, _ := w.Write(h.statusMethodNotAllowed)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			}
		} else {
			h.logger.Log("settings - bad token " + r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			length, _ := w.Write(h.statusUnauthorized)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	}
}
func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	h.logger.Log(r.URL.Path)
	if cookie, err := r.Cookie("Authorization"); err != nil {
		h.logger.Log("root - redirecting bad token " + r.RemoteAddr)
		http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if h.verifyToken(token) {
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
					w.WriteHeader(http.StatusInternalServerError)
					length, _ := w.Write(h.statusInternalServerError)
					w.Header().Set("Content-Length", strconv.Itoa(length))
				}
			} else if err == nil {
				http.ServeFile(w, r, absPth)
			} else {
				w.WriteHeader(http.StatusNotFound)
				length, _ := w.Write(h.statusNotFound)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			}
		} else {
			h.logger.Log("root - redirecting bad token " + r.RemoteAddr)
			http.Redirect(w, r, "/auth", http.StatusMovedPermanently)
		}
	}
}
func (h *Handler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		pass := r.FormValue("password")
		user := r.FormValue("username")
		if h.verifyCredentials(user, pass) {
			token := h.makeToken()
			cookie := http.Cookie{Name: "Authorization", Value: fmt.Sprintf("Bearer %s", token), Path: "/", Expires: time.Now().Add(24 * time.Hour)}
			h.logger.Log("auth - authorized " + r.RemoteAddr)
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "/auth", http.StatusMovedPermanently)
		}
	case http.MethodGet:
		if cookie, err := r.Cookie("Authorization"); err == nil {
			token := strings.Split(cookie.Value, "Bearer ")[1]
			if h.verifyToken(token) {
				http.Redirect(w, r, "/", http.StatusMovedPermanently)
				break
			}
		}
		loginPage := utils.RenderLoginPage()
		w.WriteHeader(http.StatusOK)
		length, _ := w.Write(loginPage)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	}
}

func (h *Handler) makeToken() string {
	expires := time.Now().Unix() + int64(24*time.Hour)
	//expires := time.Now().Unix()
	type JSTClaims struct {
		Data string `json:"data"`
		jwt.StandardClaims
	}
	// TODO: data
	claims := JSTClaims{
		"bar",
		jwt.StandardClaims{ExpiresAt: expires, Issuer: "issuer.7aske.com"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(h.config.GetSecret())
	return tokenString
}
func (h *Handler) verifyCredentials(user string, pass string) bool {
	configUser := h.GetConfig().GetUser()
	configPassHash := utils.Hash(h.GetConfig().GetPass())
	passHash := utils.Hash(pass)
	return configPassHash == passHash && configUser == strings.ToLower(user)
}
func (h *Handler) verifyToken(tokenString string) bool {
	if _, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method %v", token.Header["alg"])
		}
		return h.config.GetSecret(), nil
	}); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func queryApps(search string, apps *[]AppJSON) *[]AppJSON {
	var out []AppJSON
	for _, a := range *apps {
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

//func unmarshalApp(reader *io.ReadCloser, app *AppJSON) {
//	body := getJsonMap(reader)
//	_ = json.Unmarshal([]byte(body["app"]), app)
//}
