package main

import (
	"github.com/lyj404/iggen/internal/cli"
	"github.com/lyj404/iggen/internal/generator"
	"github.com/lyj404/iggen/internal/github"
)

func main() {
	ghClient := github.NewGitHubClient()
	fileGen := generator.NewGitignoreGenerator()
	cli.Run(ghClient, fileGen)
}
