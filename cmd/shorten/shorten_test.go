package shorten

import (
	"context"
	"testing"
)

func TestShortFetchComandGetData(t *testing.T) {
	type fields struct {
		urlPairs map[string]string
		Fecther  IFetchShUrl
	}
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sfc := &ShortFetchComand{
				urlPairs: tt.fields.urlPairs,
				Fecther:  tt.fields.Fecther,
			}
			if err := sfc.GetData(tt.args.ctx, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("ShortFetchComand.GetData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
