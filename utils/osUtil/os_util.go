package osUtil

import (
	"fmt"
	"runtime"
)

// GetOsType 获取操作系统类型
func GetOsType() string {
	if runtime.GOOS == "windows" {
		return "windows"
	}
	if runtime.GOOS == "linux" {
		return "linux"
	}
	panic(fmt.Sprintf("不支持的操作系统: %s", runtime.GOOS))
}
