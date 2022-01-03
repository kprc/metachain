package multistream

import "io"

type MultiStream interface {

}

type MultiConn interface {
	io.Reader
	io.Writer
	io.Closer
}

type MultiStreamDialer interface {
	Dial() MultiConn
}



