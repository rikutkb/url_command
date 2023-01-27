package undo

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestUndoFetchCommand_GetData(t *testing.T) {

	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name        string
		args        args
		fileName    string
		baseUrl     string
		expectedUrl string
		status      int
		wantErr     bool
	}{
		{
			name: "TinyURLの正常処理",
			args: args{
				ctx: context.Background(),
				url: "tinyurl.com/2kvjk6jm",
			},
			fileName:    `/tiny-response.txt`,
			expectedUrl: "https://github.com/rikutkb/url_command",
			status:      http.StatusMovedPermanently,
			wantErr:     false,
		}, {
			name: "Bitlyの正常処理",
			args: args{
				ctx: context.Background(),
				url: "bit.ly/3Hthhhn",
			},
			fileName:    `/bitly-response.txt`,
			expectedUrl: "https://cobra.dev/",
			status:      http.StatusMovedPermanently,
			wantErr:     false,
		},
	}
	dir, _ := os.Getwd()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := dir + tt.fileName
			response, _ := ioutil.ReadFile(filePath)
			ufc := NewUndoFetchCommand()
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", tt.args.url,
				httpmock.NewStringResponder(tt.status, string(response)),
			)

			if err := ufc.GetData(tt.args.ctx, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("UndoFetchCommand.GetData() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expectedUrl, ufc.urlPairs[tt.args.url])
		})
	}
}
