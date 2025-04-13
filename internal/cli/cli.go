package cli

import (
	"bufio"
	"context"
	"fmt"
	"iggen/pkg/utils"
	"os"
	"strings"

	"iggen/internal/generator"
	"iggen/internal/github"
)

const helpText = `iggen - .gitignore生成器

用法:
  iggen <命令> [proxy <地址>] [参数...]

基础命令:
  list        查看所有模板
  search <关键词> 搜索模板
  gen <模板名> 生成文件
`

func Run(gh *github.GitHubClient, gen *generator.GitignoreGenerator) {
	// 参数数量校验（至少需要1个命令参数）
	if len(os.Args) < 2 {
		fmt.Print(helpText)
		return
	}
	var (
		proxyAddr string   // 代理服务器地址（格式：host:port）
		command   string   // 主命令（list/search/gen）
		cmdArgs   []string // 命令专属参数
	)
	// 参数解析流程
	args := os.Args[1:]       // 剥离程序自身名称
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
	templates, err := gh.ListTemplates(ctx)
	if err != nil {
		exitWithError("获取模板列表失败", err)
	}

	fmt.Println("可用模板列表:")
	for _, t := range templates {
		fmt.Printf("  - %s\n", t)
	}
}

// handleSearch用于执行search命令来搜索.gitignore模板
func handleSearch(ctx context.Context, gh *github.GitHubClient, terms []string) {
	if len(terms) == 0 {
		exitWithError("搜索需要关键词", nil)
	}

	templates, err := gh.ListTemplates(ctx)
	if err != nil {
		exitWithError("获取模板失败", err)
	}

	results := utils.FuzzySearch(templates, terms[0])
	if len(results) == 0 {
		fmt.Println("没有找到匹配模板")
		return
	}

	fmt.Printf("找到%d个匹配结果:\n", len(results))
	for _, r := range results {
		fmt.Printf("  \033[32m%s\033[0m\n", r)
	}
}

// handleGenerate用于执行gen命令来生成.gitignore模板
func handleGenerate(ctx context.Context, gh *github.GitHubClient, gen *generator.GitignoreGenerator, names []string) {
	if len(names) == 0 {
		exitWithError("请指定至少一个模板", nil)
	}

	if gen.FileExists() && !confirmOverwrite() {
		fmt.Println("操作已取消")
		return
	}

	contents := make([][]byte, 0)
	for _, name := range names {
		content, err := gh.GetTemplate(ctx, name)
		if err != nil {
			exitWithError(fmt.Sprintf("获取模板 %s 失败", name), err)
		}
		contents = append(contents, content)
	}

	if err := gen.Generate(contents, names...); err != nil {
		exitWithError("生成文件失败", err)
	}

	fmt.Printf("成功生成.gitignore（包含: %s）\n", strings.Join(names, ", "))
}

// confirmOverwrite覆盖已存在.gitignore文件
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
