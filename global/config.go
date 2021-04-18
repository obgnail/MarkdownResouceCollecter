package global

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Cfg *Config

type Config struct {
	ResourceDirName            string   `yaml:"resourceDirName"`
	MarkdownDirName            string   `yaml:"markdownDirName"`
	ToolDirName                string   `yaml:"toolDirName"`
	MarkdownFileSuffix         string   `yaml:"markdownFileSuffix"`
	NewMarkdownDirName         string   `yaml:"newMarkdownDirName"`
	Handlers                   []string `yaml:"handlers"`
	DoesLocalPictureUseAbsPath bool     `yaml:"doesLocalPictureUseAbsPath"`
}

func getConfig() *Config {
	c := new(Config)
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return c
}

func init() {
	Cfg = getConfig()
}
