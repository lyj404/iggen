# IgGen

IgGen是一个用于生成.gitignore文件的小工具

# 📜 使用说明

**命令格式**： `iggen <主命令> [proxy <IP:PORT>] [参数...]`

## 命令列表

| 命令   | 参数                       | 作用                         | 示例                                    |
| ------ | -------------------------- | ---------------------------- | --------------------------------------- |
| list   | 无                         | 列出所有可用的.gitignore模板 | iggen list                              |
| search | 关键词（支持正则匹配）     | 搜索模板                     | iggen search 关键词 / iggen search "^z" |
| gen    | 模板名称（可不区分大小写） | 生成对应的.gitignore模板     | iggen gen 模板名称                      |
| proxy  | 代理地址                   | 设置HTTP代理                 | iggen list proxy 127.0.0.1:1080         |
| help/h | 无                         | 显示帮助信息                 | iggen help                              |

> 在不能直连GitHub时，可以使用`proxy`命令来确保上诉命令的正确执行，`proxy`命令只能搭配`list、serach、gen`这三个主命令一起使用

⚠️ **注意事项**

- 启用代理时请确保网络可达性
- 代理参数必须直接跟在命令词之后
- 地址格式需严格遵循 IP:PORT 模式
- 支持HTTP代理协议（暂不支持SOCKS）
- `search`命令的正则匹配不区分大小写，且需要使用双引号包裹

# 📖 数据来源

本工具所有.gitignore模板均来自 [GitHub官方gitignore仓库](https://github.com/github/gitignore)，通过GitHub API实时获取最新版本。
