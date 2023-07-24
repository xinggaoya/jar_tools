package jarUtil

import (
	"bufio"
	"fmt"
	"jar_tools/config"
	"jar_tools/consts"
	"jar_tools/utils/fileUtil"
	"jar_tools/utils/inputUtil"
	"jar_tools/utils/osUtil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/**
  @author: XingGao
  @date: 2023/5/29
**/

func RunJar(jarPath string, port int) error {
	wsPath := fileUtil.GetCurrentDirectory()
	jarName := jarPath
	jarPath = wsPath + "/" + jarPath
	if pid := GetJarPidByPort(port); pid != 0 {
		// 端口被占用，询问是否杀掉进程
		msg := fmt.Sprintf("Error: 端口 %d 已被占用，进程ID为 %d\n是否杀掉进程 %d (y/n)?", port, pid, pid)
		input := inputUtil.GetInputWithPromptAndDefault(msg, "n")
		if input == "y" || input == "Y" {
			if err := KillProcess(pid); err != nil {
				return err
			}
			fmt.Printf("进程 %d 已杀掉\n", pid)
		} else {
			return fmt.Errorf("\n端口 %d 已被占用，操作取消", port)
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
		break
	case consts.OsTypeLinux:
		// 用nohup命令，不挂断运行
		ary = append(ary, "&")
		// ary前面加上nohup
		ary = append([]string{"java"}, ary...)
		cmd = exec.Command("nohup", ary...)
		break
	default:
		return fmt.Errorf("不支持的操作系统: %s", osUtil.GetOsType())
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	fmt.Printf("执行命令: %s\n", ary)
	fmt.Printf("JAR程序已启动，端口为 %d\n", f.Port)
	// 写入 PID 文件
	cof := config.NewConfig()
	cof.Port = port
	cof.JarPath = jarName
	cof.Jvm = f.Jvm
	config.SetConfig(cof)
	return nil
}

// GetJarPidByPort 获取指定端口的进程ID
func GetJarPidByPort(port int) int {
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

// KillProcess 根据pid杀掉进程
func KillProcess(pid int) error {
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
