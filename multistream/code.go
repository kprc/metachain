package multistream

import (
	"encoding/binary"
)

const(

	ConnectionSuccess uint8 = iota
	ConnectionEof
	OtherError

	ConnSync uint8 = 255

)

type Code uint32

func (c *Code)Binary() []byte  {
	buf := make([]byte, 4)

	binary.BigEndian.PutUint32(buf,uint32(*c))

	return buf
}

func SyncCode(connId uint32) *Code {
	var c Code

	cc:=&c

	cc.Encode(connId,ConnSync)

	return cc
}

func (c *Code)Decode() (uint32, uint8)  {
	id:=uint32(*c)

	errId := id >> 24

	connId := 0x00FFFFFF & id

	return connId,uint8(errId)
}

func (c *Code)Encode(connId uint32, errId uint8)  {
	id := connId | (uint32(errId) << 24)

	*c = Code(id)
}

func Buf2Code(buf []byte) *Code {
	id:=binary.BigEndian.Uint32(buf)
	c:=Code(id)
	return &c
}