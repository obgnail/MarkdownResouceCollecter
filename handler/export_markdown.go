package handler

// ExportMarkdownStrategy 导出需要的md文件/文件夹
type ExportMarkdownStrategy struct{}

func (s *ExportMarkdownStrategy) BeforeRewrite(h *BaseHandler) error {
	return nil
}

func (s *ExportMarkdownStrategy) AfterRewrite(h *BaseHandler) error {
	return nil
}
