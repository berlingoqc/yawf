package conv

import (
	"os"
	"testing"
)

type TestStruct struct {
	About bool
	Name  string
	Value int
}

type TestStruct2 struct {
	Name string
	Root string
	V    int
	V2   int
}

const (
	FileJson = "data.json"
)

func TestConversion(t *testing.T) {
	defer func() {
		os.Remove(FileJson)
	}()
	oi := &TestStruct{
		About: true,
		Name:  "1995/07/05",
		Value: 1,
	}

	ts2 := &TestStruct2{
		Name: "dasdasdas",
		Root: "/dsadas/csada/sadas",
		V:    32312312,
		V2:   23111,
	}

	bigMap := make(map[string]interface{})

	// Ajoute ma struct dans ma map

	err := AddStructToMap(bigMap, oi)
	if err != nil {
		t.Fatal(err)
	}
	err = AddStructToMap(bigMap, oi)
	if err == nil {
		t.Fatal("Should throw key already set error")
	}

	err = AddStructToMap(bigMap, ts2)
	if err != nil {
		t.Fatal(err)
	}

	oi_2 := &TestStruct{}
	// Recupere ma struct depuis ma map
	err = FindStructMap(bigMap, oi_2)
	if err != nil {
		t.Fatal(err)
	}
	ts2_r := &TestStruct2{}
	err = FindStructMap(bigMap, ts2_r)
	if err != nil {
		t.Fatal(err)
	}

	// Enregistre ma map dans un fichier
	err = Save(FileJson, bigMap)
	if err != nil {
		t.Fatal(err)
	}

	// Load et get ma struct depuis fichier
	bigMapL, err := Load(FileJson)
	if err != nil {
		t.Fatal(err)
	}
	if len(bigMapL) == 0 {
		t.Fatal("No data in the map")
	}

}
