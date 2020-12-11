package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"shuax-chrome-auto-update/service"
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
}

func main() {
	configStr := viper.Get(`app`)
	jsonStr, e := json.Marshal(configStr)
	if e != nil {
		log.Fatal("json Marshal error  ", e)
	}
	json.Unmarshal(jsonStr, &_config)

	fmt.Printf("欢迎使用shuax_chrome_update工具\n")
	fmt.Printf("当前定义的本地chrome的安装路径为：" + _config.LocalChromePath + "\n")
	fmt.Printf("请根据提示输入相关指令进行操作\n")
	fmt.Printf("是否检查更新？1：是 2：否\n")

	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		line := input.Text()

		fmt.Printf("输入了：" + line + "\n")

		if line == "1" {

			f := &service.FileInfo{
				FileDir: _config.LocalChromePath + "\\App\\",
			}
			localVersionName := service.GetLocalVersionName(f)

			chromeFileName, latestVersionName := service.GetLatestVersionName(_config.AutoDownload)

			if service.CompareVersion(latestVersionName, localVersionName) == 1 {
				fmt.Printf("当前本地chrome的版本为：" + localVersionName + "，" + "最新chrome版本为：" + latestVersionName + " 是否进行升级？1：是 2：否\n")
			} else {
				fmt.Printf("当前本地chrome的版本为：" + localVersionName + "，" + "最新chrome版本为：" + latestVersionName + " 无需升级\n")
				break
			}
			var (
				isUpdate string
				isDelete string
			)
			fmt.Scanln(&isUpdate)
			fmt.Printf("输入了：" + isUpdate + "\n")

			if isUpdate != "1" {
				break
			}
			fmt.Printf("升级中，请等待，此过程中请不要做任何输入。\n")
			service.DownloadChrome(latestVersionName, localVersionName, chromeFileName)

			fmt.Printf("升级成功，是否删除下载/解压的文件？（建议先检查是否升级成功在执行此操作！！！）1：是 2：否\n")
			fmt.Scanln(&isDelete)
			fmt.Printf("输入了：" + isDelete + "\n")
			if isDelete != "1" {
				break
			}
			fmt.Printf("文件删除中......\n")
			break

		} else {
			break
		}

		// 输入bye时 结束
		if line == "exit" {
			break
		}
	}

	return
	f := &service.FileInfo{
		FileDir: _config.LocalChromePath + "\\App\\",
	}
	localVersionName := service.GetLocalVersionName(f)

	chromeFileName, latestVersionName := service.GetLatestVersionName(_config.AutoDownload)

	fmt.Println("1111111", chromeFileName, latestVersionName, localVersionName)

	service.DownloadChrome(latestVersionName, localVersionName, chromeFileName)

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
