package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/scrypt"
	"io"
)

func ReinforcementPassword(password, salt []byte) ([]byte,error) {
	return scrypt.Key(password,salt,32768, 8, 1, 32)
}

func Encrypt(key []byte,plainText []byte) ([]byte,error)  {
	blk,err:=aes.NewCipher(key)
	if err!=nil{
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv:=cipherText[:aes.BlockSize]
	_,err = io.ReadFull(rand.Reader, iv)
	if err!=nil{
		return nil, err
	}

	stream:=cipher.NewCFBEncrypter(blk,iv)

	stream.XORKeyStream(cipherText[aes.BlockSize:],plainText)

	return cipherText,nil

}

func Decrypt(key, cipherText []byte) (plainText []byte, err error)  {
	if len(cipherText)<aes.BlockSize{
		return nil,errors.New("cipher text too short")
	}

	var blk cipher.Block

	blk, err = aes.NewCipher(key)
	if err!= nil{
		return nil, err
	}

	iv:=cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream:=cipher.NewCFBDecrypter(blk,iv)
	stream.XORKeyStream(cipherText,cipherText)

	return cipherText,nil
}