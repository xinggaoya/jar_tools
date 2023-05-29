package app

import (
	"fmt"
	"jar_tools/config"
	"jar_tools/utils/jarUtil"
	"os"
)

/**
  @author: XingGao
  @date: 2023/5/29
**/

// Start 程序处理
func Start(input string) {
	if input == "1" {
		f := config.GetConfig()
		if _, err := os.Stat(f.JarPath); os.IsNotExist(err) {
			fmt.Printf("Error: 找不到Jar文件，请检查: %s\n", f.JarPath)
			return
		}

		// 执行 JAR 程序
		if err := jarUtil.RunJar(f.JarPath, f.Port); err != nil {
			fmt.Println(err)
			return
		}
	} else if input == "2" {
		// 获取 JAR 文件路径和端口号
		f := config.GetConfig()
		pid := jarUtil.GetJarPidByPort(f.Port)
		if pid != 0 {
			err := jarUtil.KillProcess(pid)
			if err != nil {
				return
			} else {
				fmt.Printf("端口 %d ,进程 %d 已杀掉\n", f.Port, pid)
			}
		}
	} else if input == "3" {
		// 创建配置文件
		config.InitConfig()
		fmt.Println("配置文件初始化成功")
		return
	} else {
		fmt.Println("Error: 无效的输入")
		return
	}
}
