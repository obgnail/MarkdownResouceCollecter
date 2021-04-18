package handler

import (
	"fmt"
	. "github.com/obgnail/MarkdownResouceCollecter/global"
	"github.com/obgnail/MarkdownResouceCollecter/utils"
	"io/ioutil"
	"path/filepath"
)

func (h *CollectNetWorkPictureHandler) Adjust() {

}

func (h *CollectNetWorkPictureHandler) Extra() error {
	if err := h.PullNetWorkPictures(); err != nil {
		return err
	}
	return nil
}

func (h *CollectNetWorkPictureHandler) PullNetWorkPictures() error {
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			if err := pic.StoreNetWorkPicture(); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func (p *Picture) StoreNetWorkPicture() error {
	if !p.FromNet || !p.IsExist {
		return nil
	}

	p.IsExist = false
	body, err := utils.Request(p.AbsPath)
	if err != nil || body == nil {
		return fmt.Errorf("[WARN] Cant Pull NetWork File %s, Match:%s", err, p.OldMatch)
	}
	if err := utils.Mkdir(DirPath.ResourceDirPath); err != nil {
		return err
	}

	newFilePath := filepath.Join(DirPath.ResourceDirPath, p.RealName)
	if err := ioutil.WriteFile(newFilePath, body, 0644); err != nil {
		return fmt.Errorf("[Error] Write Picture File: %s", err)
	}
	relPath, err := filepath.Rel(p.BelongFile.Dir, DirPath.ResourceDirPath)
	if err != nil {
		return fmt.Errorf("[ERROR] Get RelPath: %s", err)
	}

	fmt.Println("Pull Picture Success:", p.OldPath)

	p.FromNet = false
	p.IsExist = true
	p.AbsPath = newFilePath
	p.NewPath = newFilePath
	p.NewMatch = fmt.Sprintf("![%s](%s)", p.ShowName, filepath.Join(relPath, p.RealName))
	return nil
}
