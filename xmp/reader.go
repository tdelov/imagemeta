package xmp

import (
	"bufio"
	"bytes"
	"io"

	"github.com/pkg/errors"

	"github.com/tdelov/imagemeta/xmp/xmpns"
)

const (
	xmpBufferLength = 1538 // (1.5kb)

	// Reader blocksizes
	maxTagValueSize  = 512
	maxTagHeaderSize = 128
)

var (
	// Reader errors
	ErrNoValue      = errors.New("error property has no value")
	ErrNegativeRead = errors.New("error negative read")
	ErrBufferFull   = bufio.ErrBufferFull

	// xmpRootTag starts with "<x:xmpmeta" and ends with "</x:xmpmeta>"
	xmpRootTag      = [...]byte{'<', 'x', ':', 'x', 'm', 'p', 'm', 'e', 't', 'a'}
	xmpRootCloseTag = [...]byte{'<', '/', 'x', ':', 'x', 'm', 'p', 'm', 'e', 't', 'a', '>'}
)

type xmpReader struct {
	r *bufio.Reader
	a bool
}

func newXMPReader(r io.Reader) xmpReader {
	br, ok := r.(*bufio.Reader)
	if !ok || br.Size() < xmpBufferLength {
		br = bufio.NewReaderSize(r, xmpBufferLength)
	}
	return xmpReader{r: br}
}

// readRootTag reads and returns the xmpRootTag from the bufReader.
// If the xmpRootTag is not found returns the error ErrNoXMP.
func (br *xmpReader) readRootTag() (tag Tag, err error) {
	var buf []byte
	for {
		if _, err = br.r.ReadSlice(xmpRootTag[0]); err != nil {
			if err == io.EOF {
				err = ErrNoXMP
				return
			}
			if err == bufio.ErrBufferFull {
				continue
			}
		}
		if buf, err = br.r.Peek(10); err != nil {
			return
		}

		if bytes.Equal(xmpRootTag[1:], buf[0:9]) {
			_, err = br.r.ReadSlice('>') // Read until end of the StartTag (RootTag)
			tag.t = startTag
			tag.self = xmpns.XMPRootProperty
			//fmt.Println("XMP Discarded:", discarded)
			return tag, err
		}
	}
}

func (br *xmpReader) Discard(n int) (discarded int, err error) {
	return br.r.Discard(n)
}

// hasAttribute returns true when the bufReader's next read is
// an attribute.
func (br *xmpReader) hasAttribute() bool {
	return br.a
}

func (br *xmpReader) Peek(n int) (buf []byte, err error) {
	if buf, err = br.r.Peek(n); err == io.EOF {
		if len(buf) > 4 {
			return buf, nil
		}
		return buf, err
	}
	return
}

// readAttribute reads an attribute from the bufReader and Tag.
func (br *xmpReader) readAttribute(tag *Tag) (attr Attribute, err error) {
	var buf []byte
	attr.pt = attrPType
	attr.parent = tag.self

	// Attribute Name
	if buf, err = br.Peek(maxTagHeaderSize); err != nil {
		err = errors.Wrap(err, "Attr")
		return
	}

	var d int
	if attr.self, d, err = parseAttrName(buf); err != nil {
		err = errors.Wrap(ErrNegativeRead, "Attr (name)")
		return
	}
	if _, err = br.Discard(d); err != nil {
		err = errors.Wrap(err, "Attr (discard)")
		return
	}

	// Attribute Value
	attr.val, err = br.readAttrValue(tag)

	return attr, err
}

// readAttrValue reada an Attributes value from the Tag.
// Needs improvement for performance
func (br *xmpReader) readAttrValue(tag *Tag) (buf []byte, err error) {
	d, i := 0, 2
	s := maxTagValueSize / 2
	for {
		if buf, err = br.Peek(s); err != nil {
			err = errors.Wrap(err, "Attr Value")
			return
		}

		if buf[0] == '=' && (buf[1] == '"' || buf[1] == '\'') {
			delim := buf[1]
			if b := bytes.IndexByte(buf[i:], delim); b >= 0 {
				i += b
				d = i + 1
				if buf[i+1] == '>' {
					d++
					br.a = false
				} else if buf[i+1] == '/' && buf[i+2] == '>' {
					d += 2
					tag.t = soloTag
					br.a = false
				}
				if _, err = br.Discard(d); err != nil {
					err = errors.Wrap(err, "Attr Value (discard)")
				}
				return buf[2:i], err
			}
		}
		s += maxTagValueSize
	}
}

// readTagHeader reads an xmp tag's header and returns the tag.
func (br *xmpReader) readTagHeader(parent Tag) (tag Tag, err error) {
	tag.pt = tagPType
	tag.parent = parent.self

	s := maxTagHeaderSize
	// Read Tag Header
	var buf []byte
	var i int
	for {
		if buf, err = br.Peek(s); err != nil {
			err = errors.Wrap(err, "Tag Header")
			return
		}

		// Find Start of Tag
		for ; i < len(buf); i++ {
			if buf[i] == '<' {
				if buf[i+1] == '/' {
					tag.t = stopTag
					i += 2
				} else if buf[i+1] == '?' {
					err = io.EOF
					return
				} else {
					tag.t = startTag
					i++
				}
				buf = buf[i:]
				goto end
			}
		}
		// large white spaces in xmp files

		s += maxTagHeaderSize
	}
end:
	var d int
	tag.self, d, err = parseTagName(buf)
	if err != nil {
		err = errors.Wrap(err, "Tag Header (tag name)") // Err finding tag name
		return
	}
	if buf[d] == '>' {
		br.a = false // No Attributes
		d++
	} else if buf[d] == ' ' || buf[d] == '\n' { // Attributes
		br.a = true
	} else if buf[d] == '/' && buf[d+1] == '>' { // SoloTag
		br.a = false // No Attributes
		tag.t = soloTag
		d += 2
	}
	if _, err = br.Discard(d + i); err != nil {
		err = errors.Wrap(err, "Tag Header (discard)")
	}
	return
}

// readTagValue reads the Tag's Value from the bufReader. Returns
// a temporary []byte.
func (br *xmpReader) readTagValue() (buf []byte, err error) {
	var i, j int
	s := maxTagValueSize
	for {
		if buf, err = br.Peek(s); err != nil {
			err = errors.Wrap(err, "Tag Value")
			return
		}
		if i == 0 {
			if buf[i] == '>' {
				i++
			} else if buf[i] == '/' && buf[i+1] == '>' {
				i += 2
			}
			// removes white space and new lines prefixes
			for ; i < len(buf); i++ {
				if buf[i] == ' ' || buf[i] == '\n' {
					continue
				}
				break
			}
			j = i
		}
		// Search buffer.
		for ; j < len(buf); j++ {
			if buf[j] == '<' {
				if _, err = br.Discard(j); err != nil {
					err = errors.Wrap(err, "Tag Value (discard)")
					return nil, err
				}
				return buf[i:j], nil
			}
		}
		s += maxTagValueSize
	}
}

func (br *xmpReader) readTag(xmp *XMP, parent Tag) (tag Tag, err error) {
	for {
		if tag, err = br.readTagHeader(parent); err != nil {
			break
		}
		if tag.isEndTag(parent.self) {
			break
		}
		var attr Attribute
		for br.hasAttribute() {
			if attr, err = br.readAttribute(&tag); err != nil {
				return
			}
			// Parse Attribute Value
			if err = xmp.parser(attr.property); err != nil {
				return
			}
		}
		if tag.isStartTag() {
			if tag.Is(xmpns.RDFSeq) || tag.Is(xmpns.RDFAlt) || tag.Is(xmpns.RDFBag) {
				if err = br.readSeqTags(xmp, tag); err != nil {
					return
				}
			} else {
				tag.val, err = br.readTagValue()
				if err != nil {
					return
				}
				// Parse Tag Value
				if err = xmp.parser(tag.property); err != nil {
					return
				}

				if tag, err = br.readTag(xmp, tag); err != nil {
					return
				}
			}
		}
		if tag.isRootStopTag() {
			return
		}
	}
	return
}

// Special Tags
// xmpMM:History -> stEvt
// rdf:Bag -> rdf:li
// rdf:Seq -> rdf:li
// rdf:Alt -> rdf:li
func (br *xmpReader) readSeqTags(xmp *XMP, parent Tag) (err error) {
	var tag Tag
	for {
		if tag, err = br.readTagHeader(parent); err != nil {
			return
		}

		if tag.isEndTag(parent.self) {
			break
		}

		if tag.isStartTag() {
			var attr Attribute
			for br.hasAttribute() {
				if attr, err = br.readAttribute(&tag); err != nil {
					return
				}

				attr.parent = attr.self
				attr.self = parent.parent
				// Parse Attribute Value
				if err = xmp.parser(attr.property); err != nil {
					return
				}
			}

			if tag.val, err = br.readTagValue(); err != nil {
				return
			}
			tag.self = parent.parent
			tag.parent = parent.self
			// Parse Tag Value
			if err = xmp.parser(tag.property); err != nil {
				return
			}
		}
	}
	return
}

func parseAttrName(buf []byte) (xmpns.Property, int, error) {
	var a, b, c int
	for ; a < len(buf); a++ {
		if buf[a] == ' ' || buf[a] == '\n' {
			continue
		}
		break
	}
	for b = a + 1; b < len(buf); b++ {
		if buf[b] == ':' {
			break
		}
	}
	for c = b + 2; c < len(buf); c++ {
		if buf[c] == '=' || buf[c] == ' ' {
			return xmpns.IdentifyProperty(buf[a:b], buf[b+1:c]), c, nil
		}
	}
	return xmpns.Property{}, -1, ErrNegativeRead
}

func parseTagName(buf []byte) (xmpns.Property, int, error) {
	var a, b int
	for ; a < len(buf); a++ {
		if buf[a] == ':' {
			break
		}
	}
	for b = a + 1; b < len(buf); b++ {
		if buf[b] == '>' || buf[b] == ' ' || buf[b] == '\n' || buf[b] == '/' {
			return xmpns.IdentifyProperty(buf[:a], buf[a+1:b]), b, nil
		}
	}
	return xmpns.Property{}, -1, ErrNegativeRead
}
