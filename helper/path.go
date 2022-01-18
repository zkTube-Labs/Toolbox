package helper

import (
	"os"
)

func GetProjectRoot() (path string, err error) {
	path, err = os.Getwd()
	return
}
