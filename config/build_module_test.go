package config

import "testing"

func TestBuildModule(t *testing.T) {
	bm, err := GoBuildModule("github.com/berlingoqc/yawf", "module/blog", "/home/wq/go/src/github.com/berlingoqc/yawf/website/root/module")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(*bm)
}
