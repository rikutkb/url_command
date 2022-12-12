package shorten

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
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
	body_json, _ := json.Marshal(TinyURLReq{URL: baseUrl})
	API_KEY := t.apiKey
	serviceUrl := "https://api.tinyurl.com/create"

	method := "POST"
	body := bytes.NewBuffer(body_json)
	req, err = http.NewRequest(method, serviceUrl, body)
	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+API_KEY)

	return req, err
}
func (t *TinyURL) ParseResp(resp *http.Response) (shUrl string, err error) {
	bodyText, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		var resp = TinyURLResp{}
		if err := json.Unmarshal(bodyText, &resp); err != nil {
			return "", err
		}
		return resp.TinyURLRespData.Url, nil
	} else {
		var resp = TinyURLRespFail{}
		if err := json.Unmarshal(bodyText, &resp); err != nil {
			return "", err
		}
		for e := range resp.Errors {
			fmt.Print(e)
		}

		return "", fmt.Errorf("通信中に不具合がありました。:%s", resp.Errors)
	}

}
