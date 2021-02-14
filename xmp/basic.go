package xmp

import (
	"time"

	"github.com/evanoberholster/imagemeta/xmp/xmpns"
	uuid "github.com/satori/go.uuid"
)

// UUID is a [16]byte Universally Unique Identifier (UUID).
// Based on github.com/satori/go.uuid
type UUID uuid.UUID

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (u UUID) MarshalText() (text []byte, err error) {
	return uuid.UUID(u).MarshalText()
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (u *UUID) UnmarshalText(text []byte) (err error) {
	uid, err := uuid.FromString(string(text))
	*u = UUID(uid)
	return err
}

func (basic *Basic) decode(p property) (err error) {
	switch p.Name() {
	case xmpns.CreateDate:
		basic.CreateDate, err = parseDate(p.val)
	case xmpns.CreatorTool:
		basic.CreatorTool = parseString(p.val)
	case xmpns.Label:
		basic.Label = parseString(p.val)
	case xmpns.MetadataDate:
		basic.MetadataDate, err = parseDate(p.val)
	case xmpns.ModifyDate:
		basic.ModifyDate, err = parseDate(p.val)
	case xmpns.Rating:
		basic.Rating = int8(parseInt(p.val))
	default:
		return ErrPropertyNotSet
	}
	return
}

func (mm *XMPMM) decode(p property) (err error) {
	switch p.Name() {
	case xmpns.DocumentID:
		mm.DocumentID = parseUUID(p.val)
	case xmpns.OriginalDocumentID:
		mm.OriginalDocumentID = parseUUID(p.val)
	case xmpns.InstanceID:
		mm.InstanceID = parseUUID(p.val)
	default:
		return ErrPropertyNotSet
	}
	return
}

// Basic - the XMP basic namespace contains properties that provide basic descriptive information.
// XMP spec Section 8.4
// xmlns:xmp="http://ns.adobe.com/xap/1.0/"
type Basic struct {
	// The date and time the resource was created. For a digital file, this need not match a
	// file-system  creation time. For a freshly created resource, it should be close to that time,
	// modulo the time taken to write the file. Later file transfer, copying, and so on, can make the
	// file-system time arbitrarily different.
	CreateDate time.Time `xml:"CreateDate"`
	// The name of the first known tool used to create the resource.
	CreatorTool string `xml:"CreatorTool"`
	// A word or short phrase that identifies a resource as a member of a user-defined collection.
	Label string `xml:"Label,attr"`
	// The date and time that any metadata for this resource was last changed.
	// It should be the same as or more recent than xmp:ModifyDate.
	MetadataDate time.Time `xml:"MetadataDate,attr"`
	// The date and time the resource was last modified.
	ModifyDate time.Time `xml:"ModifyDate,attr"`
	// A user-assigned rating for this file. The value shall be -1 or in the range [0..5],
	// where -1 indicates “rejected” and 0 indicates “unrated”. If xmp:Rating is not present,
	// a value of 0 should be assumed.
	Rating int8 `xml:"Rating,attr"`
}

// XMPMM - The XMP Media Management namespace contains properties that provide information
// regarding the identification, composition, and history of a resource.
// XMP spec Section 8.6
// Incomplete
type XMPMM struct {
	// DocumentId is the common identifier for all versions and renditions of a resource.
	DocumentID UUID
	// InstanceId is an identifier for a specific incarnation of a resource,
	// updated each time a file is saved.
	InstanceID UUID
	// OriginalDocumentId is the common identifier for the original resource from which the current
	// resource is derived. For example, if you save a resource to a different format,
	// then save that one to another format, each save operation should generate a new
	// xmpMM:DocumentID that uniquely identifies the resource in that format,
	// but should retain the ID of the source file here.
	OriginalDocumentID UUID
}