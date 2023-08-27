package testutil

import (
	"log"
	"os"
	"path"
	"runtime"
)

// https://brandur.org/fragments/testing-go-project-root
func init() {
	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled
	dir := path.Join(path.Dir(filename), "../../..")

	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
}
