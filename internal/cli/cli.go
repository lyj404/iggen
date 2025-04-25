package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lyj404/iggen/pkg/utils"

	"github.com/lyj404/iggen/internal/generator"
	"github.com/lyj404/iggen/internal/github"
)

func Run(gh *github.GitHubClient, gen *generator.GitignoreGenerator) {
	// 创建主命令集
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	genCmd := flag.NewFlagSet("gen", flag.ExitOnError)
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	// 为list、search、gen、help添加proxy子命令
	proxyHelp := "代理地址（格式： IP:PORT）"
	listProxy := listCmd.String("proxy", "", proxyHelp)
	searchProxy := searchCmd.String("proxy", "", proxyHelp)
	genProxy := genCmd.String("proxy", "", proxyHelp)

	if len(os.Args) < 2 {
		utils.ExitWithError("错误：缺少命令参数\n\n请使用 help命令 查看帮助信息", nil)
	}

	// 创建请求上下文（可用于超时控制/取消请求）
	ctx := context.Background()

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		client := getClient(gh, *listProxy)
		handleList(ctx, client)

	case "search":
		searchCmd.Parse(os.Args[2:])
		client := getClient(gh, *searchProxy)
		handleSearch(ctx, client, searchCmd.Args())

	case "gen":
		genCmd.Parse(os.Args[2:])
		client := getClient(gh, *genProxy)
		handleGenerate(ctx, client, gen, genCmd.Args())

	case "help", "h":
		helpCmd.Parse(os.Args[2:])
		if len(helpCmd.Args()) > 0 {
			if helpText, ok := commandHelps[helpCmd.Args()[0]]; ok {
				fmt.Print(helpText)
				return
			}
			utils.ExitWithError(fmt.Sprintf("未知命令：%s", helpCmd.Args()[0]), nil)
		}
		fmt.Print(helpText)

	default:
		fmt.Printf("未知命令: %s\n\n%s", os.Args[1], helpText)
	}

}

// getClient根据是否设置代理返回不同的GitHub客户端
func getClient(gh *github.GitHubClient, proxy string) *github.GitHubClient {
	if proxy != "" {
		return github.NewGitHubClientWithProxy(proxy)
	}
	return gh
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
