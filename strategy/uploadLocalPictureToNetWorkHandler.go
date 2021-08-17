package strategy

import "github.com/obgnail/MarkdownResouceCollecter/process"

// UploadLocalPictureToNetWorkStrategy 将本地图片上传至网络
type UploadLocalPictureToNetWorkStrategy struct{}

func (s *UploadLocalPictureToNetWorkStrategy) Adjust(h *process.BaseHandler) error {
	return nil
}

func (s *UploadLocalPictureToNetWorkStrategy) Extra(h *process.BaseHandler) error {
	return nil
}
