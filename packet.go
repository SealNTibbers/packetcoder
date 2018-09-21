package packetcoder

import (
	. "bytes"
	"errors"
	"github.com/SealNTibbers/go-bitstream"
	"math/bits"
	"net"
)

func ReverseBytes(bytes []byte) []byte {
	newBytes := make([]byte, len(bytes))
	var i, j int
	for i, j = 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		newBytes[i], newBytes[j] = bytes[j], bytes[i]
	}
	size := len(bytes)
	if size%2 != 0 {
		newBytes[i] = bytes[i]
	}
	return newBytes
}

type SmartPacket interface {
	SetScheme(*BitScheme)
	GetScheme() *BitScheme
	WriteValue64(fieldName string, value uint64) error
	WriteBytes(fieldName string, value []byte) error
	WriteStuff(fieldName string) error
	ReadValue64(fieldName string) (uint64, error)
	ReadBytesValue(fieldName string) ([]byte, error)
	GetData() *Buffer
	EncodeTo(buffer *Buffer) (SmartPacket, error)
	DecodeFrom(buffer *Buffer) (SmartPacket, error)
	GetName() string

	ProcessDecoded(rawData []byte, conn net.Conn)
}

type Packet struct {
	Scheme      *BitScheme
	writeBuffer *Buffer
	readBuffer  *Reader
	bitWriter   *bitstream.BitWriter
	bitReader   *bitstream.BitReader
}

func (p *Packet) ProcessDecoded(rawData []byte, conn net.Conn) {
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

func (p *Packet) GetScheme() *BitScheme {
	return p.Scheme
}

func (p *Packet) WriteValue64(fieldName string, value uint64) error {
	field, err := p.Scheme.GetField(fieldName)
	if err != nil {
		return err
	}
	err = field.WriteValue64Into(p, value)
	return err
}

func (p *Packet) WriteBytes(fieldName string, value []byte) error {
	field, err := p.Scheme.GetField(fieldName)
	if err != nil {
		return err
	}
	err = field.WriteBytesInto(p, value)
	return err
}

func (p *Packet) WriteStuff(fieldName string) error {
	return p.WriteValue64(fieldName, 0)
}

func (p *Packet) ReadValue64(fieldName string) (uint64, error) {
	field, err := p.Scheme.GetField(fieldName)
	if err != nil {
		return 0, err
	}
	value, err := field.ReadValue64From(p)
	return value, err
}

func (p *Packet) ReadBytesValue(fieldName string) ([]byte, error) {
	field, err := p.Scheme.GetField(fieldName)
	if err != nil {
		return nil, err
	}
	value, err := field.ReadBytesValueFrom(p)
	return value, err
}

func (p *Packet) GetData() *Buffer {
	return p.writeBuffer
}

func (p *Packet) EncodeTo(buffer *Buffer) (SmartPacket, error) {
	_, err := buffer.Write(p.writeBuffer.Bytes())
	return p, err
}

func (p *Packet) DecodeFrom(buffer *Buffer) (SmartPacket, error) {
	sizeOfPacketInByte := p.Scheme.BitSize() / 8
	var packetByteArray []byte
	if int(sizeOfPacketInByte) <= len(buffer.Bytes()) {
		packetByteArray = buffer.Next(int(sizeOfPacketInByte))
	} else {
		return nil, errors.New("can't create packet because sizeInBits of the packet is larger than sizeInBits of the input buffer")
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
	sizeInBits   uint
	sizeInBytes  uint
	offset       uint
	littleEndian bool
}

func (f *bitfield) bitSize() uint {
	return f.sizeInBits
}

func (f *bitfield) IsLittleEndian() bool {
	return f.littleEndian
}

func (f *bitfield) SetLittleEndian(littleEndian bool) *bitfield {
	f.littleEndian = littleEndian
	return f
}

func (b *bitfield) WriteValue64Into(packet *Packet, value uint64) error {
	var data uint64

	if b.littleEndian {
		data = bits.ReverseBytes64(value)
		data = data >> (64 - b.sizeInBits)
	} else {
		data = value
	}
	err := packet.bitWriter.WriteBits(data, int(b.sizeInBits))
	return err
}

func (b *bitfield) WriteBytesInto(packet *Packet, value []byte) error {
	var data []byte
	var err error
	if b.littleEndian {
		data = ReverseBytes(value)
	} else {
		data = value
	}
	err = packet.bitWriter.WriteBytes(data)
	return err
}

func (b *bitfield) ReadValue64From(packet *Packet) (uint64, error) {
	var value uint64

	packet.readBuffer.Reset(packet.writeBuffer.Bytes())
	packet.bitReader.Reset(packet.readBuffer)
	tmp := packet.writeBuffer.Bytes()
	if len(tmp) > 0 {
		//
	}
	_, err := packet.bitReader.ReadBits(int(b.offset))
	if err != nil {
		return 0, err
	}
	data, err := packet.bitReader.ReadBits(int(b.sizeInBits))
	if b.littleEndian {
		value = bits.ReverseBytes64(data)
		value = value >> (64 - b.sizeInBits)
	} else {
		value = data
	}
	return value, err
}

func (b *bitfield) ReadBytesValueFrom(packet *Packet) ([]byte, error) {
	var data, value []byte

	packet.readBuffer.Reset(packet.writeBuffer.Bytes())
	packet.bitReader.Reset(packet.readBuffer)

	_, err := packet.bitReader.ReadBits(int(b.offset))
	if err != nil {
		return nil, err
	}
	var i uint
	for i = 0; i < b.sizeInBytes; i++ {
		dataByte, err := packet.bitReader.ReadByte()
		data = append(data, dataByte)
		if err != nil {
			return nil, err
		}
	}
	if b.littleEndian {
		value = ReverseBytes(data)
	} else {
		value = data
	}
	return value, err
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

func (s *BitScheme) AddBitField(fieldName string, sizeInBits uint) *bitfield {
	field := new(bitfield)
	field.name = fieldName
	field.sizeInBits = sizeInBits
	field.offset = s.size
	s.size += field.sizeInBits

	s.fields[fieldName] = field

	return field
}

func (s *BitScheme) AddBitFieldLittleEndian(fieldName string, sizeInBits uint) *bitfield {
	field := s.AddBitField(fieldName, sizeInBits)
	field.littleEndian = true
	return field
}

func (s *BitScheme) AddByteField(fieldName string, sizeInBytes uint) *bitfield {
	field := new(bitfield)
	field.name = fieldName
	field.sizeInBytes = sizeInBytes
	field.sizeInBits = sizeInBytes * 8
	field.offset = s.size
	s.size += field.sizeInBits

	s.fields[fieldName] = field

	return field
}

func (s *BitScheme) AddByteFieldLittleEndian(fieldName string, sizeInBytes uint) *bitfield {
	field := s.AddByteField(fieldName, sizeInBytes)
	field.littleEndian = true
	return field
}

func (s *BitScheme) AddStuffBits(fieldName string, fieldSize uint) {
	s.AddBitField(fieldName, fieldSize)
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

func (s *BitScheme) BitSizeOf(fieldName string) (uint, error) {
	field, err := s.GetField(fieldName)
	if err != nil {
		return 0, err
	}
	return field.sizeInBits, nil
}

func (s *BitScheme) ByteSizeOf(fieldName string) (uint, error) {
	field, err := s.GetField(fieldName)
	if err != nil {
		return 0, err
	}
	return field.sizeInBytes, nil
}

func (s *BitScheme) SizeAndOffsetOf(fieldName string) (uint, uint, error) {
	field, err := s.GetField(fieldName)
	if err != nil {
		return 0, 0, err
	}
	return field.sizeInBits, field.offset, nil
}

func (s *BitScheme) GetFields() map[string]*bitfield {
	return s.fields
}

func (s *BitScheme) GetName() string {
	return s.name
}
