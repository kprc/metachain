package keystore

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/kprc/metachain/chain/core/account"
	"github.com/kprc/metachain/chain/core/crypto"
	"github.com/kprc/metachain/chain/utils"
	"io/ioutil"
	"os"
)

type KeyProtectSalt [16]byte

func (kp KeyProtectSalt)MarshalText() ([]byte, error)  {
	b64:=base64.StdEncoding.EncodeToString(kp[:])

	return []byte(b64),nil
}

func (kp *KeyProtectSalt)UnmarshalText(text []byte)  error  {
	b,err:=base64.StdEncoding.DecodeString(string(text))
	if err!=nil{
		return err
	}

	kps:=KeyProtectSalt{}

	copy(kps[:],b)

	*kp = kps

	return nil
}

type AccountJson struct {
	Salt KeyProtectSalt `json:"salt"`
	Addr account.MetaAddr  `json:"addr"`
	CipherText []byte `json:"cipher_text"`
	Version int32     `json:"version"`
}

type KeyStore struct {
	savePath string
	acctJson *AccountJson
	account *account.Account
}

var New = func(savePath string) *KeyStore {
	return &KeyStore{
		savePath: savePath,
	}
}

func (ks *KeyStore)SetAccount(act *account.Account)  {
	ks.account = act
}

func (ks *KeyStore)GetAccount() *account.Account  {
	return ks.account
}

func (ks *KeyStore)Load() error  {
	if _,err:=os.Stat(ks.savePath);err!=nil{
		return err
	}

	var data []byte

	if f,err:=os.OpenFile(ks.savePath,os.O_RDONLY, 0755);err!=nil{
		return err
	}else{
		defer f.Close()
		if data, err = ioutil.ReadAll(f);err!=nil{
			return err
		}
	}

	aj := &AccountJson{}

	if err:=json.Unmarshal(data,aj);err!=nil{
		return err
	}

	ks.acctJson = aj

	return nil

}

func (ks *KeyStore)Open(passwd string) error  {
	if ks.acctJson == nil && ks.savePath == ""{
		return errors.New("key store not initialized")
	}

	if ks.acctJson == nil{
		if err:=ks.Load();err!=nil{
			return err
		}
	}

	var (
		key []byte
		err error
		plainText []byte
	)

	if key,err=crypto.ReinforcementPassword([]byte(passwd),ks.acctJson.Salt[:]);err!=nil{
		return err
	}

	cipherText := make([]byte,len(ks.acctJson.CipherText))
	copy(cipherText,ks.acctJson.CipherText)

	if plainText, err = crypto.Decrypt(key,cipherText);err!=nil{
		return err
	}

	act:=&account.Account{}

	var priv account.PrivKey

	if priv,err = account.Bytes2Priv(plainText);err!=nil{
		return err
	}

	act.SetPriv(priv)

	ks.account = act

	return nil
}

func (ks *KeyStore)Save() error  {
	if ks.savePath == ""{
		return errors.New("no save path")
	}

	if ks.acctJson == nil{
		return errors.New("key store have not loaded")
	}

	var (
		data []byte
		err error
		f *os.File
	)

	if data, err= json.MarshalIndent(*ks.acctJson," ","\t");err!=nil {
		return err
	}
	if f, err = os.OpenFile(ks.savePath,os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755);err!=nil{
		return err
	}
	defer f.Close()

	if _,err=f.Write(data); err!=nil{
		return err
	}

	return nil
}

func (ks *KeyStore)SaveByNewPassword(passwd string) error  {

	if ks.savePath == ""{
		return errors.New("no save path")
	}

	if ks.account == nil{
		return errors.New("key store have not opened")
	}

	ajold := ks.acctJson

	aj:= &AccountJson{}

	salt := utils.GenRandomBytes(len(aj.Salt))
	copy(aj.Salt[:],salt)

	aesk,err:=crypto.ReinforcementPassword([]byte(passwd),aj.Salt[:])
	if err!=nil{
		return err
	}
	var ciphertext []byte
	ciphertext,err = crypto.Encrypt(aesk,ks.account.GetPrivBytes())
	if err!=nil{
		return err
	}

	aj.CipherText = ciphertext
	aj.Addr = ks.account.Addr()

	ks.acctJson = aj

	if err = ks.Save();err!=nil{
		ks.acctJson = ajold
		return err
	}

	return nil
}