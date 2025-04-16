# IgGen

IgGen是一个用于生成.gitignore文件的小工具

# 📜 使用说明

**命令格式**： `iggen <主命令> [proxy <IP:PORT>] [参数...]`
## 使用方法
```
主命令：
    list        列出所有可用的 .gitignore 模板
    search      搜索特定的 .gitignore 模板（支持正则匹配）
    gen         生成 .gitignore 文件
    help, h     显示帮助信息

子命令：
    -proxy      代理地址（格式：IP:PORT），该子命令不能单独使用

示例：
    # 列出所有模板
    iggen list
    iggen list -proxy 127.0.0.1:7890

    # 搜索模板
    iggen search go  
    iggen search "^z"                 
    iggen search java -proxy 127.0.0.1:7890

    # 生成gitignore文件
    iggen gen java
    iggen gen node -proxy 127.0.0.1:7890

    # 查看帮助
    iggen help
    iggen h
    iggen help list    # 查看list命令的详细帮助
    iggen help search  # 查看search命令的详细帮助
    iggen help gen     # 查看gen命令的详细帮助

注意：
    1. list、search、gen命令都支持代理设置
    2. search和gen命令需要提供额外的参数
    3. search命令的正则匹配不区分大小写，且需要使用双引号包裹
    4. 如果不提供任何参数，将显示错误信息和帮助信息
```

⚠️ **注意事项**

- 启用代理时请确保网络可达性
- 只支持支持HTTP代理协议（暂不支持SOCKS）

# 📖 数据来源

本工具所有.gitignore模板均来自 [GitHub官方gitignore仓库](https://github.com/github/gitignore)，通过GitHub API实时获取最新版本。
