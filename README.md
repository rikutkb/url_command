# HOW TO USE
## option
### urlを短くして出力

url [-s][--short] url

### qrコードとしてurlを保存
url [-q]|[-qr] url　[-o] url.png

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

export TINYURL_API_KEY=xxxxxxxxxxxxxxxxxxxxxxx 
