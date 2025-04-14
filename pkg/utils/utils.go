package utils

import (
	"strings"
	"unicode"
)

// FuzzySearch模糊搜索模板名称（支持部分匹配和相似字符匹配）
func FuzzySearch(templates []string, term string) []string {
	// 统一转为小写实现大小写不敏感搜索[^1]
	term = strings.ToLower(term)
	matches := make([]string, 0)

	// 遍历所有模板进行双阶段匹配
	for _, t := range templates {
		lower := strings.ToLower(t)

		// 第一阶段：直接包含匹配
		if strings.Contains(lower, term) {
			matches = append(matches, t)
			continue // 匹配成功则跳过后续判断[^2]
		}

		// 第二阶段：特殊字符处理后匹配
		if compareSimilar(lower, term) {
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
