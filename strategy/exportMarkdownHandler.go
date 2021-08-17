package strategy

import "github.com/obgnail/MarkdownResouceCollecter/process"

// ExportMarkdownStrategy 导出需要的md文件/文件夹
type ExportMarkdownStrategy struct{}

func (s *ExportMarkdownStrategy) Adjust(h *process.BaseHandler) error {
	return nil
}

func (s *ExportMarkdownStrategy) Extra(h *process.BaseHandler) error {
	return nil
}
