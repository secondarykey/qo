package qo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
	"golang.org/x/xerrors"
)

const (
	QiitaDomain     = "https://qiita.com"
	QiitaItemNum    = "p.UserCounterList__UserCounterItemCount-sc-1xyqx6o-2"
	QiitaItemClass  = "div.AllArticleList__Item-mhtjc8-2"
	QiitaTitleClass = "a.AllArticleList__ItemBodyTitle-mhtjc8-6"
	QiitaTagClass   = "a.AllArticleList__TagListTag-mhtjc8-4"
	QiitaDateClass  = "div.AllArticleList__Timestamp-mhtjc8-8"
)

var NoLongerError = fmt.Errorf("no longer")

func getItemNum(id string) (int, error) {

	log.Println("記事数の取得を行います")
	num := -1
	url := QiitaDomain + "/" + id

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return num, xerrors.Errorf("記事数の取得に失敗しました(%s): %w", url, err)
	}

	numSels := doc.Find(QiitaItemNum)

	if numSels.Length() <= 0 {
		return num, xerrors.Errorf("記事数の要素の取得に失敗しました(%s)", url)
	}

	numSel := numSels.First()
	numBuf := numSel.Text()

	num, err = strconv.Atoi(numBuf)
	if err != nil {
		return num, xerrors.Errorf("記事数の要素に問題があります(%s): %w", numBuf, err)
	}

	log.Println(fmt.Sprintf("記事数:%d", num))

	return num, nil
}

func getItems(id string, num int) ([]*Item, error) {

	base := QiitaDomain + "/" + id
	page := 0

	items := make([]*Item, 0, num)

	for {

		page++
		url := fmt.Sprintf(base+"?page=%d", page)

		log.Println("Access:" + url)

		wkItems, err := getUserItem(url)
		if err != nil {
			if errors.Is(err, NoLongerError) {
				break
			}
			return nil, xerrors.Errorf("一覧の取得に失敗(%s):%w", url, err)
		}

		items = append(items, wkItems...)

		//ページ数で５をかけて終了
		if (page * 5) >= num {
			break
		}
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

func generateItems(op *Option, items []*Item) error {

	log.Println(fmt.Sprintf("記事のダウンロードを開始します"))

	for _, item := range items {
		err := generateItem(op, item)
		if err != nil {
			return xerrors.Errorf("generate item: %w", err)
		}

		time.Sleep(time.Duration(op.Duration) * time.Second)
	}

	return nil
}

func generateItem(op *Option, item *Item) error {

	url := QiitaDomain + item.URL + ".md"

	log.Println("Access:" + url)
	resp, err := http.Get(url)
	if err != nil {
		return xerrors.Errorf("item request[%s]: %w", url, err)
	}
	defer resp.Body.Close()

	itemId := item.ID()
	if itemId == "" {
		return xerrors.Errorf("item id error[%s]: %w", item.URL, err)
	}

	path := op.GetPath()
	name := filepath.Join(path, itemId) + ".md"

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
