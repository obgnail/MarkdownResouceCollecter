package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/obgnail/MarkdownResouceCollecter/config"
	"github.com/obgnail/MarkdownResouceCollecter/utils"
)

const (
	CollectLocalPicture         = "CollectLocalPicture"
	CollectNetWorkPicture       = "CollectNetWorkPicture"
	UploadLocalPictureToNetWork = "UploadLocalPictureToNetWork"
	ExportMarkdown              = "ExportMarkdown"
)

func New(cfg *config.Config) *BaseHandler {
	handler := &BaseHandler{Config: cfg}
	if err := handler.setDefaultStrategies(); err != nil {
		log.Fatal(err)
	}
	return handler
}

func (h *BaseHandler) setDefaultStrategies() error {
	mapTypeNameToStrategy := map[string]Strategy{
		CollectLocalPicture:         &CollectLocalPictureStrategy{},
		CollectNetWorkPicture:       &CollectNetWorkPictureStrategy{},
		UploadLocalPictureToNetWork: &UploadLocalPictureToNetWorkStrategy{},
		ExportMarkdown:              &ExportMarkdownStrategy{},
	}
	for _, name := range h.Config.Strategies {
		s := mapTypeNameToStrategy[name]
		if s == nil {
			return fmt.Errorf("[ERROR] No Such Strategy: %s", name)
		}
		h.AppendStrategy(s)
	}
	return nil
}

func (h *BaseHandler) ListStrategies() []Strategy {
	return h.strategies
}

func (h *BaseHandler) ClearStrategies() {
	h.strategies = []Strategy{}
}

func (h *BaseHandler) AppendStrategy(s Strategy) {
	h.strategies = append(h.strategies, s)
}

func (h *BaseHandler) Run() error {
	fmt.Println("---------------- Start ----------------")
	if err := h.Collect(); err != nil {
		return nil
	}
	if err := h.BaseAdjust(); err != nil {
		return err
	}
	if err := h.ExecuteStrategies(); err != nil {
		return err
	}
	if err := h.Rewrite(); err != nil {
		return err
	}
	fmt.Println("---------------- END ----------------")
	h.Report()
	fmt.Printf("\nPLEASE CHECK DIR: %s", h.NewMarkdownDirPath)
	return nil
}

func (h *BaseHandler) Collect() error {
	if err := h.CollectMarkdownFiles(); err != nil {
		return err
	}
	if err := h.CollectPictureFiles(); err != nil {
		return err
	}
	return nil
}

func (h *BaseHandler) BaseAdjust() error {
	allPic := make(map[string]struct{})
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			pic.SetBelongFile(file)
			if err := pic.EscapePath(); err != nil {
				return err
			}
			if err := pic.SetAbsPath(); err != nil {
				return err
			}

			if _, exist := allPic[pic.RealName]; exist {
				pic.ResetRealNameByMD5()
			}
			allPic[pic.RealName] = struct{}{}
		}
	}
	return nil
}

func (h *BaseHandler) ExecuteStrategies() error {
	for _, s := range h.strategies {
		if err := s.Adjust(h); err != nil {
			return err
		}
		if err := s.Extra(h); err != nil {
			return err
		}
	}
	return nil
}

func (h *BaseHandler) Rewrite() error {
	if err := utils.Mkdir(h.NewMarkdownDirPath); err != nil {
		return err
	}
	for _, file := range h.Files {
		if err := file.RewriteMarkdownFile(); err != nil {
			return fmt.Errorf("[Error] Rewrite Markdown File: %s", err)
		}
	}
	return nil
}

func (h *BaseHandler) Report() {
	h.SetTrashBin()
	fmt.Println("\n============ TrashBin ============")
	for idx1, file := range h.TrashBin {
		fmt.Printf("--- %d. %s\n", idx1+1, file.OriginPath)
		for _, pic := range file.Pictures {
			fmt.Printf("------  No.%d %s\n", pic.LineIndex, pic.OriginMatch)
		}
		fmt.Println()
	}
	fmt.Println("==================================")
}

func (h *BaseHandler) CollectMarkdownFiles() error {
	filePaths, err := utils.WalkDir(h.MarkdownDirPath, h.MarkdownFileSuffix)
	if err != nil {
		return fmt.Errorf("[ERROR] Walk MarkdownDirPath %s", err)
	}

	markDownFiles := make([]*MarkdownFile, len(filePaths))
	for idx, fp := range filePaths {
		markDownFiles[idx] = BuildMarkdown(fp, h.MarkdownDirPath, h.NewMarkdownDirPath)
	}
	h.Files = markDownFiles
	return nil
}

func (h *BaseHandler) CollectPictureFiles() error {
	for _, markdownFile := range h.Files {
		pics, err := GetPictureInFile(markdownFile.OriginPath)
		if err != nil {
			return fmt.Errorf("[ERROR] Get Picture In File %s", err)
		}
		markdownFile.Pictures = pics
	}
	return nil
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
				file.OriginDir,
				file.OriginPath,
				file.NewDir,
				file.NewPath,
				trashPics,
			}
			h.TrashBin = append(h.TrashBin, f)
		}
	}
}

func BuildMarkdown(originFilePath, originDirPath, newDirPath string) *MarkdownFile {
	markdownFile := new(MarkdownFile)
	markdownFile.OriginPath = originFilePath
	markdownFile.OriginDir = filepath.Dir(originFilePath)
	markdownFile.NewPath = strings.Replace(originFilePath, originDirPath, newDirPath, 1)
	markdownFile.NewDir = filepath.Dir(markdownFile.NewPath)
	return markdownFile
}

func BuildPicture(oldMatch, showName, oldPath string, lineIdx int) *Picture {
	pic := new(Picture)
	pic.ShowName = showName
	pic.OriginPath = oldPath
	pic.OriginMatch = oldMatch
	pic.LineIndex = lineIdx
	pic.RealName = path.Base(pic.OriginPath)

	if strings.HasPrefix(pic.OriginPath, "http") {
		pic.FromNet = true
	}
	return pic
}

func (f *MarkdownFile) RewriteMarkdownFile() error {
	fi, err := os.Open(f.OriginPath)
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
		if picIdx != len(f.Pictures) && bytes.Contains(line, []byte(f.Pictures[picIdx].OriginMatch)) {
			// pic可能在同一行
			for picIdx != len(f.Pictures) && bytes.Contains(line, []byte(f.Pictures[picIdx].OriginMatch)) {
				line = bytes.ReplaceAll(line, []byte(f.Pictures[picIdx].OriginMatch), []byte(f.Pictures[picIdx].NewMatch))
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
	fmt.Println("Rewrote:", f.OriginPath)
	return nil
}

func (p *Picture) ResetRealNameByMD5() {
	p.RealName = utils.MD5(p.AbsPath)[:8] + "-" + p.RealName
}

func (p *Picture) SetBelongFile(f *MarkdownFile) {
	p.BelongFile = f
}

func (p *Picture) EscapePath() error {
	// url 解码
	unescapePath, err := url.QueryUnescape(p.OriginPath)
	if err != nil {
		return fmt.Errorf("[WARN] Unescape Abs OriginPath %s", err)
	} else {
		p.OriginPath = unescapePath
	}
	return nil
}

func (p *Picture) SetAbsPath() error {
	if p.FromNet {
		p.IsExist = true
		p.AbsPath = p.OriginPath
		return nil
	}

	picAbsPath, err := GetAbsPath(p.BelongFile.OriginDir, p.OriginPath)
	if err != nil {
		p.AbsPath = p.OriginPath
		p.IsExist = false
		return fmt.Errorf("[WARN] Get Pic Abs OriginPath %s", err)
	}
	if exist, _ := utils.PathExists(picAbsPath); !exist {
		p.AbsPath = p.OriginPath
		p.IsExist = false
		return fmt.Errorf("[WARN]: Cant Find Pic File:%s, Match:%s", p.BelongFile.OriginPath, p.OriginMatch)
	}
	p.AbsPath = picAbsPath
	p.IsExist = true
	return nil
}

func MovePicturesToResourceDir(h *BaseHandler) error {
	if err := utils.Mkdir(h.ResourceDirPath); err != nil {
		return err
	}
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			if pic.FromNet || !pic.IsExist {
				continue
			}
			if exist, _ := utils.PathExists(pic.NewPath); exist {
				continue
			}
			err := utils.CopyFile(pic.AbsPath, pic.NewPath)
			if err != nil {
				fmt.Println(fmt.Errorf("[Error] Copy: %s, File:%s, Match:%s\n", err, file.OriginPath, pic.OriginMatch))
				continue
			}
		}
	}
	return nil
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
	// ./assets/XXX
	if strings.HasPrefix(relPath, "./") {
		return GetAbsPath(basePath, relPath[2:])
		// ../assets/XXX
	} else if strings.HasPrefix(relPath, "../") {
		return GetAbsPathCore(basePath, relPath)
		// /assets/XXX
	} else if strings.HasPrefix(relPath, "/") {
		return relPath, nil
		// assets/XXX
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
		// relPath[3:]: 去除../
		return GetAbsPathCore(parent, relPath[3:])
	} else {
		return "", fmt.Errorf("[ERROR] Get Rel OriginPath. Base:%s Rel:%s", basePath, relPath)
	}
}
