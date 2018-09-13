package packetcoder

import (
	. "bytes"
	"errors"
	"fmt"
	"github.com/dgryski/go-bitstream"
)

type Packet struct {
	scheme      *BitScheme
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

func (p *Packet) SetScheme(scheme *BitScheme) {
	p.scheme = scheme
}

func (p *Packet) WriteValue(fieldName string, value uint64) error {
	size, err := p.scheme.SizeOf(fieldName)
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
	size, offset, err := p.scheme.SizeAndOffsetOf(fieldName)
	if err != nil {
		return 0, err
	}

	p.readBuffer.Reset(p.writeBuffer.Bytes())
	p.bitReader.Reset(p.readBuffer)

	skippedValue, err := p.bitReader.ReadBits(int(offset))
	fmt.Println(skippedValue)
	if err != nil {
		return 0, err
	}
	value, err := p.bitReader.ReadBits(int(size))
	fmt.Println(value)
	return value, err
}

func (p *Packet) GetData() *Buffer {
	return p.writeBuffer
}

func (p *Packet) EncodeTo(writeBuffer *Buffer) {

}

func (p *Packet) DecodeFrom(writeBuffer *Buffer) {

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
