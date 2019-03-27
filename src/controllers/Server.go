package controllers

import (
	"../config"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func NewServer() {
	cfg := config.LoadConfig()
	deployer := NewDeployer(cfg)
	handler := NewHandler(cfg, &deployer)
	cli := NewCli(&deployer)
	routerHandler := NewRouterHandler(&deployer, cfg)
	devMux := http.NewServeMux()
	devMux.HandleFunc("/auth", handler.HandleAuth)
	devMux.HandleFunc("/api/deploy", func(writer http.ResponseWriter, request *http.Request) {
		handler.HandleDeploy(writer, request)
		routerHandler.UpdateHosts()
		fmt.Println(*routerHandler.GetHosts())
	})
	devMux.HandleFunc("/api/update", handler.HandleUpdate)
	devMux.HandleFunc("/api/run", handler.HandleRun)
	devMux.HandleFunc("/api/find", handler.HandleFind)
	devMux.HandleFunc("/api/kill", handler.HandleKill)
	devMux.HandleFunc("/api/remove", handler.HandleRemove)
	devMux.HandleFunc("/api/settings", handler.HandleSettings)
	devMux.HandleFunc("/", handler.HandleRoot)
	routerMux := http.NewServeMux()
	routerMux.HandleFunc("/", routerHandler.HandleRoot)
	go func() {
		fmt.Printf("starting dev server on port %d\n", cfg.GetPort())
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetPort()), devMux)
		if err != nil {
			panic(fmt.Sprintf("error starting server on port %d", cfg.GetPort()))
		}
	}()
	go func() {
		fmt.Printf("starting router server on port %d\n", cfg.GetRouterPort())
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetRouterPort()), routerMux)
		if err != nil {
			panic(fmt.Sprintf("error starting server on port %d", cfg.GetRouterPort()))
		}
	}()

	fmt.Println("type \"help\" from help...")
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := reader.ReadLine()
		args := strings.Split(string(line), " ", )
		cli.ParseCommand(args...)
		cli.lastCommand = line
	}

}
