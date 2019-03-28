package main

import (
	"../../internal/controllers"
	"../../internal/utils"
	"os"
)

func main() {
	if utils.Contains("-h", &os.Args) != -1 || utils.Contains("-help", &os.Args) != -1 || utils.Contains("--help", &os.Args) != -1 || utils.Contains("help", &os.Args) != -1 {
		utils.PrintHelp()
	}
	controllers.NewServer()
}
