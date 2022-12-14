package undo

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func init() {

}

func GetRedirect(ctx context.Context, url string) (longUrl string, err error) {
	client := &http.Client{}
	// goの標準のリダイレクト機能をオーバライド
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	req, err := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)
	if err != nil {
		return "", fmt.Errorf("リクエスト作成中にエラーが起きました: %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("通信中にエラーが起きました: %s", err)
	}
	if resp.StatusCode != http.StatusMovedPermanently {
		return "", fmt.Errorf("通信中にエラーが起きました: status code :%d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	reg := regexp.MustCompile(`(<a href=")(.*)(">.*)(</a>)`)
	parsed := reg.FindStringSubmatch(string(body))
	if len(parsed) < 4 {
		return "", fmt.Errorf("レスポンスのパースに失敗しました。")
	}
	return parsed[2], nil
	// <a href="https://xxxxxxxxxxxxxxxxxxx">https://xxxxxxxxxxxxxxxxxxx</a>
	// <a href="https://xxxxxxxxxxxxxxxxxxx">
	// https://xxxxxxxxxxxxxxxxxxx
	// </a>
}
