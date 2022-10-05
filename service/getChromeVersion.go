package service

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const ShuaxHost = "https://chrome.noki.eu.org"
const CheckoutVersionApi = "https://chrome.noki.eu.org/api/c/info"

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

	cVersion := ""
	cUrl := ""

	c.OnResponse(func(resp *colly.Response) {

		//log.Println(string(resp.Body))

		var resInfo map[string]interface{}
		json.Unmarshal(resp.Body, &resInfo)

		//log.Println(resInfo)
		cVersion = resInfo["win_stable_x64"].(map[string]interface{})["version"].(string)
		cUrl = resInfo["win_stable_x64"].(map[string]interface{})["chromewithgc"].(string)
		fileName = cUrl

		//log.Println(cVersion, cUrl)

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

	//c.OnHTML(".fb-n", func(e *colly.HTMLElement) {
	//	if e.Index == 2 {
	//		fileName = e.Text
	//	}
	//})

	c.OnXML(`//*[@id="windows-x64"]/div/div[2]/div/div/div/div/p[6]/a`, func(element *colly.XMLElement) {
		log.Println(element)
		fileName = element.Text
	})

	proxyUrl := viper.GetString(`app.proxy_url`)

	if proxyUrl != "" {
		c.SetProxy("http://127.0.0.1:7890")
	}
	//visitError := c.Visit(ShuaxHost)
	visitError := c.Visit(CheckoutVersionApi)

	if visitError != nil {
		log.Println("访问" + CheckoutVersionApi + "失败")
		panic(visitError)
	}
	c.Wait()

	//version := ""

	// GoogleChrome_X64_87.0.4280.88_shuax.com.7z
	// https://chrome.noki.eu.org/download/106.0.5249.91_chrome_stable_x64.zip
	if fileName != "" {
		FStrSplit := strings.Split(fileName, "https://chrome.noki.eu.org/download/")[1]
		//version = strings.Split(FStrSplit, "_chrome_stable_x64")[0]
		fileName = FStrSplit // 106.0.5249.91_chrome_stable_x64.zip
	}

	//return fileName, version
	//log.Println(fileName)
	return fileName, cVersion
}
