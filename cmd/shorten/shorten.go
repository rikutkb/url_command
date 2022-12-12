package shorten

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func CreateShortUrl(url string, fetcher IFetchShUrl) (shortUrl string, err error) {
	// 10秒でタイムアウトを行う。
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	request, err := fetcher.CreateReq(url)
	if err != nil {
		return "", fmt.Errorf("リクエストの作成に失敗しました。:%s", err)
	}
	client := new(http.Client)
	resp, err := client.Do(request.WithContext(ctx))
	if err != nil {
		// レスポンスヘッダー取得に10秒以上かかった場合
		return "", fmt.Errorf("ヘッダーの取得に失敗しました。:%s", err)
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
