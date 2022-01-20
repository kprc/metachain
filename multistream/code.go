package multistream

import (
	"encoding/binary"
)

const(

	ConnectionSuccess uint8 = iota
	ConnectionEof
	OtherError
	ConnectionClose
	ConnSync uint8 = 255
)

type Code uint64

func (c *Code)Bytes() []byte  {
	buf := make([]byte, 8)

	binary.BigEndian.PutUint64(buf,uint64(*c))

	return buf
}

func SyncCode(connId uint32) *Code {
	var c Code

	cc:=&c

	cc.Encode(connId,ConnSync,0)

	return cc
}

func SyncCloseCode(connId uint32) *Code {
	var c Code

	cc:=&c

	cc.Encode(connId,ConnectionClose,0)

	return cc
}


func (c *Code)Decode() (uint32, uint8,uint32)  {
	id:=uint64(*c)

	l:=id & 0x00000000FFFFFFFF

	errId := (id>>32) >> 24

	connId := 0x00FFFFFF & uint32(id>>32)

	return connId,uint8(errId),uint32(l)
}

func (c *Code)Encode(connId uint32, errId uint8, datalen uint32)  {
	id := connId | (uint32(errId) << 24)

	t:=uint64(id)

	t = t << 32

	t = t|uint64(datalen)

	*c = Code(t)
}

func Buf2Code(buf []byte) *Code {
	id:=binary.BigEndian.Uint64(buf)
	c:=Code(id)
	return &c
}