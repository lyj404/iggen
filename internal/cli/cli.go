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
  iggen list                   查看所有模板
  iggen search <关键词>         搜索模板
  iggen gen <模板名> 			 生成文件

示例:
  iggen gen Go Python         生成Go和Python的合并配置
  iggen search Rust           查找Rust相关模板
`

func Run(gh *github.GitHubClient, gen *generator.GitignoreGenerator) {
	if len(os.Args) < 2 {
		fmt.Print(helpText)
		return
	}

	ctx := context.Background()
	command := os.Args[1]

	switch command {
	case "list":
		handleList(ctx, gh)
	case "search":
		handleSearch(ctx, gh, os.Args[2:])
	case "gen":
		handleGenerate(ctx, gh, gen, os.Args[2:])
	default:
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
