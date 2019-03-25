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
	fmt.Println(fmt.Sprintf("port: %d", cfg.GetPort()))
	handler := NewHandler(cfg)
	deployer := NewDeployer(cfg)
	handler.SetDeployer(&deployer)
	cli := NewCli(&deployer)
	http.HandleFunc("/auth", handler.HandleAuth)
	http.HandleFunc("/api/deploy", handler.HandleDeploy)
	http.HandleFunc("/api/update", handler.HandleUpdate)
	http.HandleFunc("/api/run", handler.HandleRun)
	http.HandleFunc("/api/find", handler.HandleFind)
	http.HandleFunc("/api/kill", handler.HandleKill)
	http.HandleFunc("/api/remove", handler.HandleRemove)
	http.HandleFunc("/api/settings", handler.HandleSettings)
	http.HandleFunc("/", handler.HandleRoot)
	fmt.Printf("starting http server on port %d\n", cfg.GetPort())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetPort()), nil)
		if err != nil {
			panic(fmt.Sprintf("error starting server on port %d", cfg.GetPort()))
		}
	}()
	fmt.Println("type \"help\" from help...")
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := reader.ReadLine()
		args := strings.Split(string(line), " ", )
		//if args[0] != "!" {
		cli.ParseCommand(args...)
		cli.lastCommand = line
		//} else {
		//	fmt.Println("lc")
		//	cli.PutLastCommand()
		//}
	}

}
