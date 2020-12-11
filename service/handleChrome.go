package service

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func DownloadChrome(latestVersionName, localVersionName, chromeFileName string) {

	//AutoUpdate := viper.Get(`app.auto_update`)

	needDownload := false

	if latestVersionName != "" && localVersionName != "" {

		if CompareVersion(latestVersionName, localVersionName) == 1 {
			needDownload = true
		}

	} else {
		panic("版本异常！")
	}

	url := ShuaxHost + "/" + chromeFileName
	path := viper.Get(`app.local_chrome_path`)

	filename := path.(string) + "\\" + chromeFileName

	if needDownload && !fileExists(filename) {
		fmt.Println("开始下载........")
		err := DownloadFile(filename, url)

		if err != nil {
			fmt.Println("下载文件" + url + "失败")
			panic(err)
		}

		fmt.Println("下载完成........")
	}

	if fileExists(filename) {

		// 先删除旧版本升级遗留的文件夹
		os.RemoveAll(path.(string) + "\\" + "App2")
		os.RemoveAll(path.(string) + "\\" + "chrome")

		fmt.Println("解压文件........")

		_, e1 := exec.Command("./7z.exe", "x", filename, "-o"+path.(string)).Output()

		if e1 != nil {
			fmt.Println("解压文件失败")
			panic(e1)
		}

		fmt.Println("解压完成")

		renameErr := os.Rename(path.(string)+"\\"+"App", path.(string)+"\\"+"App2")

		if renameErr != nil {
			fmt.Println("重命名文件失败")
			panic(renameErr)
		}

		e2 := copyDir(path.(string)+"\\"+"chrome\\App", path.(string)+"\\"+"App")

		if e2 != nil {
			fmt.Println("复制目录失败")
			panic(e2)
		}

		return
	}

	panic("升级失败")

}

func fileForCopyDir(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func copyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = fileForCopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func DownloadFile(filepath string, url string) error {

	// https://studygolang.com/articles/26441 进度条
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

// https://leetcode-cn.com/problems/compare-version-numbers/solution/golangshi-xian-by-he-qing-ping/
func CompareVersion(version1 string, version2 string) int {
	versionA := strings.Split(version1, ".")
	versionB := strings.Split(version2, ".")

	for i := len(versionA); i < 4; i++ {
		versionA = append(versionA, "0")
	}
	for i := len(versionB); i < 4; i++ {
		versionB = append(versionB, "0")
	}
	for i := 0; i < 4; i++ {
		version1, _ := strconv.Atoi(versionA[i])
		version2, _ := strconv.Atoi(versionB[i])
		if version1 == version2 {
			continue
		} else if version1 > version2 {
			return 1
		} else {
			return -1
		}
	}
	return 0
}
