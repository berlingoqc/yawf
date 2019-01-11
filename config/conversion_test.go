package config

import (
	"testing"
)

func TestConversion(t *testing.T) {
	oi := &OwnerInformation{
		About:        true,
		FAQ:          true,
		PersonalPage: true,
		ContactPage:  true,
		Birth:        "1995/07/05",
		FullName:     "William Quintal",
		Location:     "Quebec",
		ThumbnailURL: "/static/img/me.png",
	}

	n, mapOI, err := StructToMap(oi)
	if err != nil {
		t.Fatal(err)
	}

	bigMap := make(map[string]interface{})
	bigMap[n] = mapOI

	err = MapToStruct(bigMap, oi)
	if err != nil {
		t.Fatal(err)
	}

}
