package qrcode

import (
	"context"
	"net/http"

	"github.com/rikutkb/url_command/cmd/abstract"
)

var _ abstract.IFetchCommand = &QRFetchCommand{}

type QRFetchCommand struct {
	urlPairs map[string][]byte
}

func NewQRFetchCommand() *QRFetchCommand {
	return &QRFetchCommand{urlPairs: map[string][]byte{}}
}

func (qrf *QRFetchCommand) GetData(ctx context.Context, url string) error {
	return nil
}

func (sfc QRFetchCommand) WriteData(reqUrls []string) error {

	return nil
}

type IFetchQRUrl interface {
	Init() (err error)
	CreateRequest(url string) (err error, request *http.Request)
	ParseResponse(response *http.Response) (err error, image []byte)
}
