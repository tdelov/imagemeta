package imagemeta

import (
	"io"
	"os"

	"github.com/tdelov/imagemeta/exif2"
	"github.com/tdelov/imagemeta/isobmff"
	"github.com/tdelov/imagemeta/jpeg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is the logger
	logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel)
)

func SetLogger(w io.Writer, level zerolog.Level) {
	logger = log.Output(w).Level(level)
	jpeg.Logger = logger
	exif2.Logger = logger
	isobmff.Logger = logger
}
