// Copyright 2023 The Falco Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
