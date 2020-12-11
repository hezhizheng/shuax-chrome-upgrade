package service

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func DownloadChrome(latestVersionName, localVersionName, chromeFileName string) {

	//AutoUpdate := viper.Get(`app.auto_update`)
	DeleteDownloadFile := viper.Get(`app.delete_download_file`)

	needDownload := false

	if latestVersionName != "" && localVersionName != "" {

		if  compareVersion(latestVersionName,localVersionName) == 1 {
			needDownload = true
		}

	}else{
		return
	}

	url := ShuaxHost+"/"+chromeFileName
	path := viper.Get(`app.local_chrome_path`)

	if needDownload{
		fmt.Println("开始下载")
		downloadFile(url,path.(string)+"\\", chromeFileName,func(length, downLen int64) {})
	}


	if AutoState == DeleteDownloadFile {

	}
}


func isFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(info)
		return false
	}
	if filesize == info.Size() {
		fmt.Println("安装包已存在！", info.Name(), info.Size(), info.ModTime())
		return true
	}
	del := os.Remove(filename)
	if del != nil {
		fmt.Println(del)
	}
	return false
}

// https://blog.csdn.net/SHIXINGYA/article/details/88951782
func downloadFile(url string, localPath string,chromeFileName string, fb func(length, downLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	tmpFilePath := localPath + ".download"
	fmt.Println(tmpFilePath)
	//创建一个http client
	client := new(http.Client)
	//client.Timeout = time.Second * 60 //设置超时时间
	//get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	//读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	//if isFileExist(localPath, fsize) {
	//	return err
	//}
	fmt.Println("fsize", fsize)
	//创建文件
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		//没有错误了快使用 callback
		fb(fsize, written)
	}
	fmt.Println(err)
	if err == nil {
		file.Close()
		err = os.Rename(tmpFilePath, localPath+chromeFileName)
		fmt.Println(err)
	}
	return err
}


// https://leetcode-cn.com/problems/compare-version-numbers/solution/golangshi-xian-by-he-qing-ping/
func compareVersion(version1 string, version2 string) int {
	versionA:= strings.Split(version1,".")
	versionB:= strings.Split(version2,".")

	for i:= len(versionA);i<4;i++{
		versionA = append(versionA,"0")
	}
	for i:= len(versionB);i<4;i++{
		versionB = append(versionB,"0")
	}
	for i:= 0;i<4;i++{
		version1,_:= strconv.Atoi(versionA[i])
		version2,_:= strconv.Atoi(versionB[i])
		if version1 == version2{
			continue
		}else if version1>version2{
			return 1
		}else{
			return -1
		}
	}
	return 0
}
