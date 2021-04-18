package handler

const (
	PictureGrammar = `!\[(.+?)\]\((.+?)\)`

	CollectLocalPicture         = "CollectLocalPicture"
	CollectNetWorkPicture       = "CollectNetWorkPicture"
	UploadLocalPictureToNetWork = "UploadLocalPictureToNetWork"
	ExportMarkdown              = "ExportMarkdown"
)

type Picture struct {
	ShowName   string        // md中的图片文件名
	RealName   string        // 真实的名字
	OldPath    string        // 用于获取图片文件
	OldMatch   string        // 从文件中匹配出来的string，用于后续替换
	LineIndex  int           // 所在文件行数
	FromNet    bool          // 是不是网络图片
	IsExist    bool          // 是否存在
	AbsPath    string        // 绝对路径,用于找到图片
	BelongFile *MarkdownFile // 所属文件
	NewPath    string        // 用于生成新的图片
	NewMatch   string        // 未来用于替换文件中的OldMatch
}

type MarkdownFile struct {
	Dir      string // 文件目录
	Path     string // 文件路径
	NewDir   string // 新的文件目录
	NewPath  string // 新的文件路径
	Pictures []*Picture
}

type BaseHandler struct {
	Files    []*MarkdownFile
	TrashBin []*MarkdownFile
}

type Handler interface {
	Collect() error
	BaseAdjust()
	Adjust()
	Extra() error
	Rewrite() error
	Report()
}

type CollectLocalPictureHandler struct {
	*BaseHandler
}

type CollectNetWorkPictureHandler struct {
	*BaseHandler
}

type UploadLocalPictureToNetWorkHandler struct {
	*BaseHandler
}

type ExportMarkdownHandler struct {
	*BaseHandler
}
