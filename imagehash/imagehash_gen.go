package imagehash

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Ahash) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 uint64
		zb0001, err = dc.ReadUint64()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = Ahash(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Ahash) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint64(uint64(z))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Ahash) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendUint64(o, uint64(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Ahash) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 uint64
		zb0001, bts, err = msgp.ReadUint64Bytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = Ahash(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Ahash) Msgsize() (s int) {
	s = msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Phash) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 uint64
		zb0001, err = dc.ReadUint64()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = Phash(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Phash) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint64(uint64(z))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Phash) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendUint64(o, uint64(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Phash) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 uint64
		zb0001, bts, err = msgp.ReadUint64Bytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = Phash(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Phash) Msgsize() (s int) {
	s = msgp.Uint64Size
	return
}
