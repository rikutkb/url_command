package shorten_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/rikutkb/url_command/cmd/shorten"
)

func TestTinyURLCreateRequest(t *testing.T) {
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
			b := &shorten.TinyURL{}
			_, err := b.CreateReq(tt.args.baseUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tiny.CreateReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTinyURLParseResp(t *testing.T) {
	sucRespJson := []byte(`{
		"data": {
		  "domain": "tinyurl.com",
		  "alias": "example-alias",
		  "deleted": false,
		  "archived": false,
		  "tags": [
			"tag1",
			"tag2"
		  ],
		  "analytics": [
			{
			  "enabled": true,
			  "public": false
			}
		  ],
		  "tiny_url": "https://tinyurl.com/example-alias",
		  "url": "http://google.com"
		},
		"code": 0,
		"errors": []
	  }`)
	failRespJson := []byte(`
	{
		"code": 4,
		"data": [],
		"errors": [
		  "Method not allowed"
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
			wantShUrl: "https://tinyurl.com/example-alias",
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
			b := &shorten.TinyURL{}
			gotShUrl, err := b.ParseResp(tt.args.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tiny.ParseResp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotShUrl != tt.wantShUrl {
				t.Errorf("Tiny.ParseResp() = %v, want %v", gotShUrl, tt.wantShUrl)
			}
		})
	}
}
