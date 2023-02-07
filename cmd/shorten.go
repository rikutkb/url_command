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

	"github.com/rikutkb/url_command/cmd/abstract"
	"github.com/rikutkb/url_command/cmd/shorten"
	"github.com/spf13/cobra"
)

func resolveUrls(ctx context.Context, reqUrls []string, cmd abstract.IFetchCommand) error {
	var wg sync.WaitGroup
	sem := make(chan bool, 3)
	defer close(sem)
	_ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, url := range reqUrls {
		done := make(chan error, 1)
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		wg.Add(1)
		sem <- true
		go func(_url string) {
			defer wg.Done()
			err := cmd.GetData(_ctx, _url)
			done <- err
			if err != nil {
				return
			}
			<-sem
		}(url)
		if err := <-done; err != nil {
			return err
		}
	}
	wg.Wait()

	if err := cmd.WriteData(reqUrls); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

var shortenCmd = &cobra.Command{
	Use:   "shorten",
	Short: "APIを使用しurlの短縮化を行います。",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		switch service {
		case "TinyURL", "bitly":
		default:
			fmt.Fprintln(os.Stderr, "指定されたサービスは使用できません。使用できるサービスはbitly or TinyURLです。")
			os.Exit(2)
		}
		var fetcher = shorten.NewFecher(service)
		if err := fetcher.Init(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		splitedUrls := strings.Split(urls, ",")
		sfc := shorten.NewShortFetchCommand(fetcher)
		err := resolveUrls(ctx, splitedUrls, sfc)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	},
}

func init() {
	shortenCmd.PersistentFlags().StringVarP(&urls, "url", "u", "", "短縮API")
	shortenCmd.PersistentFlags().StringVarP(&service, "service", "s", "", "短縮API")
	rootCmd.AddCommand(shortenCmd)
}
