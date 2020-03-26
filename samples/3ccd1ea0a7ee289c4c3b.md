---
title: OpenStackを入れてみた！(Solaris編)
tags: openstack
author: secondarykey
slide: false
---
### [Tokyo OpenSolaris 勉強会 2014.11 静岡合宿](http://connpass.com/event/9466/)に行ってきた！

ので久々にQiitaを書こうと思ったわけです。
本当は後編にも出たかったのですが、用事があったのでキャンセル＞＜。

OpenStack自体、漠然としたイメージだったのですが、
なんとなくわかった気がしました。

今回の勉強会（前半）はインストールが目標だったわけですが、
それがうまく行かなかったので、同じトラブルで諦めた方に
少しでもお力になれればとメモがてらに。

*Solarisでの話なので違う場合も多いって事なのでご注意を*

## OpenStackとは

どうやらIaaS（PaaS、SaaSも）環境を作り出す事ができるそうです。
※まぁこれが漠然としたイメージだったのですが。
まぁ簡単にいうとAWSが作れるよ。って事らしいです。

勉強会時のプロジェクト名は「Icehouse」、開発中の最新は「Juno」。
Solarisは大体、一個昔のプロジェクトだそうです。

Python実装だそうで少し驚きました。
RESTfulインターフェースで動作するそうです。
OpenStackサミット2015は東京で行われるらしい。

複数のコンポーネントで動作を実現しているらしく、
かなり多くのプロジェクトがありました。

* keystone
* glance
* neutron
* nova
* cinder
* swift
* horizon

辺りの設定を行います。
これらでそれぞれネットワークなり、ストレージなりを管理するそう。

## インストール

私はVirtualBoxにイメージを入れて挑みました。

* SolarisのパッケージにあるOpenStackをインストール。
* RabbitMQをインストール（AMQPプロトコルで使用）
* rad-evs-controller をインストール（evs=Elastic Virtual Switch）

これらをインストール後、
ここから各種コンポーネントの設定を行っていく感じで行いました。

設定する資料が非公開らしいので、コマンド等はかけません＞＜。すみません。

## トラブル

途中、動作確認用の
/usr/demo/openstack/keystone/sample_data.sh
でデータを投入するのですが、
何故かここで、409エラーやら500エラー（HTTPステータス）を
吐いてしまいました。

Authorization Failed: An unexpected error prevented the server from fulfilling your request. Multiple rows were found for one() (HTTP 500) 

みたいなやつ。

## 解決

* /var/lib/keystone/keystone.sqliteを削除 
* 同ファイルを空で作成（権限をkeustoneにする）
* 「keystone_manage db_sync」コマンドでデータベースと同期

で動作しました。
前述のエラーが出る前にSQLっぽいエラーで
409エラーを吐きます。多分何かが重複しているような。。。

実際シェルを２回叩くとエラーになって、
Keystoneが動作しませんでした。

・・・なんでこうなるんだろ。。。わかりません。

## 最後に

あくまで当時のバージョンなので、参考にされる時はお気をつけて。

非公開資料も多く、メモ書きも少なかったのでこの程度になりましたが、
[この辺り](http://www.oracle.com/technetwork/jp/server-storage/solaris11/technologies/openstack-2135773-ja.html)を見れば、詳しく載ってるんではないかと。

設定とかのクラウドっぽさを見たかったんですけど、
トラブルにより、ネットワーク周りの詳しい話を聞き逃しました。

もう少し触ってみて、何か出せればと思う次第です(_^_)

勉強会自体はSolaris愛に満ちた（と勝手に感じた）雰囲気でした。

