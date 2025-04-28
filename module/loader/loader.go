package loader

import (
	"os"
	"path/filepath"
	"strings"
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
