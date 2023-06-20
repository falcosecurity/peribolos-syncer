package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra/doc"

	"github.com/falcosecurity/peribolos-syncer/cmd"
)

const (
	docsDir      = "docs"
	fileTemplate = `---
title: %s
---	

`
)

var (
	filePrepender = func(filename string) string {
		title := strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(filename, "_", " "), ".md"), fmt.Sprintf("%s/", docsDir))
		return fmt.Sprintf(fileTemplate, title)
	}
	linkHandler = func(filename string) string {
		if filename == cmd.CommandName+".md" {
			return "_index.md"
		}
		return filename
	}
)

func main() {
	if err := doc.GenMarkdownTreeCustom(
		cmd.New(),
		docsDir,
		filePrepender,
		linkHandler,
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err := os.Rename(path.Join(docsDir, cmd.CommandName+".md"), path.Join(docsDir, "_index.md"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
