package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const outputFile = "output.txt"

func main() {
	rootDir := "." // カレントディレクトリ
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("出力ファイル作成エラー:", err)
		return
	}
	defer out.Close()

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if ignoreDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldIgnorePath(path) {
			return nil
		}

		if isBinaryFile(path) || filepath.Base(path) == outputFile {
			return nil
		}

		lang := detectLanguage(path)
		if lang == "" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "```%s \n# %s\n\n", lang, path)
		out.Write(content)
		fmt.Fprintln(out, "```")
		fmt.Fprintln(out, "---")

		return nil
	})

	if err != nil {
		fmt.Println("エラー:", err)
	} else {
		fmt.Println("出力完了:", outputFile)
	}
}

var ignoreDirs = map[string]bool{
	".git":           true,
	"node_modules":   true,
	"vendor":         true,
	".vscode":        true,
	"dist":           true,
	"build":          true,
	".idea":          true,
	"tmp":            true,
	"uploads":        true,
	".next":          true,
	".github":        true,
	".pytest_cache":  true,
	"test_resources": true,
	"venv":           true,
}

func shouldIgnorePath(path string) bool {
	lower := strings.ToLower(path)
	return strings.Contains(lower, "min.js") ||
		strings.Contains(lower, ".lock") ||
		strings.Contains(lower, ".DS_Store")
}

func detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return "go"
	case ".toml":
		return "toml"
	case ".env":
		return "dotenv"
	case ".yml", ".yaml":
		return "yaml"
	case ".json":
		return "json"
	case ".md":
		return "markdown"
	case ".txt":
		return "text"
	case ".sh":
		return "bash"
	case ".Dockerfile", "dockerfile":
		return "Dockerfile"
	case ".gitignore":
		return ""
	case ".makefile", "makefile":
		return "makefile"
	case ".html":
		return "html"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".jsx":
		return "jsx"
	case ".tsx":
		return "tsx"
	case ".sql":
		return "sql"
	default:
		return ""
	}
}

var binaryExts = map[string]bool{
	".exe":  true,
	".dll":  true,
	".so":   true,
	".bin":  true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".pdf":  true,
	".zip":  true,
	".tar":  true,
	".gz":   true,
	".7z":   true,
	".mp3":  true,
	".mp4":  true,
	".avi":  true,
	".mov":  true,
	".wav":  true,
	".ico":  true,
}

func isBinaryFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return binaryExts[ext]
}
