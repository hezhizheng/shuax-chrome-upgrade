package main

import (
	"encoding/json"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"shuax-chrome-auto-update/service"
	"time"
)

// 1、爬取 https://assets.shuax.com/ 的页面 获取最新版的chrome版本
//2、与本地chrome当前版本比较，大于当前版本则下载到本地、解压(询问提示)
//3、覆盖旧版数据，老版本App重命名为App2
//4、删除下载的文件

type Config struct {
	AutoDownload       string `json:"auto_download"`
	AutoUpdate         string `json:"auto_update"`
	DeleteDownloadFile string `json:"delete_download_file"`
	IntervalMin        string `json:"interval_min"`
	LocalChromePath    string `json:"local_chrome_path"`
}

var _config Config

func init() {
	initConfig()
	initLog()
}

func main() {
	configStr := viper.Get(`app`)
	jsonStr, e := json.Marshal(configStr)
	if e != nil {
		log.Error("json Marshal error  ", e)
	}
	json.Unmarshal(jsonStr, &_config)



	f := &service.FileInfo{
		FileDir: _config.LocalChromePath + "\\App\\",
	}
	localVersionName := service.GetLocalVersionName(f)

	chromeFileName,latestVersionName := service.GetLatestVersionName(_config.AutoDownload)

	fmt.Println("1111111",chromeFileName, latestVersionName, localVersionName)

	service.DownloadChrome("87.0.4280.89",localVersionName,chromeFileName)

	tickerRun()

	for {
	}
}

func tickerRun() {

	ticker := time.NewTicker(time.Minute * 1)

	i := 0
	go func() {
		for { //循环
			<-ticker.C
			i++
			fmt.Println("i =", i)
			//if i == 5 {
			//	ticker.Stop()
			//}
		}
	}()
}

func initConfig() {
	viper.SetConfigType("json") // 设置配置文件的类型
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Println("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}
}

func initLog() {
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	path := "./logs/"
	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	  `WithMaxAge 设置文件清理前的最长保存时间`
	  `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	// 下面配置日志每隔 1天 转一个新文件，保留最近 1周 的日志文件，多余的自动清理掉。
	LinkName := path + "shuax-chrome-auto-update.log"

	writer, _ := rotatelogs.New(
		//path+".%Y%m%d%H%M",
		path+"go-crontab-%Y-%m-%d.log",
		rotatelogs.WithLinkName(LinkName),
		rotatelogs.WithMaxAge(time.Duration(604800)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
	)
	log.SetOutput(writer)
}