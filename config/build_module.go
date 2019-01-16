package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GoBuildModule build the module of the package given to the given root
// and return the full information of the module
func GoBuildModule(modulepkg string, subpkg string, dstPath string) (*BuildModule, error) {
	// get go path
	gopath := os.Getenv("GOPATH")
	bm := &BuildModule{}
	// for the name take the last after split of package
	pkg := modulepkg + "/" + subpkg
	lp := strings.Split(pkg, "/")
	if i := len(lp); i > 0 {
		bm.Name = lp[i-1]
	} else {
		return nil, errors.New("Invalide pkg " + pkg)
	}

	pkg += "/plugin"

	bm.Path = fmt.Sprintf("%v/%v.so", dstPath, bm.Name)
	fmt.Printf("Building module %v at %v output %v \n", bm.Name, pkg, bm.Path)
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", bm.Path, pkg)
	cmd.Dir = gopath + "/src/" + modulepkg
	var b bytes.Buffer
	cmd.Stdout = &b

	err := cmd.Run()
	print(b.String())
	if err != nil {
		return nil, err
	}

	return bm, nil
}
