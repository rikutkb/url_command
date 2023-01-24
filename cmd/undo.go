package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rikutkb/url_command.git/cmd/undo"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "短縮化されたurlを元に戻します。",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		splitedUrls := strings.Split(urls, ",")
		sfc := &undo.UndoFetchCommand{}
		err := resolveUrls(ctx, splitedUrls, sfc)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	},
}

func init() {
	undoCmd.PersistentFlags().StringVarP(&urls, "url", "u", "", "短縮API")

	rootCmd.AddCommand(undoCmd)
}
