package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rikutkb/url_command.git/cmd/undo"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func fetchUrls(ctx context.Context, reqUrls []string) (error, map[string]string) {
	sem := make(chan bool, 3)
	eg, ctx := errgroup.WithContext(ctx)
	defer close(sem)
	var urlPairs = make(map[string]string)

	for _, url := range reqUrls {
		select {
		case <-ctx.Done():
			return errors.New("fetch is canceled."), nil
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
				urlPairs[_url] = longUrl
				return nil
			})
		}(url)

	}
	if err := eg.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err, nil
	}
	return nil, urlPairs
}

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "短縮化されたurlを元に戻します。",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		splitedUrl := strings.Split(urls, ",")
		if err, pairUrls := fetchUrls(ctx, splitedUrl); err != nil {
			os.Exit(2)
		} else {
			for baseUrl, resultUrl := range pairUrls {
				fmt.Fprintln(os.Stdout, "baseUrl:"+baseUrl+",resultUrl: "+resultUrl)
			}
		}
	},
}

func init() {
	undoCmd.PersistentFlags().StringVarP(&urls, "url", "u", "", "短縮API")

	rootCmd.AddCommand(undoCmd)
}
