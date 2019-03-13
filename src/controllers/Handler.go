package controllers

import (
	"../config"
	"../utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	config   config.Config
	deployer Deployer
}

func (h Handler) LoadConfig() {
	h.config = config.LoadConfig()

}
func (h *Handler) SetConfig(c *config.Config) {
	h.config = *c
}

func (h *Handler) GetConfig() *config.Config {
	return &h.config
}

func (h *Handler) HandleDeploy(w http.ResponseWriter, r *http.Request) {
}
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
}
func (h *Handler) HandleRun(w http.ResponseWriter, r *http.Request) {
}
func (h *Handler) HandleKill(w http.ResponseWriter, r *http.Request) {
}
func (h *Handler) HandleRemove(w http.ResponseWriter, r *http.Request) {
}
func (h *Handler) HandleSettings(w http.ResponseWriter, r *http.Request) {
}
func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("Authorization"); err != nil {
		//w.WriteHeader(http.StatusUnauthorized)
		//_, _ = fmt.Fprint(w, "( ͠° ͟ʖ ͡°) 401 Unauthorized")
		http.Redirect(w, r, "/auth", 301)
	} else {
		token := strings.Split(cookie.Value, "Bearer ")[1]
		if h.verifyToken(token) {
			w.WriteHeader(http.StatusOK)
			length, _ := w.Write([]byte("( ͡ᵔ ͜ʖ ͡ᵔ ) 200 OK"))
			w.Header().Set("Content-Length", strconv.Itoa(length))
		} else {
			//w.WriteHeader(http.StatusUnauthorized)
			//_, _ = fmt.Fprint(w, "( ͠° ͟ʖ ͡°) 401 Unauthorized")
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
			//type Response struct {
			//	Token string `json:"token"`
			//}
			//res := Response{Token: token}
			//jsonRes, err := json.Marshal(res)
			//if err != nil {
			//	fmt.Println(err)
			//}
			//w.Header().Set("Content-Type", "application/json")
			//_, _ = w.Write(jsonRes)
			http.Redirect(w, r, "/", 301)
		} else {
			//w.WriteHeader(http.StatusUnauthorized)
			//_, _ = fmt.Fprint(w, "( ͠° ͟ʖ ͡°) 401 Unauthorized")
			http.Redirect(w, r, "/auth", 301)
		}
		break
	case http.MethodGet:
		if cookie, err := r.Cookie("Authorization"); err == nil {
			token := strings.Split(cookie.Value, "Bearer ")[1]
			if h.verifyToken(token) {
				http.Redirect(w,r,"/", 301)
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

func NewHandler() Handler {
	return Handler{}
}
