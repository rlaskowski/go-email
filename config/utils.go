package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var executableDir, executableName string

func init() {

	initFlags()

	initExecutable()

}

func initExecutable() {
	if executable, err := os.Executable(); err != nil {
		executableDir = "./"

		switch runtime.GOOS {
		case "windows":
			executableName = ExecutableName + ".exe"
		default:
			executableName = ExecutableName
		}

	} else {
		executableName = filepath.Base(executable)
		executableDir = filepath.Dir(executable)
	}
}

func GetWorkingDirectory() string {
	if wd, err := os.Getwd(); err != nil {
		return executableDir
	} else {
		return wd
	}
}

func GetExecutableDirectory() string {
	return executableDir
}
