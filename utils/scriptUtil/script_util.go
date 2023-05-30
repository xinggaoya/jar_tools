package scriptUtil

import (
	"fmt"
	"jar_tools/config"
	"jar_tools/consts"
	"jar_tools/utils/osUtil"
	"os"
	"path/filepath"
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
	startupPath := fmt.Sprintf(`C:\Users\%s\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`, os.Getenv("USERNAME"))
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
	fmt.Println("Windows startup script created successfully!")
}

func CreateLinuxStartupScript() {
	// 脚本名称 时间戳+后缀
	scriptName := fmt.Sprintf("startup_%d.desktop", time.Now().Unix())
	exePath, _ := os.Executable()
	homeDir, _ := os.UserHomeDir()
	scriptPath := fmt.Sprintf("%s/.config/autostart/%s", homeDir, scriptName)
	// 设置启动参数
	startupParams := "-mode=1"

	scriptContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Exec=%s %s
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Name[en_US]=My App
Name=My App
Comment[en_US]=My App Startup Script
Comment=My App Startup Script`, exePath, startupParams)

	err := os.MkdirAll(filepath.Dir(scriptPath), 0755)
	if err != nil {
		fmt.Println("Failed to create startup directory:", err)
		return
	}

	err = os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		fmt.Println("Failed to create startup script:", err)
		return
	}
	config.SetConfig(&config.Config{
		ScriptName: scriptName,
	})
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
		startupPath := fmt.Sprintf(`C:\Users\%s\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`, os.Getenv("USERNAME"))
		scriptPath := fmt.Sprintf("%s\\%s", startupPath, f.ScriptName)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return false
		}
		return true
	}
	if osType == consts.OsTypeLinux {
		homeDir, _ := os.UserHomeDir()
		scriptPath := fmt.Sprintf("%s/.config/autostart/%s", homeDir, f.ScriptName)
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
		startupPath := fmt.Sprintf(`C:\Users\%s\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`, os.Getenv("USERNAME"))
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
		homeDir, _ := os.UserHomeDir()
		scriptPath := fmt.Sprintf("%s/.config/autostart/%s", homeDir, f.ScriptName)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			fmt.Println("Error: 未找到开机启动脚本")
			return
		}
		err := os.Remove(scriptPath)
		if err != nil {
			fmt.Println("Failed to delete startup script:", err)
			return
		}
		fmt.Println("Linux startup script deleted successfully!")
		return
	}
	fmt.Println("Error: 未知的操作系统类型")
}
