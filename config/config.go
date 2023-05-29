/**
  @author: XingGao
  @since: 2023/4/5
  @desc:
**/

package config

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
	Jvm     string
}

func NewConfig() *Config {
	return &Config{}
}

// GetConfig 读取配置文件
func GetConfig() *Config {
	c := &Config{}
	f, err := os.Open(filePath + "config.txt")
	if err != nil {
		fmt.Println("Error: 找不到配置文件，请检查")
		return nil
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
			port, _ := strconv.Atoi(strings.Split(line, "=")[1])
			c.Port = port
		} else if strings.HasPrefix(line, "jvm") {
			// 去掉第一个=号前字符
			c.Jvm = strings.SplitN(line, "=", 2)[1]
		}

	}
	return c
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
