package inputUtil

import (
	"bufio"
	"fmt"
	"os"
)

// GetInput 获取用户输入
func GetInput() string {

	var input string
	_, err := fmt.Scanf("%s", &input)
	if err != nil {
		_ = fmt.Errorf("error: %s", err.Error())
	}
	// 读取输入缓冲区中的换行符
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	return input
}

// GetInputWithPrompt 获取用户输入，带提示
func GetInputWithPrompt(prompt string) string {
	fmt.Print(prompt)
	return GetInput()
}

// GetInputWithPromptAndDefault 获取用户输入，带提示和默认值
func GetInputWithPromptAndDefault(prompt string, defaultValue string) string {
	fmt.Printf("%s 默认[%s]:", prompt, defaultValue)
	input := GetInput()
	if input == "" {
		return defaultValue
	}
	return input
}
