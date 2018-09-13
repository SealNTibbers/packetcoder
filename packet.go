package packetcoder

import (
	. "bytes"
	"errors"
	"github.com/dgryski/go-bitstream"
)

type Packet struct {
	scheme *BitScheme
	data   *Buffer
	writer *bitstream.BitWriter
	reader *bitstream.BitReader
}

func NewPacket() *Packet {
	packet := new(Packet)
	packet.data = NewBuffer(nil)
	packet.writer = bitstream.NewWriter(packet.data)
	packet.reader = bitstream.NewReader(packet.data)
	return packet
}

func (p *Packet) SetScheme(scheme *BitScheme) {
	p.scheme = scheme
}

func (p *Packet) SetValue(fieldName string, value uint64) error {
	size, err := p.scheme.SizeOf(fieldName)
	if err != nil {
		return err
	}
	err = p.writer.WriteBits(value, int(size))
	return err
}

func (p *Packet) GetValue(fieldName string) (uint64, error) {
	size, err := p.scheme.SizeOf(fieldName)
	if err != nil {
		return 0, err
	}
	value, err := p.reader.ReadBits(int(size))
	return value, err
}

func (p *Packet) GetData() *Buffer {
	return p.data
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
