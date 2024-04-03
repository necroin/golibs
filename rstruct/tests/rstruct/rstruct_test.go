package tests

import (
	"fmt"
	"testing"

	"github.com/necroin/golibs/rstruct"
)

func TestMain(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.AddFields(
		rstruct.NewRTField("Name", "").
			SetTag("json", "name"),
		rstruct.NewRTField("Age", 10).
			SetTag("json", "age"),
	)
	if err != nil {
		t.Fatal(err)
	}

	instance := customStruct.New()
	fmt.Println("String:", instance)
	jsonInstance, _ := instance.ToJSON()
	fmt.Println("JSON:", string(jsonInstance))

}
