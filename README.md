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
url [-s][--service] (bit|pixiv)

### 元url取得
url [-u][--undo] url

### apiトークン設定
url [-i]



### APIキーの登録
bitlyの使用時
export BIT_API_KEY=xxxxxxxxxxxxxxxxxxxxxxx
