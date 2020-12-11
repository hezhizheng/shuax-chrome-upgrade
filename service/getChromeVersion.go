package service

import (
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

const AutoState = "true"
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

func GetLatestVersionName(AutoDownload string) (string,string) {
	fileName := ""

	if AutoState == AutoDownload {
		c := colly.NewCollector(
			colly.Async(true),
		)

		c.WithTransport(&http.Transport{
			DisableKeepAlives: true,
		})

		c.OnRequest(func(r *colly.Request) {
			log.Println("Visiting", r.URL)
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

		log.Println(visitError)

		c.Wait()
	}

	version := ""

	// GoogleChrome_X64_87.0.4280.88_shuax.com.7z
	if fileName != ""{
		FStrSplit := strings.Split(fileName, "_X64_")[1]
		version = strings.Split(FStrSplit, "_shuax")[0]
	}

	return fileName , version
}
