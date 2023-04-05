package main

import (
	"fmt"
	"jar_tools/consts"
	"jar_tools/file"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	// 作者信息
	fmt.Printf("%s: %s By 2023\n\n", consts.AppName, consts.AppAuthor)
	// 选择操作
	fmt.Printf("请选择操作：\n1. 启动JAR程序\n2. 停止JAR程序\n3. 初始化配置\n")
	fmt.Print("请输入操作序号:")
	var input string
	_, err := fmt.Scanf("%s", &input)
	if err != nil {
		return
	}
	if input == "1" {
		f, err := file.GetConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if _, err := os.Stat(f.JarPath); os.IsNotExist(err) {
			fmt.Printf("Error: 找不到Jar文件，请检查: %s\n", f.JarPath)
			return
		}

		// 执行 JAR 程序
		if err := RunJar(f.JarPath, f.Port); err != nil {
			fmt.Println(err)
			return
		}
	} else if input == "2" {
		// 获取 JAR 文件路径和端口号
		f, err := file.GetConfig()
		if err != nil {
			fmt.Println("Error: 找不到PID文件，请检查")
			return
		}
		if f.Pid != 0 {
			err := killProcess(f.Pid)
			if err != nil {
				return
			} else {
				fmt.Printf("进程 %d 已杀掉\n", f.Pid)
			}
		}
	} else if input == "3" {
		// 创建配置文件
		err := file.InitConfig()
		if err != nil {
			return
		}
		fmt.Println("配置文件初始化成功")
		return
	} else {
		fmt.Println("Error: 无效的输入")
		return
	}
}

func RunJar(jarPath string, port int) error {
	if pid, err := FindProcessByPort(port); err == nil {
		fmt.Printf("Error: 端口 %d 已被占用，进程ID为 %d\n", port, pid)
		fmt.Printf("是否杀掉进程 %d (y/n)？\n", pid)
		var input string
		_, err = fmt.Scanf("\n%s", &input)
		if input == "y" || input == "Y" {
			if err := killProcess(pid); err != nil {
				return err
			}
			fmt.Printf("进程 %d 已杀掉\n", pid)
		} else {
			return fmt.Errorf("端口 %d 已被占用，操作取消", port)
		}
	}
	f, err := file.GetConfig()
	cmd := exec.Command("javaw", "-jar", jarPath, f.Jvm)
	err = cmd.Start()
	if err != nil {
		return err
	}
	pid, err := cmd.Process.Pid, cmd.Process.Release()
	if err != nil {
		return err
	}
	fmt.Printf("JAR程序已启动，进程ID为 %d\n", pid)
	// 写入 PID 文件
	config := file.NewConfig()
	config.Pid = pid
	err = file.SetConfig(config)
	if err != nil {
		return fmt.Errorf("写入 PID 文件失败: %s", err)
	}
	return nil
}

func FindProcessByPort(port int) (int, error) {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("netstat", "-aon")
		output, err := cmd.Output()
		if err != nil {
			return 0, err
		}
		// 获取输出结果中的所有行
		for _, line := range strings.Split(string(output), "\n") {
			// 按空格拆分每行，获取端口和进程ID
			fields := strings.Fields(line)
			if len(fields) >= 4 && strings.Contains(fields[2], fmt.Sprintf(":%d", port)) {
				pid, err := strconv.Atoi(strings.TrimSuffix(fields[4], "\r"))
				if err != nil {
					return 0, err
				}
				return pid, nil
			}
		}
	case "linux":
		cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port))
		output, err := cmd.Output()
		if err != nil {
			return 0, err
		}
		lines := strings.Split(string(output), "\n")
		if len(lines) < 2 {
			return 0, fmt.Errorf("process not found")
		}
		fields := strings.Fields(lines[1])
		if len(fields) < 2 {
			return 0, fmt.Errorf("process not found")
		}
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, err
		}
		return pid, nil
	}
	return 0, fmt.Errorf("unsupported operating system")
}

// 根据pid杀掉进程
func killProcess(pid int) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
	case "linux":
		cmd = exec.Command("kill", strconv.Itoa(pid))
	default:
		return os.ErrInvalid
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
