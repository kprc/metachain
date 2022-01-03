package utils

import (
	"encoding/hex"
	"testing"
)

func TestGenRandomBytes(t *testing.T) {
	randbytes:=GenRandomBytes(32)

	t.Log(len(randbytes))
	t.Log("0x"+hex.EncodeToString(randbytes))
}

func TestGenRandomBytes2(t *testing.T) {
	buf:=make([]byte,32)

	GenRandomBytes2(buf)

	t.Log(len(buf))
	t.Log("0x"+hex.EncodeToString(buf))
}