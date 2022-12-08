package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	url_flag     string
	short_flag   bool
	service_flag string
	replace_flag bool
	undo_flag    bool
	set_flag     string
	qr_flag      bool
	file_flag    string
)

func FlagInit() {
	flag.BoolVar(&short_flag, "s", false, "urlの短縮を行います。")
	flag.BoolVar(&replace_flag, "r", false, "urlの置換を行います。")
	flag.BoolVar(&undo_flag, "u", false, "短縮urlを元に戻します。")
	flag.BoolVar(&qr_flag, "q", false, "短縮urlを元に戻します。")

	flag.StringVar(&url_flag, "url", "", "url")
	flag.StringVar(&service_flag, "service", "bitly", "urlの短縮を行います。")
	flag.StringVar(&set_flag, "set", "", "")
	flag.StringVar(&file_flag, "f", "", "")
}

type BitlyBody struct {
	Domain string `json:"domain"`
	URL    string `json:"long_url"`
}
type BitlyRespBody struct {
	CreatedAt string `json:"created_at"`
	ShortURL  string `json:"link"`
}

func main() {
	FlagInit()
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %s", err)
		}
		os.Exit(2)
	}
	// TODO APIキーが設定できていない場合はエラーとして出力するようにする。
	// TODO interfaceとしてhttp通信部分を実装
	url := url_flag
	fmt.Print(url)
	body_json, _ := json.Marshal(BitlyBody{Domain: "bit.ly", URL: url})
	BITLY_API_KEY := os.Getenv("BIT_API_KEY")
	method := "POST"
	serviceUrl := "https://api-ssl.bitly.com/v4/shorten"
	fmt.Print(url)
	body := bytes.NewBuffer(body_json)
	req, err := http.NewRequest(method, serviceUrl, body)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	req.Header.Set("Authorization", "bearer"+BITLY_API_KEY)
	dump, _ := httputil.DumpRequest(req, true)
	fmt.Print(string(dump))
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("error: %s", err)
		os.Exit(2)
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var bitlyRespBody = BitlyRespBody{}
	fmt.Print(string(respBody))
	if err := json.Unmarshal(respBody, &bitlyRespBody); err != nil {
		fmt.Errorf("error:%s", err)
		os.Exit(2)

	}
	fmt.Print(bitlyRespBody)
	fmt.Print(bitlyRespBody.ShortURL)

}
