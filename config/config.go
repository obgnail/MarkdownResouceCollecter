package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	MarkdownDirPath        string   `yaml:"markdownDirPath"`
	ResourceDirPath        string   `yaml:"resourceDirPath"`
	NewMarkdownDirPath     string   `yaml:"newMarkdownDirPath"`
	Strategies             []string `yaml:"strategies"`
	MarkdownFileSuffix     string   `yaml:"markdownFileSuffix"`
	LocalPictureUseAbsPath bool     `yaml:"localPictureUseAbsPath"`
}

func InitConfigFromYaml(path string) *Config {
	c := new(Config)
	yamlFile, err := ioutil.ReadFile(path)
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
