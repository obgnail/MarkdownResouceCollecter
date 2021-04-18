package handler

import (
	"fmt"
	. "github.com/obgnail/MarkdownResouceCollecter/global"
	"path/filepath"
)

func (h *CollectLocalPictureHandler) Adjust() {
	h.SetLocalPictureNewMatchAndNewPath()
}

func (h *CollectLocalPictureHandler) SetLocalPictureNewMatchAndNewPath() {
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			if pic.FromNet || !pic.IsExist {
				pic.NewPath = pic.OldPath
				pic.NewMatch = pic.OldMatch
			} else {
				pic.NewPath = filepath.Join(DirPath.ResourceDirPath, pic.RealName)
				relPath, err := filepath.Rel(file.Dir, DirPath.ResourceDirPath)
				if err != nil {
					fmt.Println(fmt.Errorf("[ERROR] Get RelPath: %s", err))
					pic.NewPath = pic.OldPath
					pic.NewMatch = pic.OldMatch
					continue
				}

				if Cfg.DoesLocalPictureUseAbsPath {
					pic.NewMatch = fmt.Sprintf("![%s](%s)", pic.ShowName, pic.AbsPath)
				} else {
					relPath = filepath.Join(relPath, pic.RealName)
					pic.NewMatch = fmt.Sprintf("![%s](%s)", pic.ShowName, relPath)
				}
			}
		}
	}
}

func (h *CollectLocalPictureHandler) Extra() error {
	if err := h.MovePicturesToResourceDir(); err != nil {
		return err
	}

	return nil
}
