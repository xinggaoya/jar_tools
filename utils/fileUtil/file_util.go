package fileUtil

import (
	"fmt"
	"os"
	"path/filepath"
)

/**
  @author: XingGao
  @date: 2023/5/30
**/

// GetCurrentDirectory 获取当前程序目录
func GetCurrentDirectory() string {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to get executable path:", err)
		return ""
	}

	dir := filepath.Dir(exePath)
	return dir
}
