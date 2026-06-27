package test_prj

import (
	"os"
	"testing"

	"test_prj/client"
	"test_prj/config"
)

var salesServiceClient *client.SalesServiceClient

func TestMain(m *testing.M) {
	salesServiceClient = client.NewSalesClient(config.SalesServiceUrl)

	exitCode := m.Run()

	os.Exit(exitCode)
}
