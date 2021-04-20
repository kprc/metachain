package account

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

type Account struct {
	priv PrivKey
	pub  PubKey
	addr MetaAddr
}

var New = func() *Account {
	priv,pub:=GenerateKey()

	a:=&Account{
		priv: priv,
		pub: pub,
	}
	addr,err := PubKey2Addr(pub)
	if err!=nil{
		return nil
	}

	a.addr = addr

	return a
}

func (a *Account)SetPriv(priv PrivKey)  {
	a.priv = priv
	a.pub = priv.ToPublic()
	a.addr,_ = PubKey2Addr(a.pub)

	return
}

func (a *Account)GetPrivBytes() []byte  {
	return a.priv.Bytes()
}


func (a *Account)MetaChainAddr() string  {
	return a.addr.Encode()
}

func (a *Account)Addr() MetaAddr  {
	return a.addr
}


func Sign(message []byte, priv *PrivKey) ([]byte,error)  {
	hash:=sha3.Sum256(message)

	prive:=ecdsa.PrivateKey(*priv)

	return crypto.Sign(hash[:],&prive)
}

func (a *Account)Sign(message []byte) ([]byte,error)  {
	return Sign(message, &a.priv)
}

func RecoverPubKey(message []byte,sig []byte) (*PubKey,error)   {
	hash:=sha3.Sum256(message)

	pub,err:=crypto.SigToPub(hash[:],sig)
	if err!=nil{
		return nil, err
	}

	pk:=PubKey(*pub)

	return &pk,nil
}

func VerifySignature(message []byte,sig []byte, pk *PubKey) bool  {
	pub,err:=RecoverPubKey(message,sig)
	if err!=nil{
		return false
	}

	return pk.Cmp(pub)
}

func (a *Account)VerifySignature(message []byte,sig []byte) bool {
	return VerifySignature(message,sig,&a.pub)
}

func (a *Account)ShareKey(peer PubKey) []byte {
	return ShareKey(a.priv, peer)
}

func ShareKey(privKey PrivKey, peerPub PubKey) []byte  {
	pub:=ecdsa.PublicKey(peerPub)
	priv:=ecdsa.PrivateKey(privKey)

	X,_:=pub.Curve.ScalarMult(pub.X,pub.Y,priv.D.Bytes())

	key:=sha3.Sum256(X.Bytes())

	return key[:]
}
