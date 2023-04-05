package main

import (
	"bufio"
	"fmt"
	"jar_tools/consts"
	"jar_tools/file"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
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
		pid, err := getJarPidByPort(f.Port)
		if pid != 0 {
			err := killProcess(pid)
			if err != nil {
				return
			} else {
				fmt.Printf("端口 %d ,进程 %d 已杀掉\n", f.Port, pid)
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
	//if pid, err := getJarPidByPort(port); err == nil && pid != 0 {
	//	fmt.Printf("Error: 端口 %d 已被占用，进程ID为 %d\n", port, pid)
	//	fmt.Printf("是否杀掉进程 %d (y/n)？\n", pid)
	//	var input string
	//	_, err = fmt.Scanf("\n%s", &input)
	//	if input == "y" || input == "Y" {
	//		if err := killProcess(pid); err != nil {
	//			return err
	//		}
	//		fmt.Printf("进程 %d 已杀掉\n", pid)
	//	} else {
	//		return fmt.Errorf("端口 %d 已被占用，操作取消", port)
	//	}
	//}
	f, err := file.GetConfig()
	jvm := append(strings.Split(f.Jvm, " "), "--server.port="+strconv.Itoa(port))
	ary := append([]string{"-jar"}, jarPath)
	ary = append(ary, jvm...)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// windows下使用javaw命令，不弹出黑框
		cmd = exec.Command("javaw", ary...)
	} else if runtime.GOOS == "linux" {
		// 用nohup命令，不挂断运行
		ary = append(ary, "&")
		// ary前面加上nohup
		ary = append([]string{"java"}, ary...)
		cmd = exec.Command("nohup", ary...)
	} else {
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	fmt.Printf("JAR程序已启动，端口为 %d\n", f.Port)
	// 写入 PID 文件
	config := file.NewConfig()
	config.Port = port
	config.JarPath = jarPath
	config.Jvm = f.Jvm
	err = file.SetConfig(config)
	if err != nil {
		return fmt.Errorf("写入 PID 文件失败: %s", err)
	}
	return nil
}

func getJarPidByPort(port int) (int, error) {
	var pid int

	// Determine the OS type
	osType := runtime.GOOS

	// Construct the command to get the PID of the Java Jar process listening on the specified port
	var cmd *exec.Cmd
	if osType == "windows" {
		cmd = exec.Command("cmd", "/C", "netstat -ano | findstr :"+strconv.Itoa(port))
	} else {
		cmd = exec.Command("sh", "-c", "lsof -i :"+strconv.Itoa(port)+" | awk '{print $2}'")
	}

	// Execute the command and get its output
	out, err := cmd.Output()
	if err != nil {
		return pid, err
	}

	// Parse the output to get the PID of the Java Jar process
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if osType == "windows" {
			fields := strings.Fields(line)
			if len(fields) > 4 && strings.Contains(fields[1], ":") {
				p := strings.Split(fields[1], ":")
				if p[len(p)-1] == strconv.Itoa(port) {
					pid, _ = strconv.Atoi(fields[4])
					break
				}
			}
		} else {
			if p, err := strconv.Atoi(line); err == nil {
				pid = p
				break
			}
		}
	}

	return pid, nil
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
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
