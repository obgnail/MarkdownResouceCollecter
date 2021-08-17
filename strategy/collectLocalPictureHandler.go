package strategy

import (
	"github.com/obgnail/MarkdownResouceCollecter/process"
)

// CollectLocalPictureStrategy 收集本地图片,将其放到其中一个统一的地方
type CollectLocalPictureStrategy struct{}

func (s *CollectLocalPictureStrategy) Adjust(h *process.BaseHandler) error {
	//s.SetLocalPictureNewMatchAndNewPath()
	return nil
}

func (s *CollectLocalPictureStrategy) Extra(h *process.BaseHandler) error {
	//if err := s.MovePicturesToResourceDir(); err != nil {
	//	return err
	//}
	//
	return nil
}

//
//func (h *CollectLocalPictureStrategy) SetLocalPictureNewMatchAndNewPath() {
//	for _, file := range h.Files {
//		for _, pic := range file.Pictures {
//			if pic.FromNet || !pic.IsExist {
//				pic.NewPath = pic.OldPath
//				pic.NewMatch = pic.OldMatch
//			} else {
//				pic.NewPath = filepath.Join(DirPath.ResourceDirPath, pic.RealName)
//				relPath, err := filepath.Rel(file.Dir, DirPath.ResourceDirPath)
//				if err != nil {
//					fmt.Println(fmt.Errorf("[ERROR] Get RelPath: %s", err))
//					pic.NewPath = pic.OldPath
//					pic.NewMatch = pic.OldMatch
//					continue
//				}
//
//				if Cfg.DoesLocalPictureUseAbsPath {
//					pic.NewMatch = fmt.Sprintf("![%s](%s)", pic.ShowName, pic.AbsPath)
//				} else {
//					relPath = filepath.Join(relPath, pic.RealName)
//					pic.NewMatch = fmt.Sprintf("![%s](%s)", pic.ShowName, relPath)
//				}
//			}
//		}
//	}
//}
