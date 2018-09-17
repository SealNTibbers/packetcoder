package packetcoder

import (
	. "bytes"
	"errors"
	"github.com/dgryski/go-bitstream"
)

type Packet struct {
	Scheme      *BitScheme
	writeBuffer *Buffer
	readBuffer  *Reader
	bitWriter   *bitstream.BitWriter
	bitReader   *bitstream.BitReader
}

func NewPacket() *Packet {
	packet := new(Packet)
	var data []byte
	packet.writeBuffer = NewBuffer(data)
	packet.readBuffer = NewReader(data)
	packet.bitWriter = bitstream.NewWriter(packet.writeBuffer)
	packet.bitReader = bitstream.NewReader(packet.readBuffer)
	return packet
}

func NewPacketFor(scheme *BitScheme) *Packet {
	packet := NewPacket()
	packet.SetScheme(scheme)
	return packet
}

func (p *Packet) SetScheme(scheme *BitScheme) {
	p.Scheme = scheme
}

func (p *Packet) WriteValue(fieldName string, value uint64) error {
	size, _, err := p.Scheme.SizeAndOffsetOf(fieldName)
	if err != nil {
		return err
	}
	err = p.bitWriter.WriteBits(value, int(size))
	return err
}

func (p *Packet) WriteStuff(fieldName string) error {
	return p.WriteValue(fieldName, 0)
}

func (p *Packet) ReadValue(fieldName string) (uint64, error) {
	size, offset, err := p.Scheme.SizeAndOffsetOf(fieldName)
	if err != nil {
		return 0, err
	}

	p.readBuffer.Reset(p.writeBuffer.Bytes())
	p.bitReader.Reset(p.readBuffer)

	_, err = p.bitReader.ReadBits(int(offset))
	if err != nil {
		return 0, err
	}
	value, err := p.bitReader.ReadBits(int(size))
	return value, err
}

func (p *Packet) GetData() *Buffer {
	return p.writeBuffer
}

func (p *Packet) EncodeTo(buffer *Buffer) (*Packet, error) {
	_, err := buffer.Write(p.writeBuffer.Bytes())
	return p, err
}

func (p *Packet) DecodeFrom(buffer *Buffer) (*Packet, error) {
	sizeOfPacketInByte := p.Scheme.BitSize() / 8
	var packetByteArray []byte
	if int(sizeOfPacketInByte) <= len(buffer.Bytes()) {
		packetByteArray = buffer.Next(int(sizeOfPacketInByte))
	} else {
		return nil, errors.New("can't create packet because size of the packet is larger than size of the input buffer")
	}
	p.writeBuffer = NewBuffer(packetByteArray)
	p.readBuffer = NewReader(packetByteArray)
	p.bitWriter = bitstream.NewWriter(p.writeBuffer)
	p.bitReader = bitstream.NewReader(p.readBuffer)
	return p, nil
}

type bitfield struct {
	name   string
	size   uint
	offset uint
}

func (f *bitfield) bitSize() uint {
	return f.size
}

type BitScheme struct {
	fields map[string]*bitfield
	size   uint
}

func NewBitScheme() *BitScheme {
	scheme := new(BitScheme)
	scheme.fields = make(map[string]*bitfield)
	return scheme
}

func (s *BitScheme) SetBitField(fieldName string, fieldSize uint) {
	field := new(bitfield)
	field.name = fieldName
	field.size = fieldSize
	field.offset = s.size
	s.size += field.size

	s.fields[fieldName] = field
}

func (s *BitScheme) SetStuffBits(fieldName string, fieldSize uint) {
	s.SetBitField(fieldName, fieldSize)
}

func (s *BitScheme) BitSize() uint {
	return s.size
}

func (s *BitScheme) OffsetOf(fieldName string) (uint, error) {
	field := s.fields[fieldName]
	if field == nil {
		return 0, errors.New("there is no such field as" + fieldName + "in this scheme")
	}
	return field.offset, nil
}

func (s *BitScheme) SizeOf(fieldName string) (uint, error) {
	field := s.fields[fieldName]
	if field == nil {
		return 0, errors.New("there is no such field as" + fieldName + "in this scheme")
	}
	return field.size, nil
}

func (s *BitScheme) SizeAndOffsetOf(fieldName string) (uint, uint, error) {
	field := s.fields[fieldName]
	if field == nil {
		return 0, 0, errors.New("there is no such field as" + fieldName + "in this scheme")
	}
	return field.size, field.offset, nil
}

func (s *BitScheme) GetFields() map[string]*bitfield {
	return s.fields
}
