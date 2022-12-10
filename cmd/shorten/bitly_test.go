package shorten_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/rikutkb/url_command.git/cmd/shorten"
	"github.com/stretchr/testify/assert"
)

func TestBitlyInit(t *testing.T) {

	tests := []struct {
		name    string
		wantErr error
		apiEnv  string
	}{
		{
			name:    "正常ケース",
			wantErr: nil,
			apiEnv:  "xxxxxxxxx",
		},
		{
			name:    "異常ケース:APIキーがセットされていない場合",
			wantErr: errors.New("APIキーがセットされていません。"),
			apiEnv:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(shorten.BITLY_API_ENV, tt.apiEnv)
			b := &shorten.Bitly{}
			if err := b.Init(); err != nil && errors.Is(err, tt.wantErr) {
				t.Errorf("Bitly.Init() error = %v, want %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.apiEnv, b.GetApiKey())
		})
	}
}

func TestBitlyCreateReq(t *testing.T) {

	type args struct {
		baseUrl string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "正常ケース",
			args:    args{baseUrl: "https://xxxxxxxxxxxx"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &shorten.Bitly{}
			_, err := b.CreateReq(tt.args.baseUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bitly.CreateReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBitlyParseResp(t *testing.T) {
	// 実際のjsonか構造帯に代入したものどっちの方が良い？
	// 構造体にするメリット: 短くてわかりやすい、urlを変数に入れることができる。
	// jsonのメリット実際のjsonを入れるため、仕様変更があったとしてもエラーが発見しやすい
	sucRespJson := []byte(`{
		"created_at": "2022-12-08T11:53:57+0000",
		"id": "bitly.is/3FAF6T2",
		"link": "https://bitly.is/3FAF6T2",
		"custom_bitlinks": [],
		"long_url": "https://dev.bitly.com/",
		"archived": false,
		"tags": [],
		"deeplinks": [],
		"references": {
			"group": "https://api-ssl.bitly.com/v4/groups/Bmc8b5iw5tz"
		}
	}`)
	failRespJson := []byte(`
	{
		"message": "INVALID_ARG_LONG_URL",
		"resource": "bitlinks",
		"description": "The value provided is invalid.",
		"errors": [
			{
				"field": "long_url",
				"error_code": "invalid"
			}
		]
	}`)
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name      string
		args      args
		wantShUrl string
		wantErr   bool
	}{
		{
			name: "正常ケース",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBuffer(sucRespJson)),
				},
			},
			wantShUrl: "https://bitly.is/3FAF6T2",
			wantErr:   false,
		},
		{
			name: "異常ケース:400 badrequest",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       ioutil.NopCloser(bytes.NewBuffer(failRespJson)),
				},
			},
			wantShUrl: "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &shorten.Bitly{}
			gotShUrl, err := b.ParseResp(tt.args.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bitly.ParseResp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotShUrl != tt.wantShUrl {
				t.Errorf("Bitly.ParseResp() = %v, want %v", gotShUrl, tt.wantShUrl)
			}
		})
	}
}
