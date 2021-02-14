package tiff

import (
	"encoding/binary"
	"errors"
)

// Errors
var (
	// ErrInvalidHeader is an error for an Invalid Exif TiffHeader
	ErrInvalidHeader = errors.New("error TiffHeader is not valid")
)

// Header is the first 8 bytes of a Tiff Directory.
//
// A Header contains the byte Order, first Ifd Offset,
// tiff Header offset, Exif Length (0 if unknown) and
// Image type for the parsing of the Exif information from
// a Tiff Directory.
type Header struct {
	ByteOrder        binary.ByteOrder
	FirstIfdOffset   uint32
	TiffHeaderOffset uint32
	ExifLength       uint32
}

// NewHeader returns a new TiffHeader.
func NewHeader(byteOrder binary.ByteOrder, firstIfdOffset, tiffHeaderOffset uint32, exifLength uint32) Header {
	return Header{
		ByteOrder:        byteOrder,
		FirstIfdOffset:   firstIfdOffset,
		TiffHeaderOffset: tiffHeaderOffset,
		ExifLength:       exifLength,
	}
}

// IsValid returns true if the TiffHeader ByteOrder is not nil and
// the FirstIfdOffset is greater than 0
func (th Header) IsValid() bool {
	return th.ByteOrder != nil || th.FirstIfdOffset > 0
}