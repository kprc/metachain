package multistream


type MultiConnection struct {
	msd *MSDialer
	count int
}

func NewMultiConnection(msd *MSDialer, count int) MultiConn {
	return &MultiConnection{
		msd: msd,
		count: count,
	}
}

func (conn *MultiConnection)Close() error {
	return nil
}

func (conn *MultiConnection)Read(buf []byte) (n int,err error) {
	return 0,nil
}
func (conn *MultiConnection)Write(data []byte) (n int, err error) {
	return 0,nil
}