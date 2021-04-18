package handler

import (
	"bufio"
	"bytes"
	"fmt"
	. "github.com/obgnail/MarkdownResouceCollecter/global"
	"github.com/obgnail/MarkdownResouceCollecter/utils"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func (h *BaseHandler) Collect() error {
	// Collect MarkdownFiles
	filePaths, err := utils.WalkDir(DirPath.MarkdownDirPath, Cfg.MarkdownFileSuffix)
	if err != nil {
		return fmt.Errorf("[ERROR] Walk MarkdownDirPath %s", err)
	}

	markDownFiles := make([]*MarkdownFile, len(filePaths))
	for idx, fp := range filePaths {
		markDownFiles[idx] = BuildMarkdown(fp)
	}
	h.Files = markDownFiles

	// Collect PictureFiles
	for _, markdownFile := range h.Files {
		pics, err := GetPictureInFile(markdownFile.Path)
		if err != nil {
			return fmt.Errorf("[ERROR] Get Picture In File %s", err)
		}
		markdownFile.Pictures = pics
	}
	return nil
}

func (h *BaseHandler) BaseAdjust() {
	h.SetPictureBelongFile()
	h.SetLocalPictureRealName()
	h.SetLocalPictureAbsPath()
}

func (h *BaseHandler) Report() {
	h.SetTrashBin()
	fmt.Println("\n============ TrashBin ============")
	for idx1, file := range h.TrashBin {
		fmt.Printf("--- %d. %s\n", idx1+1, file.Path)
		for _, pic := range file.Pictures {
			fmt.Printf("------  No.%d %s\n", pic.LineIndex, pic.OldMatch)
		}
		fmt.Println()
	}
	fmt.Println("==================================")
}

func (h *BaseHandler) SetPictureBelongFile() {
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			pic.BelongFile = file
		}
	}
}

func (h *BaseHandler) SetLocalPictureRealName() {
	allPic := make(map[string]byte)
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			if _, exist := allPic[pic.RealName]; exist {
				pic.RealName = fmt.Sprintf("%d%s", time.Now().UnixNano(), path.Ext(pic.RealName))
			}
			allPic[pic.RealName] = 1
		}
	}
}

func (h *BaseHandler) SetLocalPictureAbsPath() {
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			if pic.FromNet {
				pic.IsExist = true
				pic.AbsPath = pic.OldPath
			} else {
				oldPath :=  pic.OldPath
				// url 解码
				unescapePath, err := url.QueryUnescape(pic.OldPath)
				if err != nil {
					fmt.Println(fmt.Errorf("[WARN] Unescape Abs Path %s", err))
				} else {
					oldPath = unescapePath
				}

				absPath, err := GetAbsPath(file.Dir, oldPath)
				if err != nil {
					fmt.Println(fmt.Errorf("[WARN] Get Abs Path %s", err))
					pic.AbsPath = pic.OldPath
					pic.IsExist = false
					continue
				}
				if exist, _ := utils.PathExists(absPath); !exist {
					fmt.Println(fmt.Errorf("[WARN]: Cant Find File:%s, Match:%s", file.Path, pic.OldMatch))
					pic.AbsPath = pic.OldPath
					pic.IsExist = false
					continue
				}
				pic.AbsPath = absPath
				pic.IsExist = true
			}
		}
	}
}

func (h *BaseHandler) SetTrashBin() {
	for _, file := range h.Files {
		var trashPics []*Picture
		for _, pic := range file.Pictures {
			if !pic.IsExist {
				trashPics = append(trashPics, pic)
			}
		}
		if trashPics != nil {
			f := &MarkdownFile{
				file.Dir,
				file.Path,
				file.NewDir,
				file.NewPath,
				trashPics,
			}
			h.TrashBin = append(h.TrashBin, f)
		}
	}
}

func (h *BaseHandler) MovePicturesToResourceDir() error {
	if err := utils.Mkdir(DirPath.ResourceDirPath); err != nil {
		return err
	}
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			if pic.FromNet || !pic.IsExist {
				continue
			}
			err := utils.CopyFile(pic.AbsPath, pic.NewPath)
			if err != nil {
				fmt.Println(fmt.Errorf("[Error] Copy: %s, File:%s, Match:%s\n", err, file.Path, pic.OldMatch))
				continue
			}
		}
	}
	return nil
}

func (h *BaseHandler) Rewrite() error {
	if err := utils.Mkdir(DirPath.NewMarkdownDirPath); err != nil {
		return err
	}
	for _, file := range h.Files {
		if err := file.RewriteMarkdownFile(); err != nil {
			return fmt.Errorf("[Error] Rewrite Markdown File: %s", err)
		}
	}
	return nil
}

func (f *MarkdownFile) RewriteMarkdownFile() error {
	fi, err := os.Open(f.Path)
	if err != nil {
		return fmt.Errorf("[Error] Open Markdown file Error: %s\n", err)
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	output := make([]byte, 0)
	picIdx := 0
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if picIdx != len(f.Pictures) && bytes.Contains(line, []byte(f.Pictures[picIdx].OldMatch)) {
			// pic可能在同一行
			for picIdx != len(f.Pictures) && bytes.Contains(line, []byte(f.Pictures[picIdx].OldMatch)) {
				line = bytes.ReplaceAll(line, []byte(f.Pictures[picIdx].OldMatch), []byte(f.Pictures[picIdx].NewMatch))
				picIdx++
			}
		}
		output = append(output, line...)
		output = append(output, []byte("\n")...)
	}

	if err := utils.Mkdir(f.NewDir); err != nil {
		return err
	}
	if err := ioutil.WriteFile(f.NewPath, output, 0644); err != nil {
		return fmt.Errorf("[Error] Write Markdown File: %s", err)
	}
	fmt.Println("Rewrote:", f.Path)
	return nil
}

func BuildPicture(oldMatch, showName, oldPath string, lineIdx int) *Picture {
	pic := new(Picture)
	pic.ShowName = showName
	pic.OldPath = oldPath
	pic.OldMatch = oldMatch
	pic.LineIndex = lineIdx

	pic.RealName = path.Base(pic.OldPath)
	if strings.HasPrefix(pic.OldPath, "http") {
		pic.FromNet = true
	}
	return pic
}

func BuildMarkdown(path string) *MarkdownFile {
	markdownFile := new(MarkdownFile)
	markdownFile.Path = path
	markdownFile.Dir = filepath.Dir(path)

	markdownFile.NewPath = strings.Replace(path, Cfg.MarkdownDirName, Cfg.NewMarkdownDirName, 1)
	markdownFile.NewDir = strings.Replace(markdownFile.Dir, Cfg.MarkdownDirName, Cfg.NewMarkdownDirName, 1)
	return markdownFile
}

func GetPictureInFile(filePath string) ([]*Picture, error) {
	fi, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("[Error] Open Markdown File: %s", err)
	}
	defer fi.Close()

	var pictures []*Picture
	br := bufio.NewReader(fi)

	lineIdx := 0
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		lineIdx++
		pics, err := FindPicturesInLine(string(line))
		if err != nil {
			return nil, fmt.Errorf("[Error] Find Picture In Markdown File: %s", err)
		}
		if len(pics) == 0 {
			continue
		}
		for _, pic := range pics {
			p := BuildPicture(pic["oldMatch"], pic["showName"], pic["picPath"], lineIdx)
			pictures = append(pictures, p)
		}
	}
	return pictures, nil
}

func FindPicturesInLine(line string) ([]map[string]string, error) {
	PictureRegexp, err := regexp.Compile(PictureGrammar)
	if err != nil {
		return nil, err
	}
	pictures := PictureRegexp.FindAllStringSubmatch(line, -1)
	if pictures == nil {
		return nil, nil
	}
	pics := make([]map[string]string, len(pictures))
	for idx, r := range pictures {
		p := map[string]string{
			"oldMatch": r[0],
			"showName": r[1],
			"picPath":  r[2],
		}
		pics[idx] = p
	}
	return pics, nil
}

func GetAbsPath(basePath, relPath string) (absPath string, err error) {
	if strings.HasPrefix(relPath, "./") {
		return GetAbsPath(basePath, relPath[2:])
	} else if strings.HasPrefix(relPath, "../") {
		return GetAbsPathCore(basePath, relPath)
	} else if strings.HasPrefix(relPath, "/") {
		return relPath, nil
	} else {
		return filepath.Join(basePath, relPath), nil
	}
}

func GetAbsPathCore(basePath, relPath string) (string, error) {
	if !strings.HasPrefix(relPath, "../") {
		return filepath.Join(basePath, relPath), nil
	}

	parent := filepath.Dir(basePath)
	if parent != "." {
		return GetAbsPathCore(parent, relPath[3:])
	} else {
		return "", fmt.Errorf("[ERROR] Get Rel Path. Base:%s Rel:%s", basePath, relPath)
	}
}
