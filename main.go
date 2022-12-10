package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
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
	init_flag    string
	kind_flag    string
	token_flag   string
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
	flag.StringVar(&kind_flag, "k", "bitly", "使用サービスの指定を行います。")
	flag.StringVar(&token_flag, "t", "", "APIトークンの設定を行います。")
}

func main() {
	FlagInit()
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %s", err)
		}
		os.Exit(2)
	}
	// 短縮orQRコード,短縮urlを元に戻す
	if url_flag != "" {
		// TODO APIキーが設定できていない場合はエラーとして出力するようにする。
		// TODO interfaceとしてhttp通信部分を実装
		url := url_flag
		var fetcher = NewFecher(kind_flag)
		if err := fetcher.Init(); err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
		shortUrl, err := CreateShorUrl(url, fetcher)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
		fmt.Println(shortUrl)
	} else {
		fmt.Errorf("不正なコマンド入力です。")
		os.Exit(2)
	}

}

func CreateShorUrl(url string, fetcher IFetchShUrl) (shortUrl string, err error) {
	request, err := fetcher.CreateReq(url)
	if err != nil {
		return "", err
	}
	client := new(http.Client)
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return fetcher.ParseResp(resp)

}

type IFetchShUrl interface {
	// 環境変数からAPIキーをセット
	Init() (err error)
	ParseResp(resp *http.Response) (shUrl string, err error)
	CreateReq(baseUrl string) (req *http.Request, err error)
}
type TinyURL struct {
	apiKey string
}
type TinyURLReq struct {
	Domain    string `json:"domain"`
	URL       string `json:"url"`
	Alias     string `json:"alias"`
	Tags      string `json:"tags"`
	ExpiresAt string `json:"expires_at"`
}
type TinyURLRespData struct {
	Url string `json:"tiny_url"`
}
type TinyURLResp struct {
	TinyURLRespData `json:"data"`
}
type TinyURLErrResp struct {
	Errors []string `json:"errors"`
}

const TINYURL_API_ENV = "TINYURL_API_KEY"

// 環境変数からAPIキーをセット
func (t *TinyURL) Init() (err error) {
	apiKey := os.Getenv(TINYURL_API_ENV)
	if apiKey == "" {
		return errors.New("APIキーがセットされていません。")
	}
	t.apiKey = apiKey
	return nil
}

func (t *TinyURL) CreateReq(baseUrl string) (req *http.Request, err error) {
	body_json, _ := json.Marshal(TinyURLReq{Domain: "bit.ly", URL: baseUrl})
	API_KEY := t.apiKey
	serviceUrl := "https://api.tinyurl.com/create"

	method := "POST"
	body := bytes.NewBuffer(body_json)
	req, err = http.NewRequest(method, serviceUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+API_KEY)
	return req, err
}
func (t *TinyURL) ParseResp(resp *http.Response) (shUrl string, err error) {
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var tinyURLResp = TinyURLResp{}
		if err := json.Unmarshal(respBody, &tinyURLResp); err != nil {
			return "", err
		}
		return tinyURLResp.Url, nil
	} else {
		var TinyURLErrResp = TinyURLErrResp{}
		if err := json.Unmarshal(respBody, &TinyURLErrResp); err != nil {
			return "", err
		}
		return "", err
	}

}

type BitlyReq struct {
	Domain string `json:"domain"`
	URL    string `json:"long_url"`
}
type BitlyResp struct {
	CreatedAt string `json:"created_at"`
	ShortURL  string `json:"link"`
}

type Bitly struct {
	apiKey string
}

func NewFecher(service string) IFetchShUrl {
	// validaionは入ってきた際に行うので、ここではエラーハンドリングを行わない。
	switch service {
	case "bitly":
		return &Bitly{}
	case "TinyURL":
		return &TinyURL{}
	default:
		panic("not implemented")
	}

}

const BITLY_API_ENV = "BIT_API_KEY"

func (b *Bitly) Init() (err error) {
	apiKey := os.Getenv(BITLY_API_ENV)
	if apiKey == "" {
		return errors.New("APIキーがセットされていません。")
	}
	b.apiKey = apiKey
	return nil
}

func (b *Bitly) CreateReq(baseUrl string) (req *http.Request, err error) {
	body_json, _ := json.Marshal(BitlyReq{Domain: "bit.ly", URL: baseUrl})
	BITLY_API_KEY := b.apiKey
	serviceUrl := "https://api-ssl.bitly.com/v4/shorten"

	method := "POST"
	body := bytes.NewBuffer(body_json)
	req, err = http.NewRequest(method, serviceUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+BITLY_API_KEY)
	return req, err
}
func (b *Bitly) ParseResp(resp *http.Response) (shUrl string, err error) {

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var bitlyRespBody = BitlyResp{}
		if err := json.Unmarshal(respBody, &bitlyRespBody); err != nil {
			return "", err
		}
		return bitlyRespBody.ShortURL, nil
	} else {
		return "", errors.New("err")
	}
}
