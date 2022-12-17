package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/rikutkb/url_command.git/cmd/undo"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "短縮化されたurlを元に戻します。",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		splitedUrl := strings.Split(urls, ",")
		var wg sync.WaitGroup
		sem := make(chan bool, 3)
		for _, url := range splitedUrl {
			wg.Add(1)

			sem <- true
			go func(_url string) {

				defer wg.Done()
				longUrl, err := undo.GetRedirect(_url)
				<-sem
				if err != nil {
					fmt.Fprintf(os.Stderr, err.Error())
				}
				fmt.Println(longUrl)

			}(url)
		}
		defer close(sem)
		wg.Wait()
	},
}

func init() {
	undoCmd.PersistentFlags().StringVarP(&urls, "url", "u", "", "短縮API")

	rootCmd.AddCommand(undoCmd)
}
