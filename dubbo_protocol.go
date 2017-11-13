package dubbo_thrift

import (
	"encoding/binary"
	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	Version                  = 1
	Magic                    = -9540
	DefaultHeaderLength      = 1
	DefaultMessageLength     = 1
	MessageLengthIndex       = 6
	MessageHeaderLengthIndex = 10
)

type TDubboProtocolFactory struct {
	serviceName string
}

func NewTDubboProtocolFactory(serviceName string) *TDubboProtocolFactory {
	return &TDubboProtocolFactory{serviceName: serviceName}
}

func (p *TDubboProtocolFactory) GetProtocol(trans thrift.TTransport) thrift.TProtocol {
	return NewTDubboProtocol(trans, p.serviceName)
}

type TDubboProtocol struct {
	*thrift.TBinaryProtocol
	buffer      [64]byte
	serviceName string
}

func NewTDubboProtocol(trans thrift.TTransport, serviceName string) *TDubboProtocol {
	p := &TDubboProtocol{
		TBinaryProtocol: thrift.NewTBinaryProtocol(trans, false, true),
		serviceName:     serviceName,
	}
	return p
}

func (p *TDubboProtocol) WriteDubboHeader(seqId int32) error {
	if err := p.WriteI32(DefaultMessageLength); err != nil { //message长度 i32最大值
		return err
	}
	if err := p.WriteI16(Magic); err != nil { //magic数
		return err
	}
	if err := p.WriteI32(DefaultMessageLength); err != nil { //message长度 i32最大值
		return err
	}
	if err := p.WriteI16(DefaultHeaderLength); err != nil { //message header长度 i16最大值
		return err
	}
	if err := p.WriteByte(Version); err != nil { //Version 版本号
		return err
	}
	if err := p.WriteString(p.serviceName); err != nil { //服务名
		return err
	}
	if err := p.WriteI64(int64(seqId)); err != nil { // request id
		return err
	}
	p.FillHeaderLength()
	return nil
}

func (p *TDubboProtocol) ReadDubboHeader() error {
	if _, err := p.ReadI32(); err != nil { //message长度 i32最大值
		return err
	}
	if _, err := p.ReadI16(); err != nil { //magic数
		return err
	}
	if _, err := p.ReadI32(); err != nil { //message长度 i32最大值
		return err
	}
	if _, err := p.ReadI16(); err != nil { //message header长度 i16最大值
		return err
	}
	if _, err := p.ReadByte(); err != nil { //Version 版本号
		return err
	}
	if _, err := p.ReadString(); err != nil { //服务名
		return err
	}
	if _, err := p.ReadI64(); err != nil { // request id
		return err
	}
	return nil
}

func (p *TDubboProtocol) FillHeaderLength() {
	if trans, ok := p.Transport().(*TDubboTransport); ok {
		size := trans.Len() - 4
		buf := trans.Bytes()[MessageHeaderLengthIndex:(MessageHeaderLengthIndex + 2)]
		binary.BigEndian.PutUint16(buf, uint16(size))
	}
}

func (p *TDubboProtocol) FillMessageLength() {
	if trans, ok := p.Transport().(*TDubboTransport); ok {
		size := trans.Len() - 4
		buf := trans.Bytes()[0:4]
		binary.BigEndian.PutUint32(buf, uint32(size))
		buf = trans.Bytes()[MessageLengthIndex:(MessageLengthIndex + 4)]
		binary.BigEndian.PutUint32(buf, uint32(size))
	}
}

func (p *TDubboProtocol) WriteMessageBegin(name string, typeId thrift.TMessageType, seqId int32) error {
	if err := p.WriteDubboHeader(seqId); err != nil {
		return err
	}
	if err := p.TBinaryProtocol.WriteMessageBegin(name, typeId, seqId); err != nil {
		return err
	}
	return nil
}

func (p *TDubboProtocol) WriteMessageEnd() error {
	p.FillMessageLength()
	return nil
}

func (p *TDubboProtocol) ReadMessageBegin() (name string, typeId thrift.TMessageType, seqid int32, err error) {
	if err := p.ReadDubboHeader(); err != nil {
		return "", typeId, 0, thrift.NewTProtocolException(err)
	}
	return p.TBinaryProtocol.ReadMessageBegin()
}
