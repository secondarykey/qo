---
title: Dockerを触ってみた
tags: Docker
author: secondarykey
slide: false
---
ちょっとApache+PHP+MySQLとかの環境を作る機会があって、いつも自分の端末に直接インストールしていたんだけど、流行りのDockerってどんなもんだろ？と思って触ってみた。

インストールした端末はUbuntu14.04です。
# インストール
多くの文献等にVirtualBox使った例があったけど、


```
curl -s https://get.docker.io/ubuntu/ | sudo sh
```

で行った。
どうやって動いているとかもう少し勉強が必要だなと
感じたけど、またそれは今度。まずは構築。

```
docker -v
```

これでバージョンが出たらインストール完了。
自分のバージョンは「0.11.1」でした。

# イメージをインストール
私の作る環境はCentOSだったので

```
sudo docker search centos
```
と行って、検索した。
ユーザ登録してリポジトリに登録するーとかがあるらしい、
なんだかなうい。

```
sudo docker pull centos:6.4
```

これで持ってくる事ができた。
イメージ名は「[name]/[os?]:[tag]」とするのが慣習っぽかったけど、これはcentosのタグが6.4となっていた。

#確認と起動

```
sudo docker images
```

とするとpullしてきたImageが確認できる

```
sudo docker run -i -t centos:6.4 /bin/bash
```
と行う事で、イメージからコンテナが起動して後ろのコマンド（bash）が起動する

このイメージは一番軽いイメージっぽいのでほとんど何も入ってない
コンテナを起動すると

```
sudo docker ps
```
これでコンテナの一覧を確認できた。
コンテナIDとかはこれで確認するみたい。

起動しているコンテナで

```
 # exit
```

とするとコンテナが終了する。
操作したコマンドやファイル等は保存されてない。

#コンテナの保存

起動したイメージはどんなにいじっても保存されない。
保存するには

```
sudo docker commit <コンテナID> centos:6.4
```

として既存のイメージが新規のイメージに永続化を行う。
逆にいうと保存しなければ残らないから何度でも試す事ができる

>あとで気付いたんだけど、
>大体開発環境の構築って間違った事を行ってしまった場合、関連性とかがぐちゃぐちゃになりながらやっとこ動き出すみたいな事が多くて、Dockerでこれをやる場合、「あっ間違った！このコンテナ捨てよう！」っていう事をやらないと勿体ないと気付いた。

# コンテナ、イメージの削除

試している間、commitの使い方とかがわからなくて
無駄にコミットしたりとかでイメージやコンテナがぽこぽこできてた。

```
sudo docker rm <コンテナID>
```
で不要になったコンテナが削除できる。
-fを使うと強制的に削除する事ができる。
ゾンビみたいな奴はこれで殺せる。

イメージは

```
sudo docker rmi <イメージIDかイメージ名>
```

で削除できる。

#感想
やり終わって感じた事は、例えば大きなチーム開発の場合、
新しいメンバーにDockerのイメージを渡す〜とかが出てくるのかな？とか感じた。
ファイルのやりとりとかが結構難しくて、（私はDockerfile作って移動してみた）ソース管理とかどうするかなーと感じているので後日その辺りを触ったらまた書こうと思う。

