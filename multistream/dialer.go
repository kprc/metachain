package multistream

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"net"
	"sync"
)

type MSDialer struct {
	Name string
	Conn []*net.Conn
	MaxCount int
	RemoteAddr string
	connCount int
	connMap map[int]int
}



type MHash [32]byte

type MSDialerStore struct {
	msd map[MHash]*MSDialer
	rwLock sync.RWMutex
}

var (
	msDialerStore *MSDialerStore
	lock sync.Mutex
)

func init()  {
	msDialerStore = &MSDialerStore{
		msd: make(map[MHash]*MSDialer),
	}
}

func dialString(dialName string, count int) []byte  {
	buf:=make([]byte,len(dialName)+4)
	binary.BigEndian.PutUint32(buf,uint32(count))

	copy(buf[4:],dialName)

	return buf
}

func hashbyte(data []byte) MHash {
	return MHash(sha256.Sum256(data))
}

func NewDialer(dialerName string, count int, remoteAddr string) (MultiStreamDialer,error) {
	data:=dialString(dialerName,count)
	hash:=hashbyte(data)

	msDialerStore.rwLock.RLock()
	if _,ok:=msDialerStore.msd[hash];ok{
		msDialerStore.rwLock.RUnlock()
		return nil,errors.New("duplicate dialer")
	}
	msDialerStore.rwLock.RUnlock()

	msDialerStore.rwLock.Lock()
	defer msDialerStore.rwLock.Unlock()

	if _,ok:=msDialerStore.msd[hash];ok{
		return nil,errors.New("duplicate dialer")
	}

	msDialerStore.msd[hash] = &MSDialer{
		Name: dialerName,
		MaxCount: count,
		RemoteAddr: remoteAddr,
		connMap: make(map[int]int),
	}

	return msDialerStore.msd[hash],nil

}


func (msd *MSDialer)Dial() MultiConn  {
	if len(msd.Conn) < msd.MaxCount{
		
	}
}

