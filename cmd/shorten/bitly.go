package shorten

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

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
