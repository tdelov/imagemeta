package xmp

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/tdelov/imagemeta/meta"
	"github.com/tdelov/imagemeta/xmp/xmpns"
)

func (xmp *XMP) parser(p property) (err error) {
	if len(p.Value()) == 0 {
		return
	}
	switch p.Namespace() {
	case xmpns.XMLnsNS:
		return // Null operation
	case xmpns.ExifNS:
		err = xmp.Exif.parse(p)
	case xmpns.AuxNS:
		err = xmp.Aux.parse(p)
	case xmpns.DcNS:
		err = xmp.DC.parse(p)
	case xmpns.XmpNS, xmpns.XapNS:
		err = xmp.Basic.parse(p)
	case xmpns.TiffNS:
		err = xmp.Tiff.parse(p)
	case xmpns.CrsNS:
		err = xmp.CRS.parse(p)
	case xmpns.XmpMMNS, xmpns.XapMMNS:
		err = xmp.MM.parse(p)
	default:
		//fmt.Println(p, ns)
		return
	}
	if err != nil {
		err = nil
		//fmt.Println(err, "\t", p)
	}

	return
}

// parseDate parses a Date and returns a time.Time or an error
//func parseDate(buf []byte) (t time.Time, err error) {
//	str := string(buf)
//	if t, err = time.Parse("2006-01-02T15:04:05Z07:00", str); err != nil {
//		if t, err = time.Parse("2006-01-02T15:04:05.00", str); err != nil {
//			return time.Parse("2006-01-02T15:04:05", str)
//		}
//	}
//	return
//}

// This function returns a EXIF date time compliant date.
// Another function may be needed to return XPM compliant dates
func parseDate(buf string) (t time.Time, err error) {

	// "2021-01-10T17:30:57.00"
        // Not sure if this time format is part of the specs. more common for IPTC or XMP format to have the dashes
        // XMP format could have the time zone
        // https://developer.adobe.com/xmp/docs/XMPNamespaces/XMPDataTypes/#date
        // https://iptc.org/std/photometadata/documentation/mappingguidelines/#exif-note-on-date-created
	if buf[4] == '-' && buf[7] == '-' &&	buf[13] == ':' && buf[16] == ':' {
	 	if buf[10] == ' ' || buf[10] == 'T' { 
			year, err := strconv.Atoi( buf[0:4]   )
			month,err := strconv.Atoi( buf[5:7]   )
			day,  err := strconv.Atoi( buf[8:10]  )
			hour, err := strconv.Atoi( buf[11:13] )
			min,  err := strconv.Atoi( buf[14:16] )
			sec,  err := strconv.Atoi( buf[17:19] )
			// mil, err := strconv.Atoi(  buf[20:] )

			//fmt.Println( year ,"|", month, "|", day, "|", hour, "|", min, "|", sec ,"|" , 0 )
			return time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), 0 , time.UTC) , err
		}
	}

	// EXIF Date format YYYY:MM:DD HH:mm:SS
	// https://exiv2.org/manpage.html#date_time_fmts
	// buf[10] should be a space according to the specs, but they also say in the Makernote, the format could be different per manufacturers.
	// Exif date-time values have no time zone information.
	if buf[4] == ':' && buf[7] == ':' && buf[10] == ' ' &&	buf[13] == ':' && buf[16] == ':' {
		 if buf[10] == ' ' || buf[10] == 'T' { 
			year, err := strconv.Atoi( buf[0:4]   )
			month,err := strconv.Atoi( buf[5:7]   )
			day,  err := strconv.Atoi( buf[8:10]  )
			hour, err := strconv.Atoi( buf[11:13] )
			min,  err := strconv.Atoi( buf[14:16] )
			sec,  err := strconv.Atoi( buf[17:19] )
			// mil, err := strconv.Atoi(  buf[20:] )

			//fmt.Println( year ,"|", month, "|", day, "|", hour, "|", min, "|", sec, "|" , 0 )
			return time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), 0 , time.UTC) , err
		}
	}
	return  
}
// --------------------------------------------------------------------- //

// parseUUID parses a UUID and returns a meta.UUID
func parseUUID(buf []byte) (uuid meta.UUID) {
	if _, b := readUntil(buf, ':'); len(b) > 0 {
		buf = b
	}
	err := uuid.UnmarshalText(buf)
	if err != nil {
		if DebugMode {
			fmt.Println("Parse UUID error: ", err)
		}
	}
	return
}

// parseInt parses a []byte of a string representation of an int64 value and returns the value
func parseInt(buf []byte) (i int64) {
	if buf[0] == '-' {
		buf = buf[1:]
		i = -1
	}
	i *= int64(parseUint(buf))
	return
}

// parseUint parses a []byte of a string representation of a uint64 value and returns the value.
func parseUint(buf []byte) (u uint64) {
	for i := 0; i < len(buf); i++ {
		u *= 10
		u += uint64(buf[i] - '0')
	}
	return
}

// parseUint32 parses a []byte of a string representation of a uint32 value and returns the value.
// If the value is larger than uint32 returns 0.
func parseUint32(buf []byte) (u uint32) {
	if i := parseUint(buf); i < math.MaxUint32 {
		return uint32(i)
	}
	return 0
}

// parseUint8 parses a []byte of a string representation of a uint8 value and returns the value.
// If the value is larger than uint8 returns 0.
func parseUint8(buf []byte) (u uint8) {
	if i := parseUint(buf); i < math.MaxUint8 {
		return uint8(i)
	}
	return 0
}

// parseFloat64 parses a []byte of a string representation of a float64 value and returns the value
func parseFloat64(buf []byte) (f float64) {
	f, err := strconv.ParseFloat(string(buf), 64)
	if err != nil {
		return 0.0
	}
	return
}

// parseString parses a []byte and returns a string
func parseString(buf []byte) string {
	return string(buf)
}

// parseRational separates a string into a fraction.
// With "n" as the numerator and "d" as the denominator.
// TODO: Improve parsing functionality
func parseRational(buf []byte) (n uint32, d uint32) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == '/' {
			if i < len(buf)+1 {
				n = uint32(parseUint(buf[:i]))
				d = uint32(parseUint(buf[i+1:]))
				return
			}
		}
	}
	return
}

func readUntil(buf []byte, delimiter byte) (a []byte, b []byte) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == delimiter || buf[i] == '>' {
			return buf[:i], buf[i+1:]
		}
	}
	return buf, nil
}
