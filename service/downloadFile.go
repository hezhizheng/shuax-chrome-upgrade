package service

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const httpProxyUrl = "http://127.0.0.1:7890"

// @link https://studygolang.com/articles/26441
type writeCounter struct {
	Total uint64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc writeCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func DownloadFile(filepath string, furl string) error {
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	proxyUrl := viper.GetString(`app.proxy_url`)

	var httpclient = http.Client{}

	if proxyUrl != "" {
		ProxyURL, _ := url.Parse(httpProxyUrl)
		httpclient = http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(ProxyURL),
			},
		}
	}

	resp, err := httpclient.Get(furl)
	//resp, err := http.Get(furl)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()
	counter := &writeCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}
	fmt.Print("\n")
	out.Close()
	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}
