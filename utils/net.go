package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

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
