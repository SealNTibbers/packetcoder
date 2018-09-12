package packetcoder

import (
	. "bytes"
)

type Packet struct {
}

func (p *Packet) EncodeTo(writeBuffer *Buffer) {

}

func (p *Packet) DecodeFrom(readBuffer *Buffer) {

}

func (p *Packet) SetBitField(fieldName string, fieldSize uint) {

}

func (p *Packet) SetStuffBits(fieldSize uint) {

}
