package model

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port      int          `yaml:"port"`
	Debug     bool         `yaml:"debug" default:"false"`
	Telegram  TelegramConf `yaml:"telegram"`
	Templates Templates    `yaml:"templates"`
}

type TelegramConf struct {
	Token           string          `yaml:"token"`
	Trigger         string          `yaml:"trigger"`
	Aliases         Aliases         `yaml:"aliases"`
	TemplateToChats TemplateToChats `yaml:"templateToChats"`
}

type Templates map[string]string
type TemplateToChats map[string][]int64
type Aliases map[string]Alias
type Alias map[string]string

func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
