package controllers

import (
	"../config"
	"../utils"
	json2 "encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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
	deployer                  Deployer
	statusInternalServerError []byte
	statusOK                  []byte
	statusUnauthorized        []byte
	statusNotFound            []byte
	statusMethodNotAllowed    []byte
}

type DeployRequest struct {
	Token string `json:"token"`
	Repo  string `json:"repo"`
}
type FindResponse struct {
	Running  []AppJSON  `json:"running"`
	Deployed *[]AppJSON `json:"deployed"`
}

//func (h Handler) LoadConfig() {
//	h.config = config.LoadConfig()
//
//}
func (h *Handler) SetConfig(c *config.Config) {
	h.config = c
}

func (h *Handler) GetConfig() *config.Config {
	return h.config
}
func (h *Handler) SetDeployer(d *Deployer) {
	h.deployer = *d
}

func (h *Handler) GetDeployer() *Deployer {
	return &h.deployer
}

func (h *Handler) HandleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		jsonBody := utils.GetJsonMap(r.Body)
		//name := utils.GetNameFromRepo(jsonBody["repo"])
		app, err := h.GetDeployer().Deploy(jsonBody["repo"], jsonBody["runner"])
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			length, _ := w.Write(h.statusInternalServerError)
			w.Header().Set("Content-Length", strconv.Itoa(length))
			return
		}
		err = h.GetDeployer().Install(app)
		app.Print()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			length, _ := w.Write(h.statusInternalServerError)
			w.Header().Set("Content-Length", strconv.Itoa(length))
			return
		}
		//err = h.GetDeployer().Run(app)
		//if err != nil {
		//	fmt.Println(err)
		//	w.WriteHeader(http.StatusInternalServerError)
		//	length, _ := w.Write(h.statusInternalServerError)
		//	w.Header().Set("Content-Length", strconv.Itoa(length))
		//	return
		//}
		w.WriteHeader(http.StatusOK)
		length, _ := w.Write(h.statusOK)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		length, _ := w.Write(h.statusMethodNotAllowed)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	}
}
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
}
func (h *Handler) HandleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body := utils.GetJsonMap(r.Body)
		name := body["app"]
		if appJson, ok := h.GetDeployer().GetAppD(name); ok {
			app := NewAppFromJson(appJson)
			if h.GetDeployer().IsAppRunning(app) {
				w.WriteHeader(http.StatusInternalServerError)
				length, _ := w.Write(h.statusInternalServerError)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			} else {
				err := h.GetDeployer().Run(app)
				if err != nil {
					fmt.Println(err)
				}
				w.WriteHeader(http.StatusOK)
				length, _ := w.Write(h.statusOK)
				w.Header().Set("Content-Length", strconv.Itoa(length))
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			length, _ := w.Write(h.statusInternalServerError)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		length, _ := w.Write(h.statusMethodNotAllowed)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	}
}
func (h *Handler) HandleFind(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Println(h.GetDeployer().GetApps())
		apps := h.GetDeployer().GetAppsAsJSON()
		appD := h.GetDeployer().GetDeployedApps()
		json, _ := json2.Marshal(&FindResponse{Running: apps, Deployed: &appD})
		length, _ := w.Write(json)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		length, _ := w.Write(h.statusMethodNotAllowed)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	}
}
func (h *Handler) HandleKill(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body := utils.GetJsonMap(r.Body)
		name := body["app"]
		if app, ok := h.GetDeployer().GetApp(name); ok {
			h.GetDeployer().Kill(app)
			w.WriteHeader(http.StatusOK)
			length, _ := w.Write(h.statusOK)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			length, _ := w.Write(h.statusInternalServerError)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		length, _ := w.Write(h.statusMethodNotAllowed)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	}

}
func (h *Handler) HandleRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body := utils.GetJsonMap(r.Body)
		name := body["app"]
		if appJson, ok := h.GetDeployer().GetAppD(name); ok {
			if app, ok := h.GetDeployer().GetApp(appJson.Id); ok {
				h.GetDeployer().Kill(app)
			}
			h.GetDeployer().Remove(appJson)
			w.WriteHeader(http.StatusOK)
			length, _ := w.Write(h.statusOK)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			length, _ := w.Write(h.statusInternalServerError)
			w.Header().Set("Content-Length", strconv.Itoa(length))
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		length, _ := w.Write(h.statusMethodNotAllowed)
		w.Header().Set("Content-Length", strconv.Itoa(length))
	}

}
func (h *Handler) HandleSettings(w http.ResponseWriter, r *http.Request) {
}
func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("Authorization"); err != nil {
		http.Redirect(w, r, "/auth", 301)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if h.verifyToken(token) {
			var absPth string
			pth := strings.Replace(r.URL.String(), "/", string(filepath.Separator), -1)
			root := "client/dist"
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
			http.Redirect(w, r, "/auth", 301)
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
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", 301)
		} else {
			http.Redirect(w, r, "/auth", 301)
		}
		break
	case http.MethodGet:
		if cookie, err := r.Cookie("Authorization"); err == nil {
			token := strings.Split(cookie.Value, "Bearer ")[1]
			if h.verifyToken(token) {
				http.Redirect(w, r, "/", 301)
				break
			}
		}
		loginPage := utils.RenderLoginPage()
		w.WriteHeader(http.StatusOK)
		length, _ := w.Write(loginPage)
		w.Header().Set("Content-Length", strconv.Itoa(length))
		break
	}
}

func (h *Handler) makeToken() string {
	expires := time.Now().Unix() + int64(24*time.Hour)
	type JSTClaims struct {
		Data string `json:"data"`
		jwt.StandardClaims
	}
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
	return configPassHash == passHash && configUser == user
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

func NewHandler(cfg *config.Config) Handler {
	h := Handler{}
	h.config = cfg
	h.statusInternalServerError = []byte("( ͠° ͟ʖ ͡°) 500 INTERNAL SERVER ERROR")
	h.statusNotFound = []byte("( ͡° ʖ̯ ͡°) 404 NOT FOUND")
	h.statusUnauthorized = []byte("( ͠° ͟ʖ ͡°) 401 UNAUTHORIZED")
	h.statusMethodNotAllowed = []byte("( ͠° ͟ʖ ͡°) 405 METHOD NOT ALLOWED")
	h.statusOK = []byte("( ͡ᵔ ͜ʖ ͡ᵔ ) 200 OK")
	return h
}
