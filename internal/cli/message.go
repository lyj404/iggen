package cli

const (
	helpText = `iggen - .gitignore生成工具

用法:
  iggen <命令> [proxy <地址>] [参数...]

基础命令:
  list        查看所有模板
  search <关键词> 搜索模板（支持正则匹配）
  gen <模板名> 生成文件
  help        显示帮助信息
`
)

var commandHelps = map[string]string{
	"gen": "gen：生成指定语言的.gitignore文件\n" +
		"格式：iggen gen <模板名称> [proxy 地址] [参数...]\n" +
		"示例：iggen gen Python proxy 127.0.0.1:7890",

	"list": "list：列出所有可用模板\n" +
		"格式：iggen list [proxy 地址]\n" +
		"示例：iggen list proxy 127.0.0.1:1080",

	"search": "search：搜索关键词相关模板\n" +
		"格式：iggen search <关键词> [proxy 地址]\n" +
		`示例：iggen search java / iggen search "^z"`,

	"proxy": "proxy：设置HTTP代理服务器\n" +
		"格式：iggen [主命令] proxy <IP:PORT>\n" +
		"注意：必须与list/search/gen配合使用\n" +
		"示例：iggen list proxy 127.0.0.1:7890",
}
