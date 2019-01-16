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

var (
	oi = &TestStruct{
		About: true,
		Name:  "1995/07/05",
		Value: 1,
	}

	ts2 = &TestStruct2{
		Name: "dasdasdas",
		Root: "/dsadas/csada/sadas",
		V:    32312312,
		V2:   23111,
	}
)

func TestStructQuery(t *testing.T) {
	query, err := StructToQuery(ts2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Query is /url%v\n", query)

	queryValues := map[string][]string{
		"Name": []string{"dasdasdas"},
		"Root": []string{"/dsadas/csada/sadas"},
		"V":    []string{"323231"},
		"V2":   []string{"113"},
	}

	ts := &TestStruct2{}
	err = QueryToStruct(queryValues, ts)
	if err != nil {
		t.Fatal(err)
	}

}

func TestConversion(t *testing.T) {
	defer func() {
		os.Remove(FileJson)
	}()
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

	oi2 := &TestStruct{}
	// Recupere ma struct depuis ma map
	err = FindStructMap(bigMap, oi2)
	if err != nil {
		t.Fatal(err)
	}
	ts2r := &TestStruct2{}
	err = FindStructMap(bigMap, ts2r)
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
