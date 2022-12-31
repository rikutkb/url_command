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

type BitlyReq struct {
	Domain string `json:"domain"`
	URL    string `json:"long_url"`
}
type BitlyResp struct {
	CreatedAt string `json:"created_at"`
	ShortURL  string `json:"link"`
}
type BitlyRespErrors struct {
	Field     string `json:"field"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}
type BitlyRespFail struct {
	Message     string            `json:"message"`
	Description string            `json:"description"`
	Resource    string            `json:"resource"`
	Errors      []BitlyRespErrors `json:"errors"`
}

type Bitly struct {
	apiKey string
}

const BITLY_API_ENV = "BIT_API_KEY"

var _ IFetchShUrl = Bitly{}

func (b Bitly) GetApiKey() string {
	return b.apiKey
}

func (b Bitly) Init() (err error) {
	apiKey := os.Getenv(BITLY_API_ENV)
	if apiKey == "" {
		return errors.New("APIキーがセットされていません。")
	}
	b.apiKey = apiKey
	return nil
}

func (b Bitly) CreateReq(baseUrl string) (req *http.Request, err error) {
	body_json, _ := json.Marshal(BitlyReq{Domain: "bit.ly", URL: baseUrl})
	BITLY_API_KEY := b.apiKey
	serviceUrl := "https://api-ssl.bitly.com/v4/shorten"

	method := "POST"
	body := bytes.NewBuffer(body_json)
	req, err = http.NewRequest(method, serviceUrl, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request : %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+BITLY_API_KEY)

	return req, nil
}
func (b Bitly) ParseResp(resp *http.Response) (shUrl string, err error) {

	bodyText, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		resp := BitlyResp{}
		if err := json.Unmarshal(bodyText, &resp); err != nil {
			return "", fmt.Errorf("failed to unmarshal json : %s", err)
		}
		return resp.ShortURL, nil
	} else {
		resp := BitlyRespFail{}
		if err := json.Unmarshal(bodyText, &resp); err != nil {
			return "", fmt.Errorf("failed to marshal json : %s", err)
		}
		return "", errors.New(resp.Errors[0].Message)
	}
}
