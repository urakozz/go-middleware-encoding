// Copyright 2016 Yury Kozyrev

package cencoding

import (
	"io"
	"gopkg.in/kothar/brotli-go.v0/enc"
	"github.com/urakozz/go-middleware-encoding"
)

type EncoderCBrotli struct {}
func (_ *EncoderCBrotli) GetName() (string){ return "br" }
func (_ *EncoderCBrotli) NewWriter(w io.Writer, level int) (ew io.Writer){
	p := enc.NewBrotliParams()
	if level > 0 && level <= 11{
		p.SetQuality(level)
	}
	return enc.NewBrotliWriter(p, w)
}

var _ encoding.Encoder = (*EncoderCBrotli)(nil)


