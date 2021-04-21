package account

import (
	"bytes"
	"errors"
	"github.com/kprc/metachain/chain/code/base58"
	"golang.org/x/crypto/sha3"
)



const(
	MetaAddrLen 	int 	= 20
	PubPrefix   	int 	= 16
	CheckSumLen 	int 	= 4
	MetaChainMagic  string 	= "0285446711860699"
	HashRound   	int 	= 3
)

type MetaAddr [MetaAddrLen]byte


func (ma *MetaAddr)Encode() string  {
	return base58.Encode(ma[:])
}

func Encode(ma MetaAddr) (string,error)  {
	if !ma.Valid(){
		return "",errors.New("check sum error, not a correct meta address")
	}

	return ma.Encode(),nil
}


func Decode(sma string) (MetaAddr,error)  {
	bma:= base58.Decode(sma)

	if len(bma) != MetaAddrLen {
		return MetaAddr{},errors.New("address length error, not a correct meta address")
	}

	ma:=MetaAddr{}

	copy(ma[:],bma)

	if !ma.Valid(){
		return MetaAddr{},errors.New("checksum error, not a correct meta address")
	}

	return ma,nil
}

func (ma MetaAddr)MarshalText() ([]byte,error)  {
	if s, err:= Encode(ma);err!=nil{
		return nil, err
	}else{
		return []byte(s),nil
	}
}

func (ma *MetaAddr)UnmarshalText(text []byte) error  {
	if ma1,err:=Decode(string(text));err!=nil{
		return err
	}else{
		*ma = ma1
	}

	return nil
}


func cshash(data []byte) ([CheckSumLen]byte,error) {
	h:=sha3.New256()

	var s []byte
	s = append(s,data...)
	s = append(s,[]byte(MetaChainMagic)...)

	for i:=0;i<HashRound;i++{
		n,err:=h.Write(s)
		if err!=nil {
			return [CheckSumLen]byte{},err
		}
		if n!= len(s){
			return [CheckSumLen]byte{},errors.New("wirte hash message failed")
		}
		s = h.Sum(nil)
	}

	cs:=[CheckSumLen]byte{}

	copy(cs[:],s[:CheckSumLen])

	return cs,nil
}

func PubKey2Addr(pk PubKey) (MetaAddr,error) {
	bprefix:=pk.Bytes()
	prefix:=bprefix[:PubPrefix]

	ma:= MetaAddr{}

	if cs,err:=cshash(prefix);err!=nil{
		return ma,err
	}else{
		copy(ma[CheckSumLen:],prefix)
		copy(ma[:CheckSumLen],cs[:])

		return ma,nil
	}
}

func (ma MetaAddr)Valid() bool  {
	if cs,err:=cshash(ma[CheckSumLen:]);err!=nil{
		return false
	}else{
		if bytes.Compare(cs[:],ma[:CheckSumLen]) == 0{
			return true
		}
	}
	return false
}
