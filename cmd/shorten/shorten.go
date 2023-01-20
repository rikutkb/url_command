package shorten

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rikutkb/url_command.git/cmd/abstract"
)

var _ abstract.IFetchCommand = &ShortFetchComand{}

type ShortFetchComand struct {
	urlPairs map[string]string
}

func (sfc *ShortFetchComand) GetData(ctx context.Context, url string) error {
	return nil
}
func (sfc ShortFetchComand) WriteData() error {
	return nil
}

func CreateShortUrl(ctx context.Context, url string, fetcher IFetchShUrl) (shortUrl string, err error) {
	// 10秒でタイムアウトを行う。
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	request, err := fetcher.CreateReq(url)
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
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
