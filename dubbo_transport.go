package dubbo_thrift

import (
	"bufio"
	"bytes"
	"git.apache.org/thrift.git/lib/go/thrift"
)

type TDubboTransportFactory struct {
	bufferSize int
}

// Return a wrapped instance of the base Transport.
func (p *TDubboTransportFactory) GetTransport(trans thrift.TTransport) thrift.TTransport {
	return NewTDubboTransport(trans, p.bufferSize)
}

func NewTDubboTransportFactory(bufferSize int) thrift.TTransportFactory {
	return &TDubboTransportFactory{bufferSize: bufferSize}
}

type TDubboTransport struct {
	trans  thrift.TTransport
	buf    bytes.Buffer
	reader *bufio.Reader
}

func NewTDubboTransport(trans thrift.TTransport, bufferSize int) *TDubboTransport {
	dubboTrans := new(TDubboTransport)
	dubboTrans.trans = trans
	dubboTrans.buf.Grow(bufferSize)
	dubboTrans.reader = bufio.NewReaderSize(trans, bufferSize)
	return dubboTrans
}

func (p *TDubboTransport) Open() error {
	return p.trans.Open()
}

func (p *TDubboTransport) IsOpen() bool {
	return p.trans.IsOpen()
}

func (p *TDubboTransport) Close() error {
	return p.trans.Close()
}

func (p *TDubboTransport) Write(buf []byte) (int, error) {
	n, err := p.buf.Write(buf)
	return n, thrift.NewTTransportExceptionFromError(err)
}

func (p *TDubboTransport) Flush() error {
	size := p.buf.Len()
	if size > 0 {
		if _, err := p.buf.WriteTo(p.trans); err != nil {
			return thrift.NewTTransportExceptionFromError(err)
		}
	}
	err := p.trans.Flush()
	return thrift.NewTTransportExceptionFromError(err)
}

func (p *TDubboTransport) Bytes() []byte {
	return p.buf.Bytes()
}

func (p *TDubboTransport) Len() int {
	return p.buf.Len()
}

func (p *TDubboTransport) Read(b []byte) (int, error) {
	n, err := p.reader.Read(b)
	if err != nil {
		p.reader.Reset(p.trans)
	}
	return n, err
}

func (p *TDubboTransport) RemainingBytes() uint64 {
	return p.trans.RemainingBytes()
}
