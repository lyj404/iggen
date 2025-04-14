package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// FuzzySearch模糊搜索模板名称（支持部分匹配和相似字符匹配）
func FuzzySearch(templates []string, term string) []string {
	// 统一转为小写实现大小写不敏感搜索[^1]
	term = strings.ToLower(term)
	var reg *regexp.Regexp
	useRegex := isRegex(term)

	// 判断是否使用正则
	if useRegex {
		// 自动添加(?!)实现大小写不敏感
		reg = regexp.MustCompile("(?i)" + term)
	}

	matches := make([]string, 0)

	// 遍历所有模板进行双阶段匹配
	for _, t := range templates {
		// 将从模板列表中遍历出来的模板名称改成小写
		lower := strings.ToLower(t)

		// 正则匹配
		if useRegex && reg.MatchString(lower) {
			matches = append(matches, t)
		}

		// 模糊匹配
		if strings.Contains(lower, term) || compareSimilar(lower, term) {
			matches = append(matches, t)
		}
	}

	return matches
}

// compareSimilar比较处理后的相似字符串（支持前缀/后缀匹配）
func compareSimilar(a, b string) bool {
	// 移除特殊字符统一格式
	a = removeSpecialChars(a)
	b = removeSpecialChars(b)
	// 允许前缀匹配（如"nodejs"匹配"node"）
	// 允许后缀匹配（如"python3"匹配"3"）
	return strings.HasPrefix(a, b) || strings.HasSuffix(a, b)
}

// removeSpecialChars过滤非字母数字字符（增强模糊匹配能力）
func removeSpecialChars(s string) string {
	return strings.Map(func(r rune) rune {
		// 保留字母和数字（包括unicode字符）
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, s)
}

// ExitWithError提示错误信息
func ExitWithError(msg string, err error) {
	fmt.Printf("\033[31m错误: %s\033[0m\n", msg)
	if err != nil {
		fmt.Printf("详情: %v\n", err)
	}
	os.Exit(1)
}

// ConfirmOverwrite用于确认覆盖已存在.gitignore文件
func ConfirmOverwrite() bool {
	fmt.Print(".gitignore已存在，是否覆盖？(y/N) ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.ToLower(scanner.Text()) == "y"
}

// isRegex判断是否正则表达式
func isRegex(pattern string) bool {
	_, err := regexp.Compile(pattern)
	return err == nil
}
