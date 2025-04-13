# IgGen

IgGen是一个用于生成.gitignore文件的小工具

# 📜 使用说明

**命令格式**： `iggen <命令> [proxy <IP:PORT>] [参数...]`

## 命令列表

| 命令   | 参数                       | 作用                         | 示例                |
| ------ | -------------------------- | ---------------------------- | ------------------- |
| list   | 无                         | 列出所有可用的.gitignore模板 | iggen list          |
| search | 关键词                     | 搜索模板                     | iggen search 关键词 |
| gen    | 模板名称（可不区分大小写） | 生成对应的.gitignore模板     | iggen gen 模板名称  |

> 在不能直连GitHub时，可以使用`proxy`命令来确保上诉命令的正确执行，proxy命令示例：`iggen list proxy 127.0.0.1:1080`，`list、search、gen`这些命令均可以使用`proxy`

⚠️ **注意事项**

- 代理参数必须直接跟在命令词之后
- 地址格式需严格遵循 IP:PORT 模式
- 支持HTTP代理协议（暂不支持SOCKS）
- 启用代理时请确保网络可达性

# 📖 数据来源

本工具所有.gitignore模板均来自 [GitHub官方gitignore仓库](https://github.com/github/gitignore)，通过GitHub API实时获取最新版本。
