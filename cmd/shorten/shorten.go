package shorten

import "net/http"

func CreateShortUrl(url string, fetcher IFetchShUrl) (shortUrl string, err error) {
	request, err := fetcher.CreateReq(url)
	if err != nil {
		return "", err
	}
	client := new(http.Client)
	resp, err := client.Do(request)
	if err != nil {
		return "", err
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
