package abstract

import "context"

type IFetchCommand interface {
	GetData(ctx context.Context, url string) error
	WriteData() error
}
