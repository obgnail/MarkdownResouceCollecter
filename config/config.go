package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	OriginMarkdownRootPath string   `yaml:"originMarkdownRootPath"`
	NewResourceRootDirPath string   `yaml:"newResourceRootDirPath"`
	NewMarkdownRootPath    string   `yaml:"newMarkdownRootPath"`
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

func InitConfig(
	OriginMarkdownRootPath, NewResourceRootDirPath, NewMarkdownRootPath string,
	Strategies []string, MarkdownFileSuffix string, LocalPictureUseAbsPath bool,
) *Config {
	return &Config{
		OriginMarkdownRootPath: OriginMarkdownRootPath,
		NewResourceRootDirPath: NewResourceRootDirPath,
		NewMarkdownRootPath:    NewMarkdownRootPath,
		Strategies:             Strategies,
		MarkdownFileSuffix:     MarkdownFileSuffix,
		LocalPictureUseAbsPath: LocalPictureUseAbsPath,
	}
}
