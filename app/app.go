package app

import (
	"fmt"
	"jar_tools/config"
	"jar_tools/consts"
	"jar_tools/utils/fileUtil"
	"jar_tools/utils/jarUtil"
	"jar_tools/utils/osUtil"
	"jar_tools/utils/scriptUtil"
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
		path := fileUtil.GetCurrentDirectory() + "\\" + f.JarPath
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Error: 找不到Jar文件，请检查: %s\n", path)
			return
		}
		// 执行 JAR 程序
		if err := jarUtil.RunJar(f.JarPath, f.Port); err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	if input == "2" {
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
		return
	}
	if input == "3" {
		if scriptUtil.CheckStartupScript() {
			scriptUtil.DeleteStartupScript()
		} else {
			osType := osUtil.GetOsType()
			if osType == consts.OsTypeWindows {
				scriptUtil.CreateWindowsStartupScript()
				return
			}
			if osType == consts.OsTypeLinux {
				scriptUtil.CreateLinuxStartupScript()
				return
			}

			fmt.Println("Error: 暂不支持该操作系统")
			os.Exit(1)
		}
		return
	}
	// 删除启动脚本
	if input == "4" {
		scriptUtil.DeleteStartupScript()
		return
	}
	fmt.Println("Error: 无效的输入")
	return
}
