// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package cencoding

import (
	"io"
	"github.com/youtube/vitess/go/cgzip"
	"github.com/urakozz/go-middleware-encoding/encoding"
)

type EncoderCGzip struct {}
func (_ *EncoderCGzip) GetName() (string){ return "gzip" }
func (_ *EncoderCGzip) NewWriter(w io.Writer, level int) (ew io.Writer){
	if level < 1 || level > 11{
		level = 1
	}
	ew, _ = cgzip.NewWriterLevel(w, level)
	return
}


var _ encoding.Encoder = (*EncoderCGzip)(nil)
