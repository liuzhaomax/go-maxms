package script

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

// url 和 path 的声明放在测试函数外面的原因：
// 因为 flag 包在整个程序运行期间只能被初始化一次，
// 而 TestGetContract 是一个测试函数，相当于整个程序的一部分。
// 在 Go 语言中，flag 包是在程序启动时解析命令行参数的，
// 而 TestGetContract 是在测试运行时执行的。
// 因此，如果你在测试函数内部声明 flag，它不会被正确地初始化。
var (
	url  = flag.String("url", "", "contract的URL")
	path = flag.String("path", "", "contract的保存位置和文件名")
)

func TestGetContract(t *testing.T) {
	flag.Parse()
	fmt.Println("path", *path)
	// 发送请求
	err := GetContract(*url, *path)
	if err != nil {
		t.Errorf("GetContract failed: %v", err)
	}
	// 检查文件是否存在
	_, err = os.Stat(*path)
	if os.IsNotExist(err) {
		t.Errorf("File was not created: %v", err)
	}
}
