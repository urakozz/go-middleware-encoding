// Copyright 2016 Yury Kozyrev

package encoding

import (
	"github.com/klauspost/compress/flate"
	"io"
	"github.com/klauspost/compress/gzip"
)

type EncoderDeflate struct {}
func (_ *EncoderDeflate) GetName() (string){ return "deflate" }
func (_ *EncoderDeflate) NewWriter(w io.Writer, level int) (ew io.Writer){
	if level < 1 || level > 11{
		level = 1
	}
	ew, _ = flate.NewWriter(w, level)
	return
}

var _ Encoder = (*EncoderDeflate)(nil)

type EncoderGzip struct {}
func (_ *EncoderGzip) GetName() (string){ return "gzip" }
func (_ *EncoderGzip) NewWriter(w io.Writer, level int) (ew io.Writer){
	if level < 1 || level > 11{
		level = 1
	}
	ew, _ = gzip.NewWriterLevel(w, level)
	return
}

var _ Encoder = (*EncoderGzip)(nil)

