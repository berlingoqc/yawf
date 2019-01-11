package main

import (
	"flag"
	"fmt"

	"github.com/berlingoqc/yawf/website"
)

func main() {

	var dir string

	flag.StringVar(&dir, "web_asset", ".", "the directory to serve file from. Default current working directory")
	flag.Parse()

	if dir == "." {
		//dir, _ = os.Getwd()
		dir = "/home/wq/go/src/github.com/berlingoqc/yawf/website/root"
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
