# HOW TO USE
## option
### urlを短くして出力

url [-s][--short] url

### csvなどのファイルから読み込み(標準出力)
url [-f] file

### ファイル内url置き換え
url [-rep][--replace] file

### 使用するサービスの指定
url [-k][--kind] (bit|tinyurl)

### 元url取得
url [-u][--undo] url


### APIキーの登録
bitlyの使用時

export BIT_API_KEY=xxxxxxxxxxxxxxxxxxxxxxx 


tinyURLの使用時

$ export TINYURL_API_KEY=xxxxxxxxxxxxxxxxxxxxxxx 
$ ~ % url_command shorten -s bitly -u https://github.com/rikutkb/url_command,https://www.youtube.com/feed/subscriptions
https://bit.ly/3R88Drn,https://bit.ly/3RsoaTh                 
APIキーがセットされていません。
$ ~ % url_command shorten -s TinyURL -u https://github.com/rikutkb/url_command,https://www.youtube.com/feed/subscriptions
https://tinyurl.com/2kvjk6jm,https://tinyurl.com/bdvnndt9

$~ % url_command undo -u https://tinyurl.com/2kvjk6jm,https://tinyurl.com/bdvnndt9
https://github.com/rikutkb/url_command,https://www.youtube.com/feed/subscriptions