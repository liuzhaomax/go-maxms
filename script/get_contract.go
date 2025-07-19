package script

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetContract(url string, localFilePath string) error {
	// 发送HTTP GET请求
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP请求失败：%v", err)
	}
	defer response.Body.Close()
	// 创建本地文件
	file, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败：%v", err)
	}
	defer file.Close()
	// 将响应体拷贝到本地文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("拷贝响应体到本地文件失败：%v", err)
	}
	fmt.Printf("代码已保存到本地文件：%s\n", localFilePath)
	return nil
}
