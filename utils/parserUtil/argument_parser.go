package parserUtil

import (
	"flag"
	"fmt"
	"jar_tools/app"
	"jar_tools/consts"
)

/**
  @author: XingGao
  @date: 2023/5/29
**/

// ParseArgs 解析命令行参数
func ParseArgs() bool {
	// 版本命令行参数
	version := flag.Bool("v", false, "打印版本信息")
	// mode 命令行参数
	mode := flag.String("mode", "", "运行模式")

	// 解析命令行参数
	flag.Parse()

	// 如果 version 标志被设置为 true，则打印版本信息并退出
	if *version {
		fmt.Println(consts.AppVersion)
		return false
	}
	if *mode != "" {
		// 指定运行模式 直接运行
		app.Start(*mode)
		return false
	}
	return true
}
