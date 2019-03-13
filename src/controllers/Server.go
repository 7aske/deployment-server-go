package controllers

import (
	"../config"
	"fmt"
	"net/http"
)

func NewServer() {
	cfg := config.LoadConfig()
	fmt.Println(fmt.Sprintf("port: %d", cfg.GetPort()))
	handler := NewHandler()
	handler.SetConfig(&cfg)
	http.HandleFunc("/auth", handler.HandleAuth)
	http.HandleFunc("/api/deploy", handler.HandleDeploy)
	http.HandleFunc("/api/update", handler.HandleUpdate)
	http.HandleFunc("/api/run", handler.HandleRun)
	http.HandleFunc("/api/kill", handler.HandleKill)
	http.HandleFunc("/api/remove", handler.HandleRemove)
	http.HandleFunc("/api/settings", handler.HandleSettings)
	http.HandleFunc("/", handler.HandleRoot)
	fmt.Printf("starting http server on port %d\n", cfg.GetPort())
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetPort()), nil)
	if err != nil {
		panic(fmt.Sprintf("error starting server on port %d", cfg.GetPort()))
	}
}