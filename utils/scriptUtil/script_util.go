package scriptUtil

import (
	"fmt"
	"jar_tools/config"
	"jar_tools/consts"
	"jar_tools/utils/osUtil"
	"os"
	"os/exec"
	"os/user"
	"time"
)

/**
  @author: XingGao
  @date: 2023/5/30
**/

func CreateWindowsStartupScript() {
	// 脚本名称 时间戳+后缀
	scriptName := fmt.Sprintf("startup_%d.bat", time.Now().Unix())
	exePath, _ := os.Executable()
	userHomeDir, _ := os.UserHomeDir()
	startupPath := fmt.Sprintf(`%s\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`, userHomeDir)
	scriptPath := fmt.Sprintf("%s\\%s", startupPath, scriptName)

	scriptContent := fmt.Sprintf(`@echo off
start "" "%s" -mode=1
exit`, exePath)

	err := os.MkdirAll(startupPath, 0755)
	if err != nil {
		fmt.Println("Failed to create startup directory:", err)
		return
	}

	err = os.WriteFile(scriptPath, []byte(scriptContent), 0644)
	if err != nil {
		fmt.Println("Failed to create startup script:", err)
		return
	}
	config.SetConfig(&config.Config{
		ScriptName: scriptName,
	})
	fmt.Println("开机启动脚本创建成功")
}

// CreateLinuxStartupScript 创建Linux开机启动脚本
func CreateLinuxStartupScript() {
	// 脚本名称 时间戳+后缀
	serviceName := fmt.Sprintf("startup_%d", time.Now().Unix())
	// 获取当前可执行文件的路径
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("无法获取可执行文件路径：", err)
		return
	}

	// 创建Shell脚本内容
	scriptContent := `#!/bin/bash
# 运行当前可执行文件
%s -mode=1
`

	// 设置脚本文件路径
	scriptPath := fmt.Sprintf("/etc/init.d/%s", serviceName)

	// 将可执行文件路径插入到脚本内容中
	scriptContent = fmt.Sprintf(scriptContent, executablePath)

	// 写入脚本内容到文件
	err = os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		fmt.Println("无法写入脚本文件：", err)
		return
	}
	currentUser, err := user.Current()
	// 检查是否管理员权限
	if err != nil || currentUser.Uid != "0" {
		fmt.Println("请使用管理员权限运行此程序")
		return
	}

	// 使用update-rc.d命令将脚本添加到启动项
	cmd := exec.Command("update-rc.d", serviceName, "defaults")
	err = cmd.Run()
	if err != nil {
		fmt.Println("无法设置开机自启：", err)
		return
	}

	// 设置配置
	config.SetConfig(&config.Config{
		ScriptName: serviceName,
	})

	// 输出成功消息
	fmt.Println("已设置开机自启")
}

// CheckStartupScript 检查开机启动脚本是否存在
func CheckStartupScript() bool {
	// 获取配置
	f := config.GetConfig()
	if f.ScriptName == "" {
		return false
	}
	// 获取操作系统类型
	osType := osUtil.GetOsType()
	if osType == consts.OsTypeWindows {
		userHome, _ := os.UserHomeDir()
		startupPath := fmt.Sprintf(`%s\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`, userHome)
		scriptPath := fmt.Sprintf("%s\\%s", startupPath, f.ScriptName)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return false
		}
		return true
	}
	if osType == consts.OsTypeLinux {
		scriptPath := fmt.Sprintf("/etc/init.d/%s", f.ScriptName)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return false
		}
	}
	return false
}

// DeleteStartupScript 删除开机启动脚本
func DeleteStartupScript() {
	// 获取配置
	f := config.GetConfig()
	if f.ScriptName == "" {
		fmt.Println("Error: 未找到开机启动脚本")
		return
	}
	// 获取操作系统类型
	osType := osUtil.GetOsType()
	if osType == consts.OsTypeWindows {
		userHome, _ := os.UserHomeDir()
		startupPath := fmt.Sprintf(`%s\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`, userHome)
		scriptPath := fmt.Sprintf("%s\\%s", startupPath, f.ScriptName)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			fmt.Println("Error: 未找到开机启动脚本")
			return
		}
		err := os.Remove(scriptPath)
		if err != nil {
			fmt.Println("Failed to delete startup script:", err)
			return
		}
		config.GetConfig()
		fmt.Println("Windows startup script deleted successfully!")
		return
	}
	if osType == consts.OsTypeLinux {
		scriptPath := fmt.Sprintf("/etc/init.d/%s", f.ScriptName)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			fmt.Println("Error: 未找到开机启动脚本")
			return
		}
		err := os.Remove(scriptPath)
		if err != nil {
			fmt.Println("Failed to delete startup script:", err)
			return
		}
		// 检查是否管理员权限
		currentUser, err := user.Current()
		if err != nil || currentUser.Uid != "0" {
			fmt.Println("请使用管理员权限运行此程序")
			return
		}
		cmd := exec.Command("update-rc.d", f.ScriptName, "remove")
		err = cmd.Run()
		if err != nil {
			fmt.Println("Failed to delete startup script:", err)
			return
		}
		config.GetConfig()
		fmt.Println("Linux startup script deleted successfully!")
	}
}
