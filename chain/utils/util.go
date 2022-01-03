package utils

import "crypto/rand"

func GenRandomBytes(n int) []byte  {
	buf:=make([]byte,n)

	for {
		nr,err:=rand.Read(buf)
		if err!=nil || nr != n{
			continue
		}
		break
	}

	return buf
}

func GenRandomBytes2(data []byte)  {
	l:=len(data)

	for {
		nr,err:=rand.Read(data)
		if err!=nil || nr != l{
			continue
		}
		break
	}
}