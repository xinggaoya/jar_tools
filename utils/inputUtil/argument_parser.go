package inputUtil

import "fmt"

// GetInput 获取用户输入
func GetInput() string {
	var input string
	_, err := fmt.Scanf("%s", &input)
	if err != nil {
		return ""
	}
	return input
}

// GetInputWithPrompt 获取用户输入，带提示
func GetInputWithPrompt(prompt string) string {
	fmt.Print(prompt)
	return GetInput()
}

// GetInputWithPromptAndDefault 获取用户输入，带提示和默认值
func GetInputWithPromptAndDefault(prompt string, defaultValue string) string {
	fmt.Printf("%s [%s]:", prompt, defaultValue)
	input := GetInput()
	if input == "" {
		return defaultValue
	}
	return input
}
