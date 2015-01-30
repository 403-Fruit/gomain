package main

import (
	"github.com/mattn/go-colorable"
	"github.com/tj/go-spin"
)

func init() {
	out = colorable.NewColorableStdout()
	box = spin.Spin1
}
