package qo

import (
	"path/filepath"
)

type Option struct {
	id             string
	Path           string
	ExistAndRemove bool
	Duration       int
}

func (op *Option) GetPath() string {
	return filepath.Join(op.Path, op.id)
}

func NewOption(id string) *Option {
	op := Option{}
	op.id = id
	op.Path = ""
	op.ExistAndRemove = false
	op.Duration = 1
	return &op
}

func Duration(d int) func(*Option) {
	return func(option *Option) {
		option.Duration = d
	}
}
