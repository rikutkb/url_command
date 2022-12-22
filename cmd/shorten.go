/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/rikutkb/url_command.git/cmd/shorten"
	"github.com/spf13/cobra"
)

func shortenUrls(ctx context.Context, reqUrls []string, fetcher shorten.IFetchShUrl) error {
	var wg sync.WaitGroup
	sem := make(chan bool, 3)
	defer close(sem)

	for _, url := range reqUrls {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		wg.Add(1)

		sem <- true
		go func(_url string) {
			defer wg.Done()
			shortUrl, err := shorten.CreateShortUrl(ctx, _url, fetcher)
			if err != nil {
				SetCancelSignal(err.Error())
				fmt.Fprintln(os.Stderr, err)
				return
			}
			fmt.Fprintln(os.Stdout, "baseUrl:"+_url+", resultUrl:"+shortUrl)
			<-sem

		}(url)
	}
	wg.Wait()
	return nil
}

var shortenCmd = &cobra.Command{
	Use:   "shorten",
	Short: "APIを使用しurlの短縮化を行います。",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		var fetcher = shorten.NewFecher(service)
		if err := fetcher.Init(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		splitedUrls := strings.Split(urls, ",")
		err := shortenUrls(ctx, splitedUrls, fetcher)
		if err != nil {

			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	},
}

func init() {
	shortenCmd.PersistentFlags().StringVarP(&urls, "url", "u", "", "短縮API")
	shortenCmd.PersistentFlags().StringVarP(&service, "service", "s", "", "短縮API")

	//shortenCmd.PersistentFlags().StringArrayVarP(&urls, "url", "u", make([]string, 0), "短縮API")
	rootCmd.AddCommand(shortenCmd)
}
