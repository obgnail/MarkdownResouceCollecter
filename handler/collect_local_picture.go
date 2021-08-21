package handler

import (
	"fmt"
	"path/filepath"
)

// CollectLocalPictureStrategy 收集本地图片,将其放到其中一个统一的地方
type CollectLocalPictureStrategy struct{}

func (s *CollectLocalPictureStrategy) Adjust(h *BaseHandler) error {
	return s.SetLocalPictureNewMatchAndNewPath(h)
}

func (s *CollectLocalPictureStrategy) Extra(h *BaseHandler) error {
	return MovePicturesToResourceDir(h)
}

func (s *CollectLocalPictureStrategy) SetLocalPictureNewMatchAndNewPath(h *BaseHandler) error {
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			// 网络图片和不存在的图片保持不变
			if pic.FromNet || !pic.IsExist {
				pic.NewPath = pic.OriginPath
				pic.NewMatch = pic.OriginMatch
				continue
			}

			pic.NewPath = filepath.Join(h.ResourceDirPath, pic.RealName)

			if h.LocalPictureUseAbsPath {
				pic.NewMatch = fmt.Sprintf("![%s](%s)", pic.ShowName, pic.NewPath)
			} else {
				relPath, err := filepath.Rel(file.NewDir, h.ResourceDirPath)
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
