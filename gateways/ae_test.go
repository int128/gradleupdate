package gateways

import (
	"fmt"
	"os"
	"testing"

	"github.com/favclip/testerator"
	_ "github.com/favclip/testerator/datastore"
)

func TestMain(m *testing.M) {
	_, _, err := testerator.SpinUp()
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	status := m.Run()

	err = testerator.SpinDown()
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}
