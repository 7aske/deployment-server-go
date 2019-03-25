package main

import (
	"./controllers"
	"./utils"
	"os"
)

func main() {
	if utils.Contains("-h", &os.Args) != -1 || utils.Contains("-help", &os.Args) != -1 || utils.Contains("--help", &os.Args) != -1 || utils.Contains("help", &os.Args) != -1 {
		utils.PrintHelp()
	}
	if utils.Contains("--client", &os.Args) != -1 {

	}
	if utils.Contains("-i", &os.Args) != -1 {

	}
	controllers.NewServer()
}
