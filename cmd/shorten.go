/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rikutkb/url_command.git/cmd/shorten"
	"github.com/spf13/cobra"
)

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
		splitedUrl := strings.Split(urls, ",")

		shortUrl, err := shorten.CreateShortUrl(splitedUrl[0], fetcher)
		if err != nil {

			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Println(shortUrl)
	},
}

func init() {
	shortenCmd.PersistentFlags().StringVarP(&urls, "url", "u", "", "短縮API")
	//shortenCmd.PersistentFlags().StringArrayVarP(&urls, "url", "u", make([]string, 0), "短縮API")
	rootCmd.AddCommand(shortenCmd)
}
