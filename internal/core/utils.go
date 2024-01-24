package core

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
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

func GetCallerName(level int) string {
	_, funcName, _ := GetCallerInfo(level)
	return funcName
}

func GetCallerFileAndLine(level int) string {
	file, _, line := GetCallerInfo(level)
	return fmt.Sprintf("\033[1;34m%s:%d\033[0m\n", file, line)
}

func GetCallerInfo(level int) (string, string, int) {
	pc, file, line, ok := runtime.Caller(level) // 读取第N层调用堆栈
	if !ok {
		return "", "", 0
	}
	// 通过函数的PC获取函数名
	functionName := runtime.FuncForPC(pc).Name()
	return file, functionName, line
}

func GetProjectPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	indexWithoutFileName := strings.LastIndex(path, string(os.PathSeparator))
	indexWithoutLastPath := strings.LastIndex(path[:indexWithoutFileName], string(os.PathSeparator))
	return strings.ReplaceAll(path[:indexWithoutLastPath], "\\", "/")
}

func GetRandomIdlePort() string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	return strconv.Itoa(port)
}
