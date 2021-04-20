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
