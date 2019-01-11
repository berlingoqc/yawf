package project

import (
	"fmt"

	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/db"
)

func TaskGitHubUpdate(c chan *config.Signal, args ...interface{}) {
	// Get instance bd
	idb, err := GetProjectDBInstance("project.db")
	if err != nil {
		// log l'erreur
		fmt.Printf("Error task github %v \n", err)
		return
	}
	defer db.CloseDatabse(idb)

	// Get les nouvelles shits
	ua, err := UpdateAccountInfo("berlingoqc")
	if err != nil {
		fmt.Printf("Error task github %v \n", err)
		return
	}
	idb.UpdateGHAccount(ua)
	repos, err := UpdateRepositoryInfo("berlingoqc", "YASE", "yawf")
	if err != nil {
		fmt.Printf("Error task github %v \n", err)
		return
	}
	// update dans la bd
	for _, r := range repos {
		idb.UpdateGitHubRepo(r)
	}
}
