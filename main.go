package main

import (
	"bufio"
	"fmt"
	"jar_tools/config"
	"jar_tools/consts"
	"jar_tools/utils/inputUtil"
	"jar_tools/utils/osUtil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	// 作者信息
	fmt.Printf("%s: %s By 2023\n\n", consts.AppName, consts.AppAuthor)
	// 选择操作
	fmt.Printf("请选择操作：\n1. 启动JAR程序\n2. 停止JAR程序\n3. 初始化配置\n")
	input := inputUtil.GetInputWithPrompt("请输入操作序号:")
	if input == "1" {
		f := config.GetConfig()
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
		f := config.GetConfig()
		pid := getJarPidByPort(f.Port)
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
		err := config.InitConfig()
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
	if pid := getJarPidByPort(port); pid != 0 {
		fmt.Printf("Error: 端口 %d 已被占用，进程ID为 %d\n", port, pid)
		input := inputUtil.GetInputWithPromptAndDefault(fmt.Sprintf("是否杀掉进程 %d (y/n)？", pid), "n")
		if input == "y" || input == "Y" {
			if err := killProcess(pid); err != nil {
				return err
			}
			fmt.Printf("进程 %d 已杀掉\n", pid)
		} else {
			return fmt.Errorf("端口 %d 已被占用，操作取消", port)
		}
	}
	f := config.GetConfig()
	jvm := append(strings.Split(f.Jvm, " "), "--server.port="+strconv.Itoa(port))
	ary := append([]string{"-jar"}, jarPath)
	ary = append(ary, jvm...)
	var cmd *exec.Cmd
	switch osUtil.GetOsType() {
	case consts.OsTypeWindows:
		// windows下使用javaw命令，不弹出黑框
		cmd = exec.Command("javaw", ary...)
	case consts.OsTypeLinux:
		// 用nohup命令，不挂断运行
		ary = append(ary, "&")
		// ary前面加上nohup
		ary = append([]string{"java"}, ary...)
		cmd = exec.Command("nohup", ary...)
	default:
		return fmt.Errorf("不支持的操作系统: %s", osUtil.GetOsType())
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	fmt.Printf("JAR程序已启动，端口为 %d\n", f.Port)
	// 写入 PID 文件
	cof := config.NewConfig()
	cof.Port = port
	cof.JarPath = jarPath
	cof.Jvm = f.Jvm
	err = config.SetConfig(cof)
	if err != nil {
		return fmt.Errorf("写入 PID 文件失败: %s", err)
	}
	return nil
}

func getJarPidByPort(port int) int {
	var pid int

	// Determine the OS type
	osType := osUtil.GetOsType()

	// Construct the command to get the PID of the Java Jar process listening on the specified port
	var cmd *exec.Cmd
	switch osType {
	case consts.OsTypeWindows:
		cmd = exec.Command("cmd", "/C", "netstat -ano | findstr :"+strconv.Itoa(port))
	case consts.OsTypeLinux:
		cmd = exec.Command("sh", "-c", "lsof -i :"+strconv.Itoa(port)+" | awk '{print $2}'")
	}

	// Execute the command and get its output
	out, err := cmd.Output()
	if err != nil {
		return pid
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

	return pid
}

// 根据pid杀掉进程
func killProcess(pid int) error {
	var cmd *exec.Cmd
	switch osUtil.GetOsType() {
	case consts.OsTypeWindows:
		cmd = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
	case consts.OsTypeLinux:
		cmd = exec.Command("kill", strconv.Itoa(pid))
	default:
		return os.ErrInvalid
	}
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
