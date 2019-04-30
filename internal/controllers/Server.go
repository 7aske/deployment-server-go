package controllers

import (
	"../config"
	"../logger"
	"../utils"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func NewServer() {
	l := logger.NewLogger(logger.LOG_SERVER)
	cfg := config.LoadConfig()
	deployer := NewDeployer(cfg)
	handler := NewHandler(cfg, &deployer)
	cli := NewCli(&deployer)
	routerHandler := NewRouterHandler(&deployer, cfg)
	devMux := http.NewServeMux()

	devMux.HandleFunc("/api/deploy", func(writer http.ResponseWriter, request *http.Request) {
		go func() {
			routerHandler.UpdateHosts()
			l.Log("updating router hosts")
			for key, value := range *routerHandler.GetHosts() {
				l.Log(fmt.Sprintf("%s %s", key, value))
			}
		}()
		handler.HandleDeploy(writer, request)
	})

	devMux.HandleFunc("/api/update", handler.HandleUpdate)
	devMux.HandleFunc("/api/run", handler.HandleRun)
	devMux.HandleFunc("/api/find", handler.HandleFind)
	devMux.HandleFunc("/api/kill", handler.HandleKill)
	devMux.HandleFunc("/api/remove", handler.HandleRemove)
	devMux.HandleFunc("/api/settings", handler.HandleSettings)
	devMux.HandleFunc("/auth", handler.HandleAuth)
	devMux.HandleFunc("/", handler.HandleRoot)
	routerMux := http.NewServeMux()
	//routerMux.HandleFunc("/", routerHandler.HandleRoot)
	routerMux.HandleFunc("/", routerHandler.HandleIndex)

	go func() {
		l.Log(fmt.Sprintf("starting dev    server on port %d", cfg.GetPort()))
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetPort()), devMux)
		if err != nil {
			l.Log(fmt.Sprintf("error starting server on port %d", cfg.GetPort()))
			panic(fmt.Sprintf("error starting server on port %d", cfg.GetPort()))
		}
	}()
	go func() {
		l.Log(fmt.Sprintf("starting router server on port %d", cfg.GetRouterPort()))
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetRouterPort()), routerMux)
		if err != nil {
			l.Log(fmt.Sprintf("error starting server on port %d", cfg.GetRouterPort()))
			panic(fmt.Sprintf("error starting server on port %d", cfg.GetRouterPort()))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	go func() {
		for sig := range c {
			deployer.GetLogger().Log("server killed with sig " + sig.String())
			for _, app := range *deployer.GetApps() {
				err := deployer.Kill(app)
				if err != nil {
					deployer.GetLogger().Log(err.Error())
				}
			}
			os.Exit(0)
		}
	}()

	if utils.Contains("-i", &os.Args) != -1 {
		fmt.Println("type \"help\" from help...")
		reader := bufio.NewReader(os.Stdin)
		for {
			line, _, _ := reader.ReadLine()
			args := strings.Split(string(line), " ", )
			cli.ParseCommand(args...)
			cli.lastCommand = line
		}

	} else {
		for {
			time.Sleep(time.Second)
		}
	}
}
