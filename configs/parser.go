package configs

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"slices"
)

func ParseYaml(filepath string, data interface{}) {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Error("Error during opening file \"%s\": %v", filepath, err.Error())
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	parser := yaml.NewDecoder(file)
	if err := parser.Decode(data); err != nil {
		slog.Error("Error during parsing config: %v", err.Error())
		panic(err)
	}
}

func ParseConfig(filepath string) *Config {
	var res Config
	ParseYaml(filepath, &res)
	return &res
}

func ParseDeadlineTasks(filepath string) DeadlineData {
	var res DeadlineData
	ParseYaml(filepath, &res)
	return res
}

func (config *Config) GetColorByCount(count int) string {
	ind := slices.IndexFunc(config.UnsolvedBorders, func(border *UnsolvedBorder) bool {
		return border.Count >= count
	})
	if ind < 0 {
		slog.Warn("Strange config: GetColorByCount returned -1; replacing it with len() - 1")
		ind = len(config.UnsolvedBorders) - 1
	}
	return config.UnsolvedBorders[ind].Color
}
