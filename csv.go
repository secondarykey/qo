package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

func generateCSV(id string, items []*Item) error {

	log.Println(fmt.Sprintf("記事数:%d", len(items)))

	name := filepath.Join(id, id) + ".csv"

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
	return nil
}
