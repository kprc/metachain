package keystore

import (
	"github.com/kprc/metachain/chain/core/account"
	"testing"
)


func TestNewKeyStore(t *testing.T)  {
	ks:=New("./key.json")

	act:=account.New()

	ks.SetAccount(act)

	t.Log(act.MetaChainAddr())

	err:=ks.SaveByNewPassword("123")
	if err!=nil{
		t.Fatal(err.Error())
	}
}

func TestLoadKeyStore(t *testing.T)  {
	ks:=New("./key.json")

	if err:=ks.Load();err!=nil{
		t.Fatal(err.Error())
	}

	if err:=ks.Open("123");err!=nil{
		t.Fatal(err.Error())
	}

	act:=ks.GetAccount()

	t.Log(act.MetaChainAddr())

}

