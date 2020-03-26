package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/xerrors"
)

var dur *int

func init() {
	dur = flag.Int("r", 2, "request duration")
}

func Usage() {

}

func main() {

	err := Run()
	if err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}

	fmt.Println("Success")
}

func Run() error {

	//ユーザID取得
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return fmt.Errorf("引数にユーザIDが必要です")
	}

	if *dur <= 0 {
		return fmt.Errorf("リクエスト区間が短いです")
	}

	if *dur > 30 {
		return fmt.Errorf("リクエスト区間が長すぎやしませんか？")
	}

	user := args[0]

	//ディレクトリ作成
	err := os.Mkdir(user, 0777)
	if err != nil {
		return xerrors.Errorf("ディレクトリが作成できませんでした: %w", err)
	}

	//記事一覧を取得
	items, err := getItems(user)
	if err != nil {
		return xerrors.Errorf("記事の一覧が作成できませんでした: %w", err)
	}

	//CSVを作成
	err = generateCSV(user, items)
	if err != nil {
		return xerrors.Errorf("一覧CSVの作成に失敗しました: %w", err)
	}

	//記事の一覧をダウンロード
	err = generateItems(user, items)
	if err != nil {
		return xerrors.Errorf("記事データの作成に失敗しました: %w", err)
	}

	return nil
}
