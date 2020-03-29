package main

import (
	"strings"
)

type Item struct {
	Date string
	Name string
	Tags []string
	URL  string
}

func NewItem() *Item {
	i := Item{}
	i.Tags = make([]string, 0)
	return &i
}

func (i *Item) AddTag(tag string) {
	i.Tags = append(i.Tags, tag)
	return
}

func (i *Item) ID() string {
	idSlc := strings.Split(i.URL, "/")
	if len(idSlc) <= 3 {
		return ""
	}
	return idSlc[3]
}

func (i *Item) Slice() []string {
	line := make([]string, 4)
	line[0] = i.Date
	line[1] = i.Name
	buf := ""
	for idx, tag := range i.Tags {
		if idx != 0 {
			buf += ";"
		}
		buf += tag
	}
	line[2] = buf

	line[3] = i.URL
	return line
}
