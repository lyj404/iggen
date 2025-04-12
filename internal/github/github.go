package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type GitHubClient struct {
	client *http.Client
}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

// ListTemplates获取GitHub所有.gitignore模板
func (c *GitHubClient) ListTemplates(ctx context.Context) ([]string, error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.github.com/repos/github/gitignore/contents",
		nil,
	)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回异常状态码: %d", resp.StatusCode)
	}

	var files []struct{ Name string }
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	templates := make([]string, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name, ".gitignore") {
			templates = append(templates, strings.TrimSuffix(f.Name, ".gitignore"))
		}
	}
	return templates, nil
}

// GetNormalizedName对模板名称进行小写
func (c *GitHubClient) GetNormalizedName(ctx context.Context, inputName string) (string, error) {
	templates, err := c.ListTemplates(ctx)
	if err != nil {
		return "", err
	}
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
	url := fmt.Sprintf(
		"https://raw.githubusercontent.com/github/gitignore/main/%s.gitignore",
		normalized, // 使用规范化后的名称
	)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("下载模板失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}
