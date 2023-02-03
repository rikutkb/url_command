package shorten

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rikutkb/url_command/cmd/abstract"
)

var _ abstract.IFetchCommand = &ShortFetchComand{}

type ShortFetchComand struct {
	urlPairs map[string]string
	Fecther  IFetchShUrl
}

func NewShortFetchCommand(fetcher IFetchShUrl) *ShortFetchComand {
	return &ShortFetchComand{Fecther: fetcher, urlPairs: map[string]string{}}
}
func (sfc *ShortFetchComand) GetData(ctx context.Context, url string) error {
	shortenUrl, err := CreateShortUrl(ctx, url, sfc.Fecther)
	if err != nil {
		return err
	}
	sfc.urlPairs[url] = shortenUrl
	return nil
}
func (sfc ShortFetchComand) WriteData(reqUrls []string) error {
	for i, url := range reqUrls {
		fmt.Fprintf(os.Stdout, sfc.urlPairs[url])
		if i+1 != len(reqUrls) {
			fmt.Fprintf(os.Stdout, ",")
		} else {
			fmt.Fprintln(os.Stdout, "")
		}
	}
	return nil
}

func CreateShortUrl(ctx context.Context, url string, fetcher IFetchShUrl) (shortUrl string, err error) {
	// 10秒でタイムアウトを行う。
	con, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	done := make(chan error, 1)
	urlChan := make(chan string, 1)
	request, err := fetcher.CreateReq(url)
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", fmt.Errorf("リクエストの作成に失敗しました。:%s", err)
	}
	client := new(http.Client)
	go func() {
		resp, err := client.Do(request.WithContext(con))
		if err != nil {
			// レスポンスヘッダー取得に10秒以上かかった場合
			done <- fmt.Errorf("ヘッダーの取得に失敗しました。:%s", err)
			return
		}
		url, err := fetcher.ParseResp(resp)
		urlChan <- url
		done <- err
		defer resp.Body.Close()
	}()
	if err := <-done; err != nil {
		return "", err
	}
	return <-urlChan, nil
}

type IFetchShUrl interface {
	// 環境変数からAPIキーをセット
	Init() (err error)
	ParseResp(resp *http.Response) (shUrl string, err error)
	CreateReq(baseUrl string) (req *http.Request, err error)
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
