package strategy

import "github.com/obgnail/MarkdownResouceCollecter/process"

// Strategy 使用不同的策略对资源进行处理
type Strategy interface {
	// Adjust 每个策略都可以调整md资源
	Adjust(h *process.BaseHandler) error
	// Extra 每个策略都可以执行一些额外的任务
	Extra(h *process.BaseHandler) error
}
