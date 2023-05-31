# GO语言编写的Java Jar部署脚本

这是一个用GO语言编写的Java Jar部署脚本，它提供了一些功能，包括开关Jar程序和设置开机自启。该脚本适用于Windows和Linux系统。

## 功能特点

- 开启和关闭Jar程序
- 设置开机自启动

## 使用方法

### 配置

在使用之前，你需要进行一些配置。

1. 确保你已经安装了Java运行环境，并将其添加到系统环境变量中。
2. 确保你已经安装了GO语言编译器，并将其添加到系统环境变量中。

### 下载

你可以通过以下方式获取脚本代码：

```shell
$ git clone https://github.com/your-username/your-repository.git
```

### 编译

在下载完代码后，进入项目目录，并执行以下命令进行编译：

```shell
$ go build deploy.go
```

### 配置文件

在项目目录中，你需要创建一个名为`config.json`的配置文件，用于配置Jar程序的相关信息。配置文件的示例内容如下：

```json
{
  "jarPath": "/path/to/your/jar/file.jar",
  "startOnBoot": true,
  "startCommand": "java -jar /path/to/your/jar/file.jar",
  "stopCommand": "pkill -f /path/to/your/jar/file.jar"
}
```

你需要将`jarPath`字段的值替换为你的Jar文件的路径。`startCommand`和`stopCommand`字段分别是启动和停止Jar程序的命令。根据你的实际情况，可能需要修改这些命令。

### 运行

在配置完成后，你可以运行编译后的脚本来启动、停止和设置开机自启动Jar程序。

- 启动Jar程序：

  ```shell
  $ ./deploy start
  ```

- 停止Jar程序：

  ```shell
  $ ./deploy stop
  ```

- 设置开机自启动：

  ```shell
  $ ./deploy enable-autostart
  ```

- 取消开机自启动：

  ```shell
  $ ./deploy disable-autostart
  ```

## 注意事项

- 请确保在运行脚本之前，已经正确配置了`config.json`文件。
- 在Linux系统中，你可能需要使用`chmod`命令为脚本文件赋予执行权限。
- 如果你需要在Windows系统中使用该脚本，请使用合适的GO语言编译器进行编译，并按照Windows环境的要求进行配置和运行。

## 贡献

如果你有任何建议或改进，请随时提出。你可以通过GitHub的Pull Request功能向该项目贡献代码。

## 许可证

该脚本采用[MIT许可证](https://opensource.org/licenses/MIT)进行许可。详情请参阅`LICENSE`文件。