package handler

import "github.com/obgnail/MarkdownResouceCollecter/config"

const (
	PictureGrammar = `!\[(.+?)\]\((.+?)\)`
)

// ![ones-framework](assets/aaa.jpg)
type Picture struct {
	ShowName    string        // md中的图片文件名, ones-framework
	RealName    string        // 真实的名字, aaa.jpg
	OriginPath  string        // 用于获取图片文件, assets/aaa.jpg
	OriginMatch string        // 从文件中匹配出来的string,用于后续替换, ![ones-framework](assets/aaa.jpg)
	LineIndex   int           // 所在文件行数
	FromNet     bool          // 是不是网络图片
	IsExist     bool          // 是否存在
	AbsPath     string        // 绝对路径,用于找到图片
	BelongFile  *MarkdownFile // 所属文件
	NewPath     string        // 用于生成新的图片
	NewMatch    string        // 未来用于替换文件中的OldMatch
}

type MarkdownFile struct {
	OriginDir  string // 文件目录
	OriginPath string // 文件路径
	NewDir     string // 新的文件目录
	NewPath    string // 新的文件路径
	Pictures   []*Picture
}

type Handler interface {
	// AppendStrategy 添加资源处理策略
	AppendStrategy(s *Strategy)
	// Collect 收集资源信息
	Collect() error
	// BaseAdjust 基础调整
	BaseAdjust()
	// ExecuteStrategies 执行策略
	ExecuteStrategies() error
	// Rewrite 重写资源信息到md文件中
	Rewrite() error
	// Report 输出执行报告
	Report()
}

type BaseHandler struct {
	Files      []*MarkdownFile
	TrashBin   []*MarkdownFile
	strategies []Strategy

	*config.Config
}

// Strategy 使用不同的策略对资源进行处理,更多地类似于中间件
type Strategy interface {
	// BeforeRewrite 在写入收集完资源,写入新的markdown文件前
	BeforeRewrite(h *BaseHandler) error
	// AfterRewrite 在写入新的markdown文件之后
	AfterRewrite(h *BaseHandler) error
}
