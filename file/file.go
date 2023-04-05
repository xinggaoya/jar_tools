/**
  @author: XingGao
  @since: 2023/4/5
  @desc:
**/

package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const filePath = "./data/"

type Config struct {
	JarPath string
	Port    int
	Pid     int
	Jvm     string
}

func NewConfig() *Config {
	return &Config{}
}

// GetConfig 读取配置文件
func GetConfig() (*Config, error) {
	c := &Config{}
	f, err := os.Open(filePath + "config.txt")
	if err != nil {
		return nil, errors.New("error: 找不到配置文件，请检查")
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	// 读取配置文件
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "jarPath") {
			c.JarPath = strings.Split(line, "=")[1]
		} else if strings.HasPrefix(line, "port") {
			port, err := strconv.Atoi(strings.Split(line, "=")[1])
			if err != nil {
				return nil, errors.New("error: 配置文件中的端口号格式错误")
			}
			c.Port = port
		} else if strings.HasPrefix(line, "pid") {
			pid, err := strconv.Atoi(strings.Split(line, "=")[1])
			if err != nil {
				return nil, errors.New("error: 配置文件中的PID格式错误")
			}
			c.Pid = pid
		} else if strings.HasPrefix(line, "jvm") {
			c.Jvm = strings.Split(line, "=")[1]
		}

	}
	return c, nil
}

// SetConfig 写入配置文件
func SetConfig(c *Config) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.Mkdir(filePath, 0755)
	}
	f, err := os.Create(filePath + "config.txt")
	if err != nil {
		return errors.New("error: 创建配置文件失败")
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	if c.JarPath != "" {
		_, err = f.WriteString("jarPath=" + c.JarPath + "\n")
	}
	if c.Port != 0 {
		_, err = f.WriteString("port=" + strconv.Itoa(c.Port) + "\n")
	}
	if c.Pid != 0 {
		_, err = f.WriteString("pid=" + strconv.Itoa(c.Pid) + "\n")
	}
	if c.Jvm != "" {
		_, err = f.WriteString("jvm=" + c.Jvm + "\n")
	}
	if err != nil {
		return errors.New("error: 写入配置文件失败")
	}
	return nil
}

// InitConfig 初始化配置文件
func InitConfig() error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.Mkdir(filePath, 0755)
	}
	f, err := os.Create(filePath + "config.txt")
	if err != nil {
		return fmt.Errorf("error: 创建配置文件失败")
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	_, err = f.WriteString("jarPath=" + "\n")
	_, err = f.WriteString("port=" + "\n")
	_, err = f.WriteString("pid=" + "\n")
	_, err = f.WriteString("jvm=" + "\n")
	if err != nil {
		return fmt.Errorf("error: 写入配置文件失败")
	}
	return nil
}
