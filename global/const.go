package global

import (
	"fmt"
	"os"
	"path/filepath"
)

var DirPath *Path

type Path struct {
	ResourceDirPath    string
	MarkdownDirPath    string
	ToolDirPath        string
	RootDirPath        string
	NewMarkdownDirPath string
}

func init() {
	ToolDirPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("Get pwd Error: %s\n", err)
	}
	rootDirPath := filepath.Dir(ToolDirPath)
	DirPath = &Path{
		RootDirPath:        filepath.Dir(ToolDirPath),
		ToolDirPath:        filepath.Join(rootDirPath, Cfg.ToolDirName),
		MarkdownDirPath:    filepath.Join(rootDirPath, Cfg.MarkdownDirName),
		ResourceDirPath:    filepath.Join(rootDirPath, Cfg.ResourceDirName),
		NewMarkdownDirPath: filepath.Join(rootDirPath, Cfg.NewMarkdownDirName),
	}
}
