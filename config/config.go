/**
  @author: XingGao
  @since: 2023/4/5
  @desc:
**/

package config

import (
	"bufio"
	"fmt"
	"jar_tools/utils/fileUtil"
	"os"
	"strconv"
	"strings"
)

// 获取程序目录
var workPath = fileUtil.GetCurrentDirectory()

var filePath = workPath + "\\data\\"

type Config struct {
	JarPath    string
	Port       int
	Jvm        string
	ScriptName string
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
			c.JarPath = strings.SplitN(line, "=", 2)[1]
		}
		if strings.HasPrefix(line, "port") {
			c.Port, _ = strconv.Atoi(strings.SplitN(line, "=", 2)[1])
		}
		if strings.HasPrefix(line, "jvm") {
			c.Jvm = strings.SplitN(line, "=", 2)[1]
		}
		if strings.HasPrefix(line, "scriptName") {
			c.ScriptName = strings.SplitN(line, "=", 2)[1]
		}

	}
	return c
}

// SetConfig 写入配置文件
func SetConfig(c *Config) {
	config := GetConfig()
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.Mkdir(filePath, 0755)
	}
	f, err := os.Create(filePath + "config.txt")
	if err != nil {
		fmt.Errorf("error: 创建配置文件失败")
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	if c.JarPath != "" {
		config.JarPath = c.JarPath
	}
	if c.Port != 0 {
		config.Port = c.Port
	}
	if c.Jvm != "" {
		config.Jvm = c.Jvm
	}
	if c.ScriptName != "" {
		config.ScriptName = c.ScriptName
	}
	_, err = f.WriteString("jarPath=" + config.JarPath + "\n")
	_, err = f.WriteString("port=" + strconv.Itoa(config.Port) + "\n")
	_, err = f.WriteString("jvm=" + config.Jvm + "\n")
	_, err = f.WriteString("scriptName=" + config.ScriptName + "\n")
}

// InitConfig 初始化配置文件
func InitConfig() {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.Mkdir(filePath, 0755)
	}
	f, err := os.Create(filePath + "config.txt")
	if err != nil {
		fmt.Errorf("error: 创建配置文件失败")
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	_, err = f.WriteString("jarPath=" + "\n")
	_, err = f.WriteString("port=" + "\n")
	_, err = f.WriteString("jvm=" + "\n")
	if err != nil {
		fmt.Errorf("error: 写入配置文件失败")
	}
}

// IsExist 判断文件或文件夹是否存在
func IsExist() bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}
