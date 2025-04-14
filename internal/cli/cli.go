package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/lyj404/iggen/pkg/utils"

	"github.com/lyj404/iggen/internal/generator"
	"github.com/lyj404/iggen/internal/github"
)

const (
	helpText = `iggen - .gitignore生成工具

用法:
  iggen <命令> [proxy <地址>] [参数...]

基础命令:
  list        查看所有模板
  search <关键词> 搜索模板
  gen <模板名> 生成文件
  help        显示帮助信息
`
)

var commandHelps = map[string]string{
	"gen":    "gen：生成指定语言的.gitignore文件\n格式：iggen gen <模板名称> [proxy 地址] [参数...]\n示例：iggen gen Python proxy 127.0.0.1:7890",
	"list":   "list：列出所有可用模板\n格式：iggen list [proxy 地址]\n示例：iggen list proxy 127.0.0.1:1080",
	"search": "search：搜索关键词相关模板\n格式：iggen search <关键词> [proxy 地址]\n示例：iggen search java",
	"proxy":  "proxy：设置HTTP代理服务器\n格式：iggen [主命令] proxy <IP:PORT>\n注意：必须与list/search/gen配合使用n示例：iggen list proxy 127.0.0.1:7890",
}

func Run(gh *github.GitHubClient, gen *generator.GitignoreGenerator) {
	// 提取iggen后面的命令参数()
	args := os.Args[1:]
	// 判断参数是否为0
	if len(args) == 0 {
		exitWithError("错误：缺少命令参数\n\n请使用 -help 查看帮助信息", nil)
	}
	if args[0] == "help" || args[0] == "h" {
		if helpText, ok := commandHelps[args[1]]; ok {
			fmt.Print(helpText)
			return
		}
		exitWithError(fmt.Sprintf("未知命令：%s", args[1]), nil)
		fmt.Print(helpText)
		return
	}
	var (
		proxyAddr string   // 代理服务器地址（格式：host:port）
		command   string   // 主命令（list/search/gen）
		cmdArgs   []string // 命令专属参数
	)
	// 参数解析流程
	command = args[0]         // 提取主命令
	remainingArgs := args[1:] // 剩余待处理参数
	// 代理配置解析（支持格式：iggen [command] proxy 127.0.0.1:7890 [...]）
	if len(remainingArgs) >= 1 && remainingArgs[0] == "proxy" {
		if len(remainingArgs) < 2 {
			exitWithError("proxy命令缺少代理地址", nil)
		}
		proxyAddr = remainingArgs[1] // 提取代理地址
		cmdArgs = remainingArgs[2:]  // 剩余参数给具体命令
	} else {
		cmdArgs = remainingArgs // 无代理参数直接传递
	}
	// 客户端实例化（根据代理配置选择基础客户端或代理客户端）
	client := gh // 默认使用基础客户端
	if proxyAddr != "" {
		// 当存在代理配置时，创建带代理的客户端
		client = github.NewGitHubClientWithProxy(proxyAddr)
	}
	// 创建请求上下文（可用于超时控制/取消请求）
	ctx := context.Background()
	// 命令路由分发
	switch command {
	case "list": // 列出所有可用模板
		handleList(ctx, client)
	case "search": // 搜索模板
		handleSearch(ctx, client, cmdArgs)
	case "gen": // 生成.gitignore文件
		handleGenerate(ctx, client, gen, cmdArgs)
	default: // 未知命令处理
		fmt.Printf("未知命令: %s\n\n%s", command, helpText)
	}
}

// handleList用于执行list命令来获取所有.gitignore模板
func handleList(ctx context.Context, gh *github.GitHubClient) {
	// 获取所有.gitignore模板
	templates, err := gh.ListTemplates(ctx)
	if err != nil {
		exitWithError("获取模板列表失败", err)
	}

	fmt.Println("可用模板列表:")
	// 循环打印所有获取到的模板
	for _, t := range templates {
		fmt.Printf("  - %s\n", t)
	}
}

// handleSearch用于执行search命令来搜索.gitignore模板
func handleSearch(ctx context.Context, gh *github.GitHubClient, terms []string) {
	// 判断搜索关键字是否为空
	if len(terms) == 0 {
		exitWithError("搜索需要关键词", nil)
	}

	// 获取所有.gitignore模板
	templates, err := gh.ListTemplates(ctx)
	if err != nil {
		exitWithError("获取模板失败", err)
	}

	// 从模板列表中搜索关键字对应的模板
	results := utils.FuzzySearch(templates, terms[0])
	// 判断是否搜索到对应的模板
	if len(results) == 0 {
		fmt.Println("没有找到匹配模板")
		return
	}

	fmt.Printf("找到%d个匹配结果:\n", len(results))
	// 循环打印所有解决
	for _, r := range results {
		fmt.Printf("  \033[32m%s\033[0m\n", r)
	}
}

// handleGenerate用于执行gen命令来生成.gitignore模板
func handleGenerate(ctx context.Context, gh *github.GitHubClient, gen *generator.GitignoreGenerator, names []string) {
	// 判断模板名是否为空
	if len(names) == 0 {
		exitWithError("请指定至少一个模板", nil)
	}

	// 判断.gitignore文件是否存在，如果存在是否选项覆盖
	if gen.FileExists() && !confirmOverwrite() {
		fmt.Println("操作已取消")
		return
	}

	// 声明变量存储模板内容
	contents := make([][]byte, 0)
	// 遍历模板名
	for _, name := range names {
		// 根据模板名获取模板内容
		content, err := gh.GetTemplate(ctx, name)
		// 判断模板内容是否获取成
		if err != nil {
			exitWithError(fmt.Sprintf("获取模板 %s 失败", name), err)
		}
		contents = append(contents, content)
	}

	// 将模板内容生成到.gitignore文件中
	if err := gen.Generate(contents, names...); err != nil {
		exitWithError("生成文件失败", err)
	}

	fmt.Printf("成功生成.gitignore（包含: %s）\n", strings.Join(names, ", "))
}

// confirmOverwrite用于确认覆盖已存在.gitignore文件
func confirmOverwrite() bool {
	fmt.Print(".gitignore已存在，是否覆盖？(y/N) ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.ToLower(scanner.Text()) == "y"
}

// exitWithError提示错误信息
func exitWithError(msg string, err error) {
	fmt.Printf("\033[31m错误: %s\033[0m\n", msg)
	if err != nil {
		fmt.Printf("详情: %v\n", err)
	}
	os.Exit(1)
}
