package cmd

import (
	"fmt"
	"jar_tools/app"
	"jar_tools/config"
	"jar_tools/consts"
	"jar_tools/utils/inputUtil"
	"jar_tools/utils/parserUtil"
)

/**
  @author: XingGao
  @date: 2023/5/29
**/

func Run() {
	// 作者信息
	fmt.Printf("%s: %s By 2023\n\n", consts.AppName, consts.AppAuthor)
	// 判断配置文件是否存在
	if !config.IsExist() {
		config.InitConfig()
		fmt.Printf("配置文件初始化成功,请先进行编辑\n")
		// 睡眠
		inputUtil.GetInputWithPrompt("按任意键继续...")
		return
	}
	// 解析命令行参数
	if parserUtil.ParseArgs() {
		// 选择操作
		fmt.Printf("请选择操作：\n1. 启动JAR程序\n2. 停止JAR程序\n3. 初始化配置\n")
		input := inputUtil.GetInputWithPrompt("请输入操作序号:")
		app.Start(input)
	}
}
