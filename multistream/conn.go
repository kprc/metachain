package multistream

import (
	"errors"
	"io"
)

type MultiConnection struct {
	msd *MSDialer
	slot int
	connId uint32
	rcv *chan *RcvData
}

func NewMultiConnection(msd *MSDialer, slot int,connId uint32) MultiConn {
	return &MultiConnection{
		msd: msd,
		slot: slot,
		connId: connId,
	}
}

func (conn *MultiConnection)Close() error {
	return conn.msd.close(conn,conn.slot,conn.connId)
}


func (conn *MultiConnection)Read(buf []byte) (n int,err error) {
	select {
	case data :=<-*conn.rcv:
		switch data.ErrId {
		case ConnectionSuccess:
			n:=copy(buf,data.Data)
			return n,nil
		case ConnectionEof:
			n:=copy(buf,data.Data)
			return n,io.EOF
		case OtherError:
			return 0,errors.New("connection fatal error")
		}
	}
	return 0,nil
}
func (conn *MultiConnection)Write(data []byte) (n int, err error) {
	if conn.msd == nil{
		panic("unexpect error")
	}
	conn.msd.lock.RLock()
	defer conn.msd.lock.RUnlock()

	if msc,ok:=conn.msd.ConnSlot[conn.slot];!ok{
		return 0,errors.New("slot not exists")
	}else{
		return msc.Write(data)
	}
}

func (conn *MultiConnection)ConnectTo() (err error)  {
	_,err = conn.Write(SyncCode(conn.connId).Binary())
	if err!=nil{
		return err
	}

	return nil
}
