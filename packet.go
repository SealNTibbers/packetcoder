package packetcoder

import (
	. "bytes"
	"errors"
)

type Packet struct {
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

func (s *Packet) EncodeTo(writeBuffer *Buffer) {

}

func (s *Packet) DecodeFrom(readBuffer *Buffer) {

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
