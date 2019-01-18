package main

// #vagin

import (
	"flag"
	"fmt"

	"github.com/berlingoqc/yawf/website"
)

func main() {

	var dir string

	flag.StringVar(&dir, "config_file", "", "the directory to serve file from. Default current working directory")
	flag.Parse()

	if dir == "" {
		fmt.Println("Enter the configuration file with --config_file")
		return
	}

	fmt.Printf("Web asset %v \n", dir)

	// Enregistre les modules dynamique supplementaire requis

	ws, err := website.StartWebServer(dir)

	if err != nil {
		panic(err)
	}
	defer ws.Stop()
	for {
		select {
		case <-ws.ChannelStop:
			return
		default:
			ws.TaskPool.LaunchNeededTask()
		}
	}
}
