package handler

// UploadLocalPictureToNetWorkStrategy 将本地图片上传至网络
type UploadLocalPictureToNetWorkStrategy struct{}

func (s *UploadLocalPictureToNetWorkStrategy) BeforeRewrite(h *BaseHandler) error {
	return nil
}

func (s *UploadLocalPictureToNetWorkStrategy) AfterRewrite(h *BaseHandler) error {
	return nil
}
