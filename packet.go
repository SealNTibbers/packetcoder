package packetcoder

import (
	. "bytes"
	"errors"
	"github.com/dgryski/go-bitstream"
	"math/bits"
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
	field, err := p.Scheme.GetField(fieldName)
	if err != nil {
		return err
	}
	err = field.WriteValueInto(p, value)
	return err
}

func (b *bitfield) WriteValueInto(packet *Packet, value uint64) error {
	var data uint64

	if b.littleEndian {
		data = bits.ReverseBytes64(value)
		data = data >> (64 - b.size)
	} else {
		data = value
	}
	err := packet.bitWriter.WriteBits(data, int(b.size))
	return err
}

func (p *Packet) WriteStuff(fieldName string) error {
	return p.WriteValue(fieldName, 0)
}

func (p *Packet) ReadValue(fieldName string) (uint64, error) {
	field, err := p.Scheme.GetField(fieldName)
	if err != nil {
		return 0, err
	}
	value, err := field.ReadValueFrom(p, fieldName)
	return value, err
}

func (b *bitfield) ReadValueFrom(packet *Packet, fieldName string) (uint64, error) {
	var data uint64

	packet.readBuffer.Reset(packet.writeBuffer.Bytes())
	packet.bitReader.Reset(packet.readBuffer)

	_, err := packet.bitReader.ReadBits(int(b.offset))
	if err != nil {
		return 0, err
	}
	value, err := packet.bitReader.ReadBits(int(b.size))
	if b.littleEndian {
		data = bits.ReverseBytes64(value)
		data = data >> (64 - b.size)
	} else {
		data = value
	}
	return data, err
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

func (p *Packet) GetName() string {
	if p.Scheme == nil {
		return ""
	}
	return p.Scheme.GetName()
}

type bitfield struct {
	name         string
	size         uint
	offset       uint
	littleEndian bool
}

func (f *bitfield) bitSize() uint {
	return f.size
}

func (f *bitfield) IsLittleEndian() bool {
	return f.littleEndian
}

func (f *bitfield) SetLittleEndian(littleEndian bool) *bitfield {
	f.littleEndian = littleEndian
	return f
}

type BitScheme struct {
	name   string
	fields map[string]*bitfield
	size   uint
}

func NewBitScheme(name string) *BitScheme {
	scheme := new(BitScheme)
	scheme.name = name
	scheme.fields = make(map[string]*bitfield)
	return scheme
}

func (s *BitScheme) SetBitField(fieldName string, fieldSize uint) *bitfield {
	field := new(bitfield)
	field.name = fieldName
	field.size = fieldSize
	field.offset = s.size
	s.size += field.size

	s.fields[fieldName] = field

	return field
}

func (s *BitScheme) SetBitFieldLittleEndian(fieldName string, fieldSize uint) *bitfield {
	field := s.SetBitField(fieldName, fieldSize)
	field.littleEndian = true
	return field
}

func (s *BitScheme) SetStuffBits(fieldName string, fieldSize uint) {
	s.SetBitField(fieldName, fieldSize)
}

func (s *BitScheme) BitSize() uint {
	return s.size
}

func (s *BitScheme) GetField(fieldName string) (*bitfield, error) {
	field := s.fields[fieldName]
	if field == nil {
		return nil, errors.New("there is no such field as" + fieldName + "in this scheme")
	}
	return field, nil
}

func (s *BitScheme) OffsetOf(fieldName string) (uint, error) {
	field, err := s.GetField(fieldName)
	if err != nil {
		return 0, err
	}
	return field.offset, nil
}

func (s *BitScheme) SizeOf(fieldName string) (uint, error) {
	field, err := s.GetField(fieldName)
	if err != nil {
		return 0, err
	}
	return field.size, nil
}

func (s *BitScheme) SizeAndOffsetOf(fieldName string) (uint, uint, error) {
	field, err := s.GetField(fieldName)
	if err != nil {
		return 0, 0, err
	}
	return field.size, field.offset, nil
}

func (s *BitScheme) GetFields() map[string]*bitfield {
	return s.fields
}

func (s *BitScheme) GetName() string {
	return s.name
}
