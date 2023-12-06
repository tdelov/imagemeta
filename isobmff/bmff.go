package isobmff

import (
	"io"

	"github.com/tdelov/imagemeta/meta"
)

type ExifReader func(r io.Reader, h meta.ExifHeader) error

const (
	optionSpeed uint8 = 1
)
