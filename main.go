package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	var jarPath string
	var port int

	// 选择操作
	fmt.Printf("请选择操作：\n1. 启动JAR程序\n2. 停止JAR程序\n ")
	fmt.Print("请输入操作序号: ")
	var input string
	_, err := fmt.Scanf("%s", &input)
	if err != nil {
		return
	}
	if input == "1" {
		fmt.Println("请输入JAR文件路径和端口号，例如：./app.jar 8080")
		_, err = fmt.Scanf("%s %d", &jarPath, &port)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		if _, err := os.Stat(jarPath); os.IsNotExist(err) {
			fmt.Printf("Error: 找不到指定文件，请检查: %s\n", jarPath)
			return
		}

		// 执行 JAR 程序
		if err := RunJar(jarPath, port); err != nil {
			fmt.Println(err)
			return
		}
	} else if input == "2" {
		// 获取 JAR 文件路径和端口号
		f, err := os.Open("./pid.txt")
		if err != nil {
			fmt.Println("Error: 找不到PID文件，请检查")
			return
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				return
			}
		}(f)
		// 读取 PID 文件
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fmt.Printf("检测到 PID 文件，进程ID为 %s\n", scanner.Text())
			pid, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := killProcess(pid); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("进程 %d 已杀掉\n", pid)
		}
	} else {
		fmt.Println("Error: 无效的输入")
		return
	}
}

func RunJar(jarPath string, port int) error {
	if pid, err := FindProcessByPort(port); err == nil {
		fmt.Printf("Error: 端口 %d 已被占用，进程ID为 %d\n", port, pid)
		fmt.Printf("是否杀掉进程 %d (y/n)？", pid)
		var input string
		_, err = fmt.Scanln("%d", &input)
		if input == "y" || input == "Y" {
			if err := killProcess(pid); err != nil {
				return err
			}
			fmt.Printf("进程 %d 已杀掉\n", pid)
		} else {
			return fmt.Errorf("端口 %d 已被占用，操作取消", port)
		}
	}
	cmd := exec.Command("javaw", "-jar", jarPath)
	err := cmd.Start()
	if err != nil {
		return err
	}
	pid, err := cmd.Process.Pid, cmd.Process.Release()
	if err != nil {
		return err
	}
	fmt.Printf("JAR程序已启动，进程ID为 %d\n", pid)
	// 写入 PID 文件
	f, err := os.Create("./pid.txt")
	_, err = f.Write([]byte(strconv.Itoa(pid)))
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	fmt.Println("PID文件已写入")
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
