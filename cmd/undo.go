package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rikutkb/url_command.git/cmd/undo"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func fetchUrls(ctx context.Context, reqUrls []string) {
	sem := make(chan bool, 2)
	eg, ctx := errgroup.WithContext(ctx)
	defer close(sem)

	for _, url := range reqUrls {
		select {
		case <-ctx.Done():
			return
		default:
		}

		sem <- true
		func(_url string) {
			eg.Go(func() error {
				defer func() {
					<-sem
				}()
				longUrl, err := undo.GetRedirect(ctx, _url)
				if err != nil {
					SetCancelSignal(err.Error())
					fmt.Fprintln(os.Stderr, err)
					return err
				}
				fmt.Fprintln(os.Stdout, "baseUrl:"+_url+", resultUrl:"+longUrl)
				return nil
			})
		}(url)

	}
	if err := eg.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "短縮化されたurlを元に戻します。",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		splitedUrl := strings.Split(urls, ",")

		fetchUrls(ctx, splitedUrl)
	},
}

func init() {
	undoCmd.PersistentFlags().StringVarP(&urls, "url", "u", "", "短縮API")

	rootCmd.AddCommand(undoCmd)
}
