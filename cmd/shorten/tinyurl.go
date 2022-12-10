package shorten

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

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
