package multistream

import (
	"io"
	"net"
	"reflect"
	"sync"
)

type MultiListen struct {
	listenAddr string
	l net.Listener

}

type VConn interface {
	io.Reader
	io.Writer
	io.Closer
}

type MultiVConn struct {
	conn net.Conn
	remoteAddr string
	rcvs map[uint32]*chan *RcvData
	rcvLock sync.RWMutex
}

type VConnection struct {
	vConnId uint32
	mv *MultiVConn
	rcv *chan *RcvData
}

func (mv *VConnection)Read(data []byte)(int,error)  {
	return 0,nil
}

func (mv *VConnection)Write(data []byte)(int,error)  {
	return 0,nil
}

func (mv *VConnection)Close() error  {
	return nil
}

func NewMultiListen(addr string) *MultiListen {
	return &MultiListen{
		listenAddr: addr,
	}
}

func (ml *MultiListen)ListenAndSrv()  error {
	var err error
	if ml.l,err = net.Listen("tcp",ml.listenAddr);err!=nil{
		return err
	}


	for{
		var conn net.Conn
		conn,err = ml.l.Accept()
		if err!=nil{
			return err
		}

		mvc:=&MultiVConn{
			conn: conn,
			remoteAddr: conn.RemoteAddr().String(),
			rcvs: make(map[uint32]*chan *RcvData),
		}

		go mvc.SrvConn()
		//go ml.SrvConn(conn)

	}

}

func (mvc *MultiVConn)BcastError()  {
	mvc.rcvLock.RLock()
	defer mvc.rcvLock.RUnlock()
	for _,v:=range mvc.rcvs{
		*v <- &RcvData{
			Length: 0,
			ErrId: OtherError,
		}
	}
}


func (mvc *MultiVConn)SrvConn() (VConn,error)  {
	defer mvc.conn.Close()

	buf:=make([]byte,8)
	if n,err:=mvc.conn.Read(buf);err!=nil || n != 8{
		mvc.BcastError()
		return nil,err
	}else{
		c:=Buf2Code(buf)
		id,eid,l:=c.Decode()
		if eid == ConnSync {

		}

	}



}