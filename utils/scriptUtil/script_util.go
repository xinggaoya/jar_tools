package scriptUtil

import (
	"fmt"
	"jar_tools/config"
	"jar_tools/consts"
	"jar_tools/utils/osUtil"
	"log"
	"os"
	"os/exec"
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
	scriptName := fmt.Sprintf("startup_%d", time.Now().Unix())
	// 获取可执行文件路径
	exePath, _ := os.Executable()
	// 脚本路径
	scriptPath := fmt.Sprintf("/etc/systemd/system/%s.service", scriptName)
	// 脚本内容
	scriptContent := fmt.Sprintf(`[Unit]
Description=jar_tools startup script
After=network.target

[Service]
Type=simple
ExecStart=%s -mode=1
Restart=always
RestartSec=1

[Install]
WantedBy=multi-user.target`, exePath)

	// 创建脚本文件夹
	err := os.MkdirAll("/etc/systemd/system", 0755)

	// 写入脚本文件
	err = os.WriteFile(scriptPath, []byte(scriptContent), 0644)
	if err != nil {
		log.Println("Failed to create startup script:", err)
	}
	// 重新加载配置
	cmd := exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	// 开启服务
	cmd = exec.Command("systemctl", "enable", scriptName)
	err = cmd.Run()

	fmt.Println("Linux startup script created successfully!")
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
		scriptPath := fmt.Sprintf("/etc/systemd/system/%s.service", f.ScriptName)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return false
		}
		return true
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
		scriptPath := fmt.Sprintf("/etc/systemd/system/%s.service", f.ScriptName)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			fmt.Println("Error: 未找到开机启动脚本")
			return
		}
		err := os.Remove(scriptPath)
		// 重新加载配置
		cmd := exec.Command("systemctl", "daemon-reload")
		err = cmd.Run()
		// 关闭服务
		cmd = exec.Command("systemctl", "disable", f.ScriptName)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Failed to delete startup script:", err)
			return
		}
		fmt.Println("Linux startup script deleted successfully!")
		return
	}
	fmt.Println("Error: 未知的操作系统类型")
}
