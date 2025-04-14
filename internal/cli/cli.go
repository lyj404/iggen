package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/lyj404/iggen/pkg/utils"

	"github.com/lyj404/iggen/internal/generator"
	"github.com/lyj404/iggen/internal/github"
)

func Run(gh *github.GitHubClient, gen *generator.GitignoreGenerator) {
	// 提取iggen后面的命令参数()
	args := os.Args[1:]
	// 判断参数是否为0
	if len(args) == 0 {
		utils.ExitWithError("错误：缺少命令参数\n\n请使用 -help 查看帮助信息", nil)
	}
	if args[0] == "help" || args[0] == "h" {
		if helpText, ok := commandHelps[args[1]]; ok {
			fmt.Print(helpText)
			return
		}
		utils.ExitWithError(fmt.Sprintf("未知命令：%s", args[1]), nil)
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

	// 寻找proxy命令在所在位置
	proxyIndex := -1
	for i, arg := range remainingArgs {
		if arg == "proxy" {
			proxyIndex = i
			break
		}
	}

	// 处理代理参数
	if proxyIndex != -1 {
		// 检查proxy后面是否有代理地址
		if proxyIndex+1 >= len(remainingArgs) {
			utils.ExitWithError("proxy命令缺少代理地址", nil)
		}

		// 提取代理地址
		proxyAddr = remainingArgs[proxyIndex+1]
		// 提取主命令后的参数
		cmdArgs = remainingArgs[:proxyIndex]
	} else {
		cmdArgs = remainingArgs
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
		utils.ExitWithError("获取模板列表失败", err)
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
		utils.ExitWithError("搜索需要关键词", nil)
	}

	// 获取所有.gitignore模板
	templates, err := gh.ListTemplates(ctx)
	if err != nil {
		utils.ExitWithError("获取模板失败", err)
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
		utils.ExitWithError("请指定至少一个模板", nil)
	}

	// 判断.gitignore文件是否存在，如果存在是否选项覆盖
	if gen.FileExists() && !utils.ConfirmOverwrite() {
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
			utils.ExitWithError(fmt.Sprintf("获取模板 %s 失败", name), err)
		}
		contents = append(contents, content)
	}

	// 将模板内容生成到.gitignore文件中
	if err := gen.Generate(contents, names...); err != nil {
		utils.ExitWithError("生成文件失败", err)
	}

	fmt.Printf("成功生成.gitignore（包含: %s）\n", strings.Join(names, ", "))
}
