package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

func In(haystack interface{}, needle interface{}) bool {
	sVal := reflect.ValueOf(haystack)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == needle {
				return true
			}
		}
		return false
	}
	return false
}

func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	function := runtime.FuncForPC(pc[0])
	return function.Name()
}

func GetProjectPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	indexWithoutFileName := strings.LastIndex(path, string(os.PathSeparator))
	indexWithoutLastPath := strings.LastIndex(path[:indexWithoutFileName], string(os.PathSeparator))
	return strings.Replace(path[:indexWithoutLastPath], "\\", "/", -1)
}
