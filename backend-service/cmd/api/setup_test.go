package main

import (
	"os"
	"testing"

	"github.com/GirishBhutiya/demo/backend-service/data"
)

var testApp Config

func TestMain(m *testing.M) {
	repo := data.NewPostgresTestRepository(nil)
	testApp.Repo = repo
	os.Exit(m.Run())
}
