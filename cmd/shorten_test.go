/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rikutkb/url_command/cmd/abstract"
)

var _ abstract.IFetchCommand = &TestFetchComand{}

type TestFetchComand struct {
	urlPairs map[string]string
}

func (tc *TestFetchComand) GetData(ctx context.Context, url string) error {

	tc.urlPairs[url] = url
	if url == "test3" {
		time.Sleep(3 * time.Second)
	} else {
		time.Sleep(1 * time.Microsecond)
	}
	return nil
}
func (tc TestFetchComand) WriteData(reqUrls []string) error {
	for i, url := range reqUrls {
		fmt.Fprintf(os.Stdout, tc.urlPairs[url])
		if i+1 != len(reqUrls) {
			fmt.Fprintf(os.Stdout, ",")
		} else {
			fmt.Fprintln(os.Stdout, "")
		}
	}
	return nil
}

func TestResolveUrls(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqUrls []string
		cmd     abstract.IFetchCommand
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "リクエスト成功",
			args: args{
				ctx:     context.Background(),
				reqUrls: []string{"test1", "test2", "test3"},
				cmd:     &TestFetchComand{urlPairs: map[string]string{}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := resolveUrls(tt.args.ctx, tt.args.reqUrls, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("resolveUrls() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
