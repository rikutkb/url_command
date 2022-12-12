package main

import "github.com/rikutkb/url_command.git/cmd"

// var (
// 	url_flag     string
// 	short_flag   bool
// 	service_flag string
// 	replace_flag bool
// 	undo_flag    bool
// 	set_flag     string
// 	qr_flag      bool
// 	file_flag    string
// 	init_flag    string
// 	kind_flag    string
// 	token_flag   string
// )

// func FlagInit() {
// 	flag.BoolVar(&short_flag, "s", false, "urlの短縮を行います。")
// 	flag.BoolVar(&replace_flag, "r", false, "urlの置換を行います。")
// 	flag.BoolVar(&undo_flag, "u", false, "短縮urlを元に戻します。")
// 	flag.BoolVar(&qr_flag, "q", false, "短縮urlを元に戻します。")

// 	flag.StringVar(&url_flag, "url", "", "url")
// 	flag.StringVar(&service_flag, "service", "bitly", "urlの短縮を行います。")
// 	flag.StringVar(&set_flag, "set", "", "")
// 	flag.StringVar(&file_flag, "f", "", "")
// 	flag.StringVar(&kind_flag, "k", "bitly", "使用サービスの指定を行います。")
// 	flag.StringVar(&token_flag, "t", "", "APIトークンの設定を行います。")
// }

func main() {
	// FlagInit()
	// if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
	// 	if err != flag.ErrHelp {
	// 		fmt.Fprintf(os.Stderr, "error: %s", err)
	// 	}
	// 	os.Exit(2)
	// }
	// // 短縮orQRコード,短縮urlを元に戻す
	// if url_flag != "" {
	// 	// TODO APIキーが設定できていない場合はエラーとして出力するようにする。
	// 	// TODO interfaceとしてhttp通信部分を実装
	// 	url := url_flag

	// 	var fetcher = shorten.NewFecher(kind_flag)
	// 	if err := fetcher.Init(); err != nil {
	// 		fmt.Fprintln(os.Stderr, err)
	// 		os.Exit(2)
	// 	}
	// 	shortUrl, err := shorten.CreateShortUrl(url, fetcher)
	// 	if err != nil {

	// 		fmt.Fprintln(os.Stderr, err)
	// 		os.Exit(2)
	// 	}
	// 	fmt.Println(shortUrl)
	// } else {
	// 	fmt.Errorf("不正なコマンド入力です。")
	// 	os.Exit(2)
	// }
	cmd.Execute()

}
