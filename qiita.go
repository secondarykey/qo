package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
	"golang.org/x/xerrors"
)

const (
	QiitaDomain     = "https://qiita.com"
	QiitaItemClass  = "div.AllArticleList__Item-mhtjc8-2"
	QiitaTitleClass = "a.AllArticleList__ItemBodyTitle-mhtjc8-6"
	QiitaTagClass   = "a.AllArticleList__TagListTag-mhtjc8-4"
	QiitaDateClass  = "div.AllArticleList__Timestamp-mhtjc8-8"
)

var NoLongerError = fmt.Errorf("no longer")

func getItems(id string) ([]*Item, error) {

	page := 1
	items := make([]*Item, 0, 100)
	base := QiitaDomain + "/" + id
	url := base

	for {

		log.Println("Access:" + url)

		wkItems, err := getUserItem(url)
		if err != nil {
			if errors.Is(err, NoLongerError) {
				log.Println("記事の一覧を抽出しました。")
				break
			}
			return nil, xerrors.Errorf("get user item error: %w", err)
		}

		items = append(items, wkItems...)
		page++
		url = fmt.Sprintf(base+"?page=%d", page)
	}

	return items, nil
}

func getUserItem(url string) ([]*Item, error) {

	html, err := getHTML(url)
	if err != nil {
		return nil, xerrors.Errorf("記事一覧の取得に失敗しました=[%s]: %w", url, err)
	}

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, xerrors.Errorf("記事一覧の取得に失敗しました=[%s]: %w", url, err)
	}

	items := make([]*Item, 0, 5)

	sel := doc.Find(QiitaItemClass)
	if sel.Length() <= 0 {
		return nil, NoLongerError
	}

	sel.Each(func(_ int, s *goquery.Selection) {

		item := NewItem()

		title := s.Find(QiitaTitleClass)
		title.Each(func(_ int, a *goquery.Selection) {
			val, ok := a.Attr("href")
			if ok {
				item.URL = val
				item.Name = a.Text()
			}
		})

		tag := s.Find(QiitaTagClass)
		tag.Each(func(_ int, a *goquery.Selection) {
			item.AddTag(a.Text())
		})

		date := s.Find(QiitaDateClass)
		date.Each(func(_ int, a *goquery.Selection) {
			item.Date = a.Text()
		})

		items = append(items, item)
	})

	return items, nil
}

func getHTML(url string) (io.Reader, error) {
	options := agouti.ChromeOptions(
		"args", []string{
			"--headless",
			"--disable-gpu",
			"--no-sandbox",
		})

	driver := agouti.ChromeDriver(options)
	defer driver.Stop()
	driver.Start()

	page, err := driver.NewPage()
	if err != nil {
		return nil, xerrors.Errorf("get user item error: %w", err)
	}

	err = page.Navigate(url)
	if err != nil {
		return nil, xerrors.Errorf("get user item error[%s]: %w", url, err)
	}
	html, err := page.HTML()
	if err != nil {
		return nil, xerrors.Errorf("get html: %w", err)
	}

	reader := bytes.NewBufferString(html)
	return reader, nil
}

func generateItems(id string, items []*Item) error {

	log.Println(fmt.Sprintf("記事のダウンロードを開始します(%d秒)", *dur))

	for _, item := range items {
		err := generateItem(id, item)
		if err != nil {
			return xerrors.Errorf("generate item: %w", err)
		}

		time.Sleep(time.Duration(*dur) * time.Second)
	}
	return nil
}

func generateItem(id string, item *Item) error {

	url := QiitaDomain + item.URL + ".md"

	log.Println("Access:" + url)
	resp, err := http.Get(url)
	if err != nil {
		return xerrors.Errorf("item request[%s]: %w", url, err)
	}
	defer resp.Body.Close()

	idSlc := strings.Split(item.URL, "/")
	itemId := idSlc[3]

	name := filepath.Join(id, itemId) + ".md"

	f, err := os.Create(name)
	if err != nil {
		return xerrors.Errorf("create item file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return xerrors.Errorf("response copy : %w", err)
	}

	log.Println(fmt.Sprintf("    -> %s", name))

	return nil
}
