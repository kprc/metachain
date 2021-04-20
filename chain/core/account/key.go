package account

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kprc/metachain/chain/code/base58"
)

type PrivKey ecdsa.PrivateKey
type PubKey ecdsa.PublicKey

func (priv PrivKey)Base58() string  {
	return base58.Encode(priv.Bytes())
}

func (pk PubKey)Base58() string  {
	return base58.Encode(pk.Bytes())
}

func (priv PrivKey)Hex() string  {
	return "0x"+hex.EncodeToString(priv.Bytes())
}

func (pk PubKey)Hex() string  {
	return "0x"+hex.EncodeToString(pk.Bytes())
}

func (pk PubKey)Bytes() []byte{
	pke:=ecdsa.PublicKey(pk)
	return crypto.FromECDSAPub(&pke)
}

func (priv PrivKey)Bytes() []byte  {
	prive:=ecdsa.PrivateKey(priv)
	return crypto.FromECDSA(&prive)
}

func (priv PrivKey)ToPublic() PubKey  {
	prive:=ecdsa.PrivateKey(priv)

	pub:=prive.PublicKey

	return PubKey(pub)
}

func Bytes2Priv(privBytes []byte) (PrivKey,error)  {
	priv,err:=crypto.ToECDSA(privBytes)
	if err!=nil{
		return PrivKey{}, err
	}
	return PrivKey(*priv),nil
}

func Bytes2Pub(pubBytes []byte) (PubKey,error)  {
	pub,err:=crypto.UnmarshalPubkey(pubBytes)
	if err!=nil{
		return PubKey{}, err
	}

	return PubKey(*pub),nil
}

func GenerateKey() (PrivKey,PubKey)  {
	var key *ecdsa.PrivateKey
	var err error

	for{
		key,err =crypto.GenerateKey()
		if err != nil{
			continue
		}
		break
	}

	pub:=(key.Public()).(*ecdsa.PublicKey)

	return PrivKey(*key),PubKey(*pub)
}

func (pk *PubKey)Cmp(pk1 *PubKey) bool  {
	pkb := pk.Bytes()
	pkb1 :=pk1.Bytes()

	if bytes.Compare(pkb,pkb1) == 0{
		return true
	}

	return false
}