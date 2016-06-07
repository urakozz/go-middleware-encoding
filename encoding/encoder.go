// Copyright 2016 Yury Kozyrev

package encoding

import "io"

type Encoder interface {
	NewWriter(w io.Writer, level int) io.Writer
	GetName() string
}

type EncoderNone struct {}
func (_ *EncoderNone) GetName() (s string){ return }
func (_ *EncoderNone) NewWriter(w io.Writer, level int) io.Writer {
	return w
}

var _ Encoder = (*EncoderNone)(nil)

