package handler

import (
	"fmt"
	"github.com/obgnail/MarkdownResouceCollecter/utils"
	"io/ioutil"
	"regexp"
	"testing"
)

func TestCollectLocalPicture(t *testing.T) {
	line := "![platform-第 4 页 (1)](assets/platform-第 4 页 (1).jpg)"
	PictureRegexp, err := regexp.Compile(PictureGrammar)
	if err != nil {
		t.Fatal(err)
	}
	pictures := PictureRegexp.FindAllStringSubmatch(line, -1)
	if pictures == nil {
		t.Fatal(err)
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
}

func TestWriteFile(t *testing.T) {
	path := "/Users/heyingliang/myTemp/root/md.bak/InAction/Link/link.md"
	content := "this is content"
	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestMD5(t *testing.T) {
	f := "/Users/heyingliang/myTemp/root/md.bak/Learning/Celery/网站__Celery博客教程/网站__Celery博客教程.md" + "1563378187404.png"
	result := utils.MD5(f)[:8]+"-"+"1563378187404.png"
	fmt.Println(result)
}