package app

import (
	"fmt"
	"jar_tools/config"
	"jar_tools/utils/jarUtil"
	"os"
	"os/exec"
	"runtime"
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
		// 设置开机自启
		var script string
		if runtime.GOOS == "windows" {
			script = "@echo off\nstart /b myprogram.exe"
			createStartupScript("C:\\Users\\10322\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup", script)
		} else if runtime.GOOS == "linux" {
			script = "#!/bin/bash\n./myprogram &"
			createStartupScript("$HOME/.config/autostart", script)
		} else {
			fmt.Println("Unsupported operating system")
			os.Exit(1)
		}
		return
	}
	fmt.Println("Error: 无效的输入")
	return
}

func createStartupScript(path, script string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%v/%v", path, "myprogram.desktop"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(script)
	if err != nil {
		return err
	}

	if runtime.GOOS == "linux" {
		cmd := exec.Command("chmod", "+x", fmt.Sprintf("%v/%v", path, "myprogram.desktop"))
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
