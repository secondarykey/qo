package qo

import (
	"fmt"
	"os"

	"golang.org/x/xerrors"
)

func Run(id string, options ...func(*Option)) error {

	op := NewOption(id)
	for _, option := range options {
		option(op)
	}

	err := createDirectory(op)
	if err != nil {
		return xerrors.Errorf("ディレクトリが作成できませんでした: %w", err)
	}

	num, err := getItemNum(op.id)
	if err != nil {
		return xerrors.Errorf("記事数が取れませんでした: %w", err)
	}

	//記事一覧を取得
	items, err := getItems(op.id, num)
	if err != nil {
		return xerrors.Errorf("記事の一覧が作成できませんでした: %w", err)
	}

	//CSVを作成
	err = generateCSV(op, items)
	if err != nil {
		return xerrors.Errorf("一覧CSVの作成に失敗しました: %w", err)
	}

	//記事の一覧をダウンロード
	err = generateItems(op, items)
	if err != nil {
		return xerrors.Errorf("記事データの作成に失敗しました: %w", err)
	}

	return nil
}

func createDirectory(op *Option) error {

	dir := op.GetPath()

	if _, err := os.Stat(dir); err == nil {
		if op.ExistAndRemove {
			err = os.RemoveAll(dir)
			if err != nil {
				return xerrors.Errorf("ディレクトリの削除に失敗しました: %w", err)
			}
		} else {
			return fmt.Errorf("ディレクトリが存在します(%s)", dir)
		}
	}
	//ディレクトリ作成
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return xerrors.Errorf("ディレクトリの作成に失敗しました: %w", err)
	}
	return nil
}
