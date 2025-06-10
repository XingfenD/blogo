package loader

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger logrus.Logger

func LoadLogger(logLevel int) {
	Logger.Out = os.Stdout
	Logger.Level = logrus.Level(logLevel)
	Logger.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
	Logger.Info("Logger initialized")
}

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
