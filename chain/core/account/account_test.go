package account

import (
	"bytes"
	"github.com/kprc/metachain/chain/code/base58"
	"testing"
)

func TestCreateAccount(t *testing.T)  {
	act:=New()

	t.Log(act.priv.Base58())
	t.Log(act.pub.Base58())
}

func TestAddr(t *testing.T)  {
	act:=New()

	ma,err:=PubKey2Addr(act.pub)
	if err!=nil{
		t.Error(err.Error())
	}

	t.Log(ma.Valid())
	addr:=ma.Encode()
	t.Log(addr)

	ma,err = Decode(addr)
	if err!=nil{
		t.Error(err.Error())
	}

	var s string
	s,err = Encode(ma)
	if err!=nil{
		t.Error(err.Error())
	}else{
		t.Log(s)
	}
}

func TestKey(t *testing.T)  {
	act:=New()

	t.Log(len(act.pub.Bytes()))
	t.Log(len(act.priv.Bytes()))
	t.Log(act.pub.Hex())
	t.Log(act.priv.Hex())

	if pub,err:=Bytes2Pub(base58.Decode(act.pub.Base58()));err!=nil{
		t.Error(err.Error())
	}else{
		t.Log(pub.Hex())
	}
	if priv,err:=Bytes2Priv(base58.Decode(act.priv.Base58()));err!=nil{
		t.Error(err.Error())
	}else{
		t.Log(priv.Hex())
	}
}

func TestKeyCmp(t *testing.T)  {
	act1:=New()
	act2:=New()

	if !act1.pub.Cmp(&act2.pub){
		t.Log("not equals")
	}

	if act1.pub.Cmp(&act1.pub){
		t.Log("equals")
	}
}

func TestMetaAddr(t *testing.T)  {
	act:=New()

	t.Log(act.MetaChainAddr())

	if ma,err:=PubKey2Addr(act.pub);err!=nil{
		t.Error(err.Error())
	}else{
		t.Log(ma.Encode())
	}
}

func TestSign(t *testing.T)  {
	message:="hello world"

	act:=New()

	if sig,err:=act.Sign([]byte(message));err!=nil{
		t.Error(err.Error())
	}else{
		t.Log(act.VerifySignature([]byte(message),sig))
	}

	t.Log(act.pub.Base58())

	if sig,err:=Sign([]byte(message),&act.priv);err!=nil{
		t.Error(err.Error())
	}else{
		t.Log(VerifySignature([]byte(message),sig,&act.pub))

		if pk,err:=RecoverPubKey([]byte(message),sig);err!=nil{
			t.Error(err.Error())
		}else{
			t.Log(pk.Base58())
		}

	}

}

func TestShareKey(t *testing.T)  {
	act1:=New()
	act2:=New()

	s1:=act1.ShareKey(act2.pub)
	s2:=act2.ShareKey(act1.pub)

	if e:=bytes.Compare(s1,s2);e==0{
		t.Log("correct share key")
	}

}

func TestCreateRandomAccuont(t *testing.T)  {
	privMem := make(map[PrivKey]struct{})
	addrMem := make(map[MetaAddr]struct{})

	duppriv := 0
	dupaddr := 0

	for i:=0;i<1000000;i++{
		act:=New()

		if _,ok:=privMem[act.priv];ok{
			duppriv ++
			t.Error("duplication priv")
		}else{
			privMem[act.priv] = struct{}{}
		}

		if _,ok:=addrMem[act.addr];ok{
			dupaddr ++
			t.Error("duplication addr")
		}else{
			addrMem[act.addr] = struct{}{}
		}

	}

	t.Log(duppriv)
	t.Log(dupaddr)

}


