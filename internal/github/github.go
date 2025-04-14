package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GitHubClient struct {
	client *http.Client
}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		client: &http.Client{Timeout: 3 * time.Second},
	}
}

func NewGitHubClientWithProxy(proxy string) *GitHubClient {
	// 创建自定义Transport以实现代理配置
	transport := &http.Transport{}
	// 代理配置逻辑
	if proxy != "" {
		// 解析代理地址（此处假设为HTTP代理，需添加http://前缀）
		proxyURL, _ := url.Parse("http://" + proxy)
		// 设置代理方法（自动为每个请求添加代理）
		transport.Proxy = http.ProxyURL(proxyURL)
	}
	// 返回带代理配置的客户端实例
	return &GitHubClient{
		client: &http.Client{
			Transport: transport,       // 注入自定义Transport
			Timeout:   3 * time.Second, // 设置全局请求超时
		},
	}
}

// ListTemplates获取GitHub所有.gitignore模板
func (c *GitHubClient) ListTemplates(ctx context.Context) ([]string, error) {
	// 创建带上下文的HTTP GET请求
	req, _ := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.github.com/repos/github/gitignore/contents",
		nil,
	)

	// 执行HTTP请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %w", err)
	}
	// 确保响应体关闭
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回异常状态码: %d", resp.StatusCode)
	}

	// 定义临时结构体解析JSON响应
	var files []struct{ Name string }
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		// JSON解析错误处理
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 过滤并格式化模板名称
	templates := make([]string, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name, ".gitignore") {
			// 移除.gitignore后缀
			templates = append(templates, strings.TrimSuffix(f.Name, ".gitignore"))
		}
	}
	return templates, nil
}

// GetNormalizedName对模板名称进行小写
func (c *GitHubClient) GetNormalizedName(ctx context.Context, inputName string) (string, error) {
	// 获取所有可用模板
	templates, err := c.ListTemplates(ctx)
	if err != nil {
		return "", err
	}
	// 将输入转换为小写进行模糊匹配
	lowerInput := strings.ToLower(inputName)
	for _, t := range templates {
		if strings.ToLower(t) == lowerInput {
			return t, nil // 返回仓库中实际存在的大小写格式
		}
	}
	return "", fmt.Errorf("找不到对应模板")
}

// GetTemplate获取GitHub具体.gitignore模板内容
func (c *GitHubClient) GetTemplate(ctx context.Context, name string) ([]byte, error) {
	// 先获取规范化名称
	normalized, err := c.GetNormalizedName(ctx, name)
	if err != nil {
		return nil, err
	}
	// 构造GitHub raw内容URL
	url := fmt.Sprintf(
		"https://raw.githubusercontent.com/github/gitignore/main/%s.gitignore",
		normalized, // 使用规范化后的名称
	)

	// 创建带上下文的GET请求
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("下载模板失败: %w", err)
	}
	defer resp.Body.Close()

	// 验证HTTP响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// 读取全部响应内容并返回
	return io.ReadAll(resp.Body)
}
