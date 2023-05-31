Java Jar部署脚本
这是一个使用GO语言编写的Java Jar部署脚本，它可以在Windows和Linux系统上运行。该脚本具有以下功能：

启动/停止Jar程序
设置Jar程序开机自启
如何使用
首先，你需要安装GO语言环境。你可以从官方网站 https://golang.org/ 下载并安装GO语言环境。

下载该脚本并将其保存为deploy.jar.go。

在脚本所在的目录中创建一个名为config.json的配置文件，配置文件内容如下：

{
  "jarPath": "/path/to/your/jar/file.jar",
  "startCommand": "java -jar /path/to/your/jar/file.jar",
  "stopCommand": "kill $(ps aux | grep '[j]ava -jar /path/to/your/jar/file.jar' | awk '{print $2}')",
  "autostart": true
}
请根据实际情况修改jarPath、startCommand和stopCommand字段的值。

运行以下命令编译脚本：
go build deploy.jar.go
运行脚本：
./deploy.jar
现在，你可以使用以下命令来启动/停止Jar程序：

./deploy.jar start
./deploy.jar stop
如果你想设置Jar程序开机自启，可以运行以下命令：

./deploy.jar autostart
如果你想取消开机自启，可以运行以下命令：

./deploy.jar noautostart
兼容性
该脚本已经在以下系统上测试通过：

Windows 10
Ubuntu 18.04
CentOS 7
注意事项
请确保config.json文件中的路径正确，否则程序将无法正常启动。
在Linux系统上，需要使用root权限运行该脚本以设置开机自启。