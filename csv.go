package qo

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

func generateCSV(op *Option, items []*Item) error {

	id := op.id

	path := op.GetPath()
	name := filepath.Join(path, id) + ".csv"

	f, err := os.Create(name)
	if err != nil {
		return xerrors.Errorf("get user item error: %w", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)

	for _, item := range items {
		err = writer.Write(item.Slice())
		if err != nil {
			return xerrors.Errorf("csv write error: %w", err)
		}
	}

	writer.Flush()

	log.Println(fmt.Sprintf("%s に一覧を出力しました", name))

	return nil
}
