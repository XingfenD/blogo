package loader

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

func LoadIcons() (map[string]string, error) {
	icons := make(map[string]string)
	err := filepath.Walk("static/icon", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".svg" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			icons[name] = string(content)
		}
		return nil
	})
	return icons, err
}

func extractFrontMatter(content string) string {
	start := strings.Index(content, "+++\n")
	end := strings.LastIndex(content, "\n+++")

	if start == -1 || end == -1 {
		return ""
	}

	start += len("+++\n")
	return content[start:end]
}

type FrontMatter struct {
	Title       string `toml:"title"`
	Description string `toml:"description"`
	Date        string `toml:"date"`
}

func LoadPages() map[string]FrontMatter {
	ret := make(map[string]FrontMatter)
	entries, err := os.ReadDir("content/page")
	if err != nil {
		log.Println("Error reading directory: ", err)
		return nil
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		indexMd := "content/page/" + entry.Name() + "/_index.md"
		pf, err := os.ReadFile(indexMd)
		if err != nil {
			log.Println(err)
			continue
		}
		content := string(pf)
		tomlContent := extractFrontMatter(content)
		if tomlContent == "" {
			log.Println("No front matter found in the file: " + indexMd)
			continue
		}

		var frontMatter FrontMatter
		_, err = toml.Decode(tomlContent, &frontMatter)
		if err != nil {
			log.Println("Error parsing TOML: ", err)
			return nil
		}
		ret[entry.Name()] = frontMatter
	}

	return ret
}
