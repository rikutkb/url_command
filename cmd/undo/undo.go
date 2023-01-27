package undo

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/rikutkb/url_command.git/cmd/abstract"
)

func init() {

}

var reg = regexp.MustCompile(`(<meta property="og:url" content=")(.*)(".*/>)`)

var _ abstract.IFetchCommand = &UndoFetchCommand{}

type UndoFetchCommand struct {
	urlPairs map[string]string
}

func NewUndoFetchCommand() *UndoFetchCommand {
	return &UndoFetchCommand{urlPairs: map[string]string{}}
}

func (ufc *UndoFetchCommand) GetData(ctx context.Context, url string) error {
	client := &http.Client{}
	// goの標準のリダイレクト機能をオーバライド
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	req, err := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)
	if err != nil {
		return fmt.Errorf("リクエスト作成中にエラーが起きました: %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("通信中にエラーが起きました: %s", err)
	}
	if resp.StatusCode != http.StatusMovedPermanently {
		return fmt.Errorf("通信中にエラーが起きました: status code :%d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	parsed := reg.FindStringSubmatch(string(body))
	if len(parsed) < 3 {
		return fmt.Errorf("レスポンスのパースに失敗しました。")
	}
	ufc.urlPairs[url] = parsed[2]
	return nil
	//	<meta property="og:url" content="https://xxxxxxxxx" />
}

func (ufc UndoFetchCommand) WriteData(reqUrls []string) error {
	return nil
}
