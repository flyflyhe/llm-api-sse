package fileHelper

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

type DownloadFileResStruct struct {
	Output string
	Err    error
}

func DownloadFile(u string, savePath string, gbk bool) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	// 从路径中获取最后一个元素
	lastPath := path.Base(parsedURL.Path)

	parentPath := savePath
	now := time.Now()
	for _, p := range []string{now.Format("20060102"), strconv.Itoa(int(now.Unix()))} {
		parentPath = path.Join(parentPath, p)
		if !PathExist(parentPath) {
			if err := os.Mkdir(parentPath, 755); err != nil {
				return "", err
			}
		}
	}

	var output string
	if gbk {
		gbkLastPath, err := Utf8ToGbk([]byte(lastPath))
		if err != nil {
			return "", err
		}
		output = path.Join(parentPath, string(gbkLastPath))
	} else {
		output = path.Join(parentPath, lastPath)
	}

	log.Println(output)
	// 发送HTTP请求获取响应
	response, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// 创建输出文件
	file, err := os.Create(output)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 将响应体复制到输出文件中
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return output, nil
}

func DownloadFileChan(urlList []string, parentPath string, gbk bool) chan map[string]DownloadFileResStruct {
	c := make(chan map[string]DownloadFileResStruct)
	go func() {
		m := make(map[string]DownloadFileResStruct)
		for _, v := range urlList {
			output, err := DownloadFile(v, parentPath, gbk)
			m[v] = DownloadFileResStruct{
				Output: output,
				Err:    err,
			}
		}
		c <- m
	}()

	return c
}

func PathExist(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// UTF-8 转 GBK

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
