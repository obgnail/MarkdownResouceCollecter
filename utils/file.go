package utils

import (
	"bytes"
	"fmt"
	. "github.com/obgnail/MarkdownResouceCollecter/global"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CheckExecuteDirExist() error {
	exist1, err1 := PathExists(DirPath.ToolDirPath)
	exist2, err2 := PathExists(DirPath.MarkdownDirPath)
	exist3, err3 := PathExists(DirPath.ResourceDirPath)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err2
	}
	if !exist1 {
		return fmt.Errorf("[ERROR] No Dir %s", DirPath.ToolDirPath)
	}
	if !exist2 {
		return fmt.Errorf("[ERROR] No Dir %s", DirPath.MarkdownDirPath)
	}
	if !exist3 {
		return fmt.Errorf("[ERROR] No Dir %s", DirPath.ResourceDirPath)
	}
	return nil
}

func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}

		return nil
	})

	return files, nil
}

func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Read Dir %s", err)
	}

	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}

	return files, nil
}

func PathExists(dirPth string) (bool, error) {
	_, err := os.Stat(dirPth)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("[ERROR] Path %s", err)
}

func Mkdir(dirPth string) error {
	exist, err := PathExists(dirPth)
	if err != nil {
		return fmt.Errorf("[ERROR] Path %s", err)
	}

	if !exist {
		fmt.Printf("No Dir![%v]\n", dirPth)
		err := os.MkdirAll(dirPth, os.ModePerm)
		if err != nil {
			return fmt.Errorf("[ERROR] Mkdir Failed %s", err)
		} else {
			fmt.Printf("Mkdir Success!\n")
		}
	}
	return nil
}

func CopyFile(originalFilePath, copiedFilePath string) error {
	originalFile, err := os.Open(originalFilePath)
	if err != nil {
		return fmt.Errorf("[ERROR] Open Origin File %s", err)
	}
	defer originalFile.Close()

	newFile, err := os.Create(copiedFilePath)
	if err != nil {
		return fmt.Errorf("[ERROR] Create File %s", err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, originalFile)
	if err != nil {
		return fmt.Errorf("[ERROR] Write File %s", err)
	}
	fmt.Printf("Copied %s.\n", originalFilePath)

	err = newFile.Sync()
	if err != nil {
		return fmt.Errorf("[ERROR] Sync File %s", err)
	}
	return nil
}

func Request(url string) ([]byte, error) {
	method := "GET"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.Close()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Close Writer %s", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, fmt.Errorf("[ERROR] New Request %s", err)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Request %s", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[ERROR] Request Status: %s", res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil || len(body) == 0 {
		return nil, fmt.Errorf("[ERROR] Read Response Body %s", err)
	}
	return body, nil
}

func Exit() {
	var data string
	for {
		fmt.Println("\npress exit to exit")
		fmt.Scanf("%s", &data)
		if data == "exit" {
			break
		}
	}
}
