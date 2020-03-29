package main

import (
	"flag"
	"fmt"
	"log"
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
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}

	log.Println("Success")
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

	err := createDirectory(user)
	if err != nil {
		return xerrors.Errorf("ディレクトリが作成できませんでした: %w", err)
	}

	num, err := getItemNum(user)
	if err != nil {
		return xerrors.Errorf("記事数が取れませんでした: %w", err)
	}

	//記事一覧を取得
	items, err := getItems(user, num)
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

func createDirectory(dir string) error {
	if _, err := os.Stat(dir); err == nil {
		err = os.RemoveAll(dir)
		if err != nil {
			return xerrors.Errorf("ディレクトリの削除に失敗しました: %w", err)
		}
	}
	//ディレクトリ作成
	err := os.Mkdir(dir, 0777)
	if err != nil {
		return xerrors.Errorf("ディレクトリの作成に失敗しました: %w", err)
	}
	return nil
}
