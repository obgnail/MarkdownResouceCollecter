package main

import (
	"fmt"
	"github.com/obgnail/MarkdownResouceCollecter/config"
	"github.com/obgnail/MarkdownResouceCollecter/handler"
)

func main() {
	cfg := config.InitConfigFromYaml("config.yaml")
	h := handler.New(cfg)
	if err := h.Run(); err != nil {
		fmt.Println(err)
	}
}
