package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/secondarykey/qo"
	"golang.org/x/xerrors"
)

var dur *int

func init() {
	dur = flag.Int("r", 2, "request duration")
}

func Usage() {
}

func main() {

	err := run()
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
	log.Println("Success")
}

func run() error {

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

	err := qo.Run(args[0], qo.Duration(*dur), Remove(true))
	if err != nil {
		return xerrors.Errorf("コマンド実行中にエラーが発生しました: %w", err)
	}

	return nil
}

func Remove(f bool) func(*qo.Option) {
	return func(op *qo.Option) {
		op.ExistAndRemove = f
	}
}

func Path(d string) func(*qo.Option) {
	return func(op *qo.Option) {
		op.Path = d
	}
}
