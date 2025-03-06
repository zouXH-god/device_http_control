# 项目名称：device_http_control

## 项目介绍：

**device_http_control** 是一款基于HTTP协议的高效远程控制解决方案，专为现代用户设计，旨在通过网络实现对Windows系统的无缝管理和操作。该服务整合了媒体控制、系统管理以及程序快速启动功能，以轻量化的架构和强大的扩展性，为用户提供卓越的远程交互体验。无论是个人娱乐还是企业级应用，device_http_control 都能以其简洁的设计和可靠的性能，赋予用户前所未有的掌控力。

## 核心功能：
- **多媒体管理**：通过HTTP请求实现播放/暂停、曲目切换等操作，支持主流媒体播放器。
- **系统指令**：远程触发关机或重启，简化设备管理流程。
- **音量控制**：精准调节系统音量，优化音频体验。
- **动态程序启动**：基于配置文件，支持一键启动指定应用程序，提升工作效率。

## 应用场景：
- **智能家居**：在家庭网络中，通过移动设备远程操控多媒体播放。
- **企业管理**：IT管理员可批量远程重启或关闭办公设备，提高运维效率。
- **生产力优化**：快速启动常用工具或软件，助力用户专注于核心任务。

## 技术架构：
device_http_control 采用Go语言开发，利用`net/http`构建高性能服务器，结合`os/exec`实现系统命令调用，并通过`robotgo`库模拟键盘事件，完成多媒体与音量控制。服务内置HTTP基本认证机制，确保操作安全性。配置文件采用JSON格式，灵活定义快速启动程序列表，支持动态扩展。

## 部署与使用：
1. 下载 [二进制文件](https://github.com/zouXH-god/device_http_control/releases) 到目录中。
2. 在目录下创建并配置`config.json`，文件路径为`./config.json`，实例可参考`config.json`。
3. 启动服务，默认运行于`0.0.0.0:8080`。
4. 通过HTTP客户端（如Postman或`curl`）发送请求，调用所需功能。

## 配置文件 `config.json`：
```json
{
  "token": "TokenExample12345678",
  "server": {
    "port": "8080",
    "host": "0.0.0.0"
  },
  "programs": [
    {
      "name": "qqMusic",
      "path": "D:\\QQMusic\\QQMusic.exe"
    }
  ]
}
```

## 请求与响应：
- **请求**：使用HTTP GET请求。
- **响应**：简单的执行结果和状态码。

### 请求示例：
```curl
curl -H 'Authorization: TokenExample12345678' http://localhost:8080/play_pause
```

### 请求列表：
- `/play_pause`：播放/暂停
- `/previous`：上一曲
- `/next`：下一曲
- `/volume_up`：音量加
- `/volume_down`：音量减
- `/shutdown`：关机
- `/reboot`：重启
- `/launch?name=`：启动程序，参数为 `programs[i].name`

## 安全保障：
- **访问控制**：集成基本AUTH认证，限制未授权访问。
- **执行限制**：快速启动功能仅限于配置文件定义的程序，杜绝潜在风险。
