package cli

var commandHelps = map[string]string{
	"helpText": "iggen - .gitignore生成工具\n" +
		"用法:\n" +
		"iggen <命令> [-proxy <地址>] [参数...]\n" +
		"基础命令:\n" +
		"list        查看所有模板\n" +
		"search <关键词> 搜索模板（支持正则匹配）\n" +
		"gen <模板名> 生成文件\n" +
		"help        显示帮助信息",
	"gen": "生成指定语言的.gitignore文件\n" +
		"格式: iggen gen <模板名称> [-proxy 地址] [参数...]\n" +
		"示例: iggen gen Python -proxy 127.0.0.1:7890",

	"list": "列出所有可用模板\n" +
		"格式: iggen list [-proxy 地址]\n" +
		"示例: iggen list -proxy 127.0.0.1:1080",

	"search": "搜索关键词相关模板\n" +
		"格式: iggen search <关键词> [-proxy 地址]\n" +
		"示例: iggen search java / iggen search \"^z\"",

	"proxy": "设置HTTP代理服务器\n" +
		"格式: iggen [主命令] -proxy <IP:PORT>\n" +
		"注意: 必须与list/search/gen配合使用\n" +
		"示例: iggen list -proxy 127.0.0.1:7890",
}
