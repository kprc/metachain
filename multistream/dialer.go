package multistream

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
)

type RcvData struct {
	Length uint32
	ErrId uint8
}

type MSConn struct {
	net.Conn
	wlock sync.Mutex
	vConnCount int
	syncTunnel chan struct{}
	rcvs map[uint32]*chan *RcvData
	rcvLock sync.RWMutex
}

type MSDialer struct {
	lock sync.RWMutex
	Name string
	ConnSlot map[int]*MSConn
	MaxCount int
	RemoteAddr string
	msConnCount int
	connCount uint32
	bufLen int
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

func (msc *MSConn)read()  {
	for{
		buf:=make([]byte,8)

		if n,err:=msc.Conn.Read(buf);err!=nil || n < 4{
			msc.rcvLock.RLock()
			for _,c:=range msc.rcvs{
				*c <- &RcvData{
					ErrId: OtherError,
				}
			}
			msc.rcvLock.RUnlock()
			return
		}else{
			c:=Buf2Code(buf)
			connid,errid,l:=c.Decode()
			msc.rcvLock.RLock()
			if c,ok:=msc.rcvs[connid];!ok{
				msc.rcvLock.RUnlock()
			}else{

				*c <- &RcvData{
					Length: l,
					ErrId: errid,
				}
				if errid == ConnectionSuccess || errid == ConnectionEof{
					<-msc.syncTunnel
				}
				msc.rcvLock.RUnlock()
			}
		}
	}
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
		ConnSlot: make(map[int]*MSConn),
	}

	return msDialerStore.msd[hash],nil
}

func findMinSlot(connMap map[int]*MSConn, maxCount int) (int,bool) {
	max:= math.MinInt32
	slot := 0
	flag:=false

	for i:=0;i<maxCount;i++{
		if v,ok:=connMap[i];!ok{
			slot = i
			flag = true
			break
		}else{
			if v.vConnCount < max{
				max = v.vConnCount
				slot = i
			}
		}
	}
	return slot,flag
}

func (msd *MSDialer)Dial() (MultiConn,error)  {
	msd.lock.Lock()
	defer msd.lock.Unlock()

	minSlot,newConnFlag := findMinSlot(msd.ConnSlot,msd.MaxCount)
	if newConnFlag{
		if conn,err:=net.Dial("tcp",msd.RemoteAddr);err!=nil{
			return nil,err
		}else{
			msd.ConnSlot[minSlot] = &MSConn{
				Conn:conn,
				rcvs:make(map[uint32]*chan *RcvData),
				syncTunnel: make(chan struct{}),
			}
		}
	}

	msc := msd.ConnSlot[minSlot]

	msd.ConnSlot[minSlot].vConnCount ++
	msd.connCount ++

	rcv:=make(chan *RcvData)
	msc.rcvLock.Lock()
	msc.rcvs[msd.connCount] = &rcv
	msc.rcvLock.Unlock()

	if msc.vConnCount == 1{
		go msc.read()
	}

	conn:=&MultiConnection{
		msd: msd,
		slot: minSlot,
		connId: msd.connCount,
		rcv: &rcv,
	}

	return conn,nil
}



func (msd *MSDialer)close(conn MultiConn,slot int, connid uint32) error {
	msd.lock.Lock()
	defer msd.lock.Unlock()

	if v,ok:=msd.ConnSlot[slot];!ok{
		return errors.New("no connection in slot")
	}else{
		if v.vConnCount <= 0{
			panic("connection module error")
		}
		if v.vConnCount == 1{
			v.rcvLock.Lock()
			c:=v.rcvs[connid]
			delete(v.rcvs, connid)
			close(*c)
			v.rcvLock.Unlock()
			if err:=v.Close();err!=nil{
				fmt.Println("close connection error",err)
			}
			close(v.syncTunnel)
			delete(msd.ConnSlot,slot)
		}else{
			v.vConnCount --
			v.rcvLock.Lock()
			c:=v.rcvs[connid]
			delete(v.rcvs, connid)
			close(*c)
			v.rcvLock.Unlock()
		}
	}
	return nil
}