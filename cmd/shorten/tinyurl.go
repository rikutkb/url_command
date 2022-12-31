package shorten

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type TinyURL struct {
	apiKey string
}
type TinyURLReq struct {
	URL string `json:"url"`
}
type TinyURLRespData struct {
	Url string `json:"tiny_url"`
}
type TinyURLResp struct {
	TinyURLRespData `json:"data"`
}
type TinyURLRespFail struct {
	Errors []string `json:"errors"`
}

const TINYURL_API_ENV = "TINYURL_API_KEY"

var _ IFetchShUrl = TinyURL{}

// 環境変数からAPIキーをセット
func (t TinyURL) Init() (err error) {
	apiKey := os.Getenv(TINYURL_API_ENV)
	if apiKey == "" {
		return errors.New("APIキーがセットされていません。")
	}
	t.apiKey = apiKey
	return nil
}

// Requstの認証方法はそれぞれ異なる可能性があるため、生成メソッドは別にする
func (t TinyURL) CreateReq(baseUrl string) (req *http.Request, err error) {
	body_json, _ := json.Marshal(TinyURLReq{URL: baseUrl})
	API_KEY := t.apiKey
	serviceUrl := "https://api.tinyurl.com/create"

	method := "POST"
	body := bytes.NewBuffer(body_json)
	req, err = http.NewRequest(method, serviceUrl, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request : %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+API_KEY)

	return req, nil
}
func (t TinyURL) ParseResp(resp *http.Response) (shUrl string, err error) {
	bodyText, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		var resp = TinyURLResp{}
		if err := json.Unmarshal(bodyText, &resp); err != nil {
			return "", fmt.Errorf("failed to unmarshal json :%s", err)
		}
		return resp.TinyURLRespData.Url, nil
	} else {
		var resp = TinyURLRespFail{}
		if err := json.Unmarshal(bodyText, &resp); err != nil {
			return "", fmt.Errorf("failed to unmarshal json :%s", err)
		}

		return "", fmt.Errorf("通信中に不具合がありました。:%s", resp.Errors)
	}

}
