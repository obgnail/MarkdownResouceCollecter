package handler

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/obgnail/MarkdownResouceCollecter/utils"
)

const (
	limit  = 5
	weight = 1
)

// CollectNetWorkPictureStrategy 收集本地图片,将其放到其中一个统一的地方
type CollectNetWorkPictureStrategy struct{}

func (s *CollectNetWorkPictureStrategy) BeforeRewrite(h *BaseHandler) error {
	s.ResolveDuplicateNameConflict(h)
	return PullNetWorkPictures(h)
}

func (s *CollectNetWorkPictureStrategy) AfterRewrite(h *BaseHandler) error { return nil }

// 避免网络图片和本地图片重名,所以网络图片一律添加md5前缀
func (s *CollectNetWorkPictureStrategy) ResolveDuplicateNameConflict(h *BaseHandler) {
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			if !pic.FromNet || !pic.IsExist {
				continue
			}
			pic.ResetRealNameByMD5()
		}
	}
}

func PullNetWorkPictures(h *BaseHandler) error {
	var wg sync.WaitGroup
	s := semaphore.NewWeighted(limit)

	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			if !pic.FromNet || !pic.IsExist {
				continue
			}
			wg.Add(1)
			go func(pic *Picture) {
				defer wg.Done()
				if err := s.Acquire(context.Background(), weight); err != nil {
					fmt.Println(err)
				}
				if err := pic.StoreNetWorkPicture(h.NewResourceRootDirPath, h.LocalPictureUseAbsPath); err != nil {
					fmt.Println(err)
				}
				s.Release(weight)
			}(pic)
		}
	}
	wg.Wait()
	return nil
}

func (p *Picture) StoreNetWorkPicture(resourceDirPath string, isLocalPictureUseAbsPath bool) error {
	newFilePath := filepath.Join(resourceDirPath, p.RealName)
	relPath, err := filepath.Rel(p.BelongFile.NewDir, resourceDirPath)
	if err != nil {
		return fmt.Errorf("[ERROR] Get RelPath: %s", err)
	}

	// 当文件不存在本地时,发起网络请求
	if !p.IsExistInLocal(resourceDirPath) {
		p.IsExist = false
		body, err := utils.Request(p.AbsPath)
		if err != nil || body == nil {
			return fmt.Errorf("[WARN] Cant Pull NetWork File %s, Match:%s", err, p.OriginMatch)
		}
		if err := utils.Mkdir(resourceDirPath); err != nil {
			return err
		}
		if err := ioutil.WriteFile(newFilePath, body, 0644); err != nil {
			return fmt.Errorf("[Error] Write Picture File: %s", err)
		}
		fmt.Println("Pull Picture Success:", p.OriginPath)
	}

	p.FromNet = false
	p.IsExist = true
	p.AbsPath = newFilePath

	if isLocalPictureUseAbsPath {
		p.NewPath = p.AbsPath
		p.NewMatch = fmt.Sprintf("![%s](%s)", p.ShowName, p.AbsPath)
	} else {
		picRelPath := filepath.Join(relPath, p.RealName)
		p.NewPath = picRelPath
		p.NewMatch = fmt.Sprintf("![%s](%s)", p.ShowName, picRelPath)
	}

	return nil
}

func (p *Picture) IsExistInLocal(resourceDirPath string) bool {
	filePath := filepath.Join(resourceDirPath, p.RealName)
	if existInLocal, _ := utils.PathExists(filePath); existInLocal {
		return true
	}
	return false
}
