package handler

import (
	"fmt"
	"path/filepath"
)

// CollectLocalPictureStrategy 收集本地图片,将其放到其中一个统一的地方
type CollectLocalPictureStrategy struct{}

func (s *CollectLocalPictureStrategy) BeforeRewrite(h *BaseHandler) error {
	if err := s.SetLocalPictureNewMatchAndNewPath(h); err != nil {
		return err
	}
	if err := MovePicturesToResourceDir(h); err != nil {
		return err
	}
	return nil
}

func (s *CollectLocalPictureStrategy) AfterRewrite(h *BaseHandler) error { return nil }

func (s *CollectLocalPictureStrategy) SetLocalPictureNewMatchAndNewPath(h *BaseHandler) error {
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			// 网络图片和不存在的图片保持不变
			if pic.FromNet || !pic.IsExist {
				pic.NewPath = pic.OriginPath
				pic.NewMatch = pic.OriginMatch
				continue
			}

			pic.NewPath = filepath.Join(h.NewResourceRootDirPath, pic.RealName)

			if h.LocalPictureUseAbsPath {
				pic.NewMatch = fmt.Sprintf("![%s](%s)", pic.ShowName, pic.NewPath)
			} else {
				relPath, err := filepath.Rel(file.NewDir, h.NewResourceRootDirPath)
				if err != nil {
					fmt.Println(fmt.Errorf("[ERROR] Get RelPath: %s", err))
					pic.NewPath = pic.OriginPath
					pic.NewMatch = pic.OriginMatch
					continue
				}
				// NOTE: 这里使用的是RealName,RealName已经去重,避免命名冲突
				picRelPath := filepath.Join(relPath, pic.RealName)
				pic.NewMatch = fmt.Sprintf("![%s](%s)", pic.ShowName, picRelPath)
			}
		}
	}
	return nil
}
