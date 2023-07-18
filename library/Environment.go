package library

import (
	"os"
	"path"
	"runtime"
)

func GetPrefixEnvironment() string {
	env := os.Getenv("ENV")
	if env != "" {
		return env + "-"
	} else {
		return ""
	}
}
func GetEnvironment() string {
	return os.Getenv("ENV")
}

func RunOnRoot() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..") // change to suit test file location
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
