package config

import "testing"

const (
	Module = "blog.so"
)

func TestLoadModule(t *testing.T) {
	name, module, err := LoadModuleDynamicly(Module)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Module name %v info %v\n", name, module.GetInfo())
}
