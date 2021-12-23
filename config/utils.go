package config

import (
	"os"
)

func GetWorkingDirectory() string {
	if wd, err := os.Getwd(); err != nil {
		return ""
	} else {
		return wd
	}
}
