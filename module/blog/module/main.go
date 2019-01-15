package main

import (
	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/module/blog"
)

// GetModule return the blog module
func GetModule() (string, config.IModule) {

	return "blog", &blog.Module{}
}

func main() {

}
