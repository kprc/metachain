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
			n,err:= conn.msd.ConnSlot[conn.slot].Read(buf[:data.Length])
			conn.msd.ConnSlot[conn.slot].syncTunnel <- struct{}{}
			return n,err
		case ConnectionEof:
			var n int
			var err error
			if data.Length > 0{
				n,err =conn.msd.ConnSlot[conn.slot].Read(buf[:data.Length])
			}
			conn.msd.ConnSlot[conn.slot].syncTunnel <- struct{}{}
			if err!=nil{
				return n,err
			}
			return n,io.EOF
		case OtherError:
			return 0,errors.New("connection fatal error")
		case ConnectionClose:
			return 0,io.EOF
		}
	}
	return 0,nil
}

func (conn *MultiConnection)Write(data []byte) (n int, err error) {
	if conn.msd == nil{
		panic("unexpect error")
	}
	conn.msd.lock.RLock()
	if msc,ok:=conn.msd.ConnSlot[conn.slot];!ok{
		conn.msd.lock.RUnlock()
		return 0,errors.New("slot not exists")
	}else{
		conn.msd.lock.RUnlock()
		return msc.Write(data)
	}
}

