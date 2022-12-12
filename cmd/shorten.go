/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rikutkb/url_command.git/cmd/shorten"
	"github.com/spf13/cobra"
)

var service string
var url string
var shortenCmd = &cobra.Command{
	Use:   "shorten",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		var fetcher = shorten.NewFecher(service)
		if err := fetcher.Init(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		shortUrl, err := shorten.CreateShortUrl(url, fetcher)
		if err != nil {

			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Println(shortUrl)
	},
}

func init() {
	shortenCmd.PersistentFlags().StringVarP(&service, "service", "s", "bitly", "使用APIサービス")
	shortenCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "短縮API")

	rootCmd.AddCommand(shortenCmd)
}
