# qo is Qiita Output

かっとなってQiitaのスクレイピングを書いたら
思いのほか、初学者向けにいい感じがしたので公開しておきます
サンプルにあるのは私の記事を抜き出したところです

# Install

  $ go get github.com/secondarykey/qo/cmd/qo

## Chome Web Driver

JavaScriptを展開する為にChromeDriverを利用しています
実行場所かパスにご自身が使用されているChromeと同じバージョンのドライバをおいてください

https://chromedriver.chromium.org/downloads

## Run

GOPATHにパスが通ってたら

  $ qo {ユーザID}

でそのユーザの記事とそのマークダウンをダウンロードして来ます。

### リクエスト時間

  $ qo -r 10 {ユーザID}

とするとアクセスは10秒毎に行います
デフォルトは2ですが、用法用量を守ってお使いください。

### 並行処理

あくまで技術検証として、goroutineブランチを用意しています
リクエストの制限も取っ払ってますので、大量に記事を書いているユーザのダウンロードは行わないでください

# 一覧

一応CSVで記事一覧を作成しています。
マークダウン内にタグとかは合ったので使えるのは日付位だと思います


