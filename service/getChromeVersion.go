package service

import (
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const ShuaxHost = "https://assets.shuax.com"

type FileInfo struct {
	FileDir string
	Version string
}

// 获取版本
func (f *FileInfo) GetLocalVersion() (err error) {
	rd, e := ioutil.ReadDir(f.FileDir)

	if e != nil {
		log.Println("目录读取失败", err, f.FileDir)
		return nil
	}

	// 第一个文件夹名字即版本号
	f.Version = rd[0].Name()

	return nil
}

func GetLocalVersionName(f *FileInfo) string {
	f.GetLocalVersion()
	return f.Version
}

func GetLatestVersionName() (string, string) {
	fileName := ""

	c := colly.NewCollector(
		colly.Async(true),
	)

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("请求", r.URL, "...")
	})

	retryCount := 0
	c.OnError(func(res *colly.Response, err error) {
		log.Println("Something went wrong:", err)
		if retryCount < 3 {
			retryCount += 1
			_retryErr := res.Request.Retry()
			log.Println("retry wrong:", _retryErr)
		}
	})

	c.OnHTML(".fb-n", func(e *colly.HTMLElement) {
		if e.Index == 2 {
			fileName = e.Text
		}
	})

	visitError := c.Visit(ShuaxHost)

	if visitError != nil {
		log.Println("访问" + ShuaxHost + "失败")
		panic(visitError)
	}
	c.Wait()

	version := ""

	// GoogleChrome_X64_87.0.4280.88_shuax.com.7z
	if fileName != "" {
		FStrSplit := strings.Split(fileName, "_X64_")[1]
		version = strings.Split(FStrSplit, "_shuax")[0]
	}

	return fileName, version
}
