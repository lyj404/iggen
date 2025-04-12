package main

import (
	"iggen/internal/cli"
	"iggen/internal/generator"
	"iggen/internal/github"
)

func main() {
	ghClient := github.NewGitHubClient()
	fileGen := generator.NewGitignoreGenerator()
	cli.Run(ghClient, fileGen)
}
