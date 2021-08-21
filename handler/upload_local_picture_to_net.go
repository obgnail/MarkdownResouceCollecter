package handler

// UploadLocalPictureToNetWorkStrategy 将本地图片上传至网络
type UploadLocalPictureToNetWorkStrategy struct{}

func (s *UploadLocalPictureToNetWorkStrategy) Adjust(h *BaseHandler) error {
	return nil
}

func (s *UploadLocalPictureToNetWorkStrategy) Extra(h *BaseHandler) error {
	return nil
}
