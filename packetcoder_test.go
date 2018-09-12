package packetcoder

import (
	"bytes"
	"testing"
)

func TestPacket(t *testing.T) {
	packet := new(Packet)
	t.Logf("%T", packet)
}

func TestPacketEncode(t *testing.T) {
	packet := new(Packet)
	buf := bytes.NewBuffer(nil)
	packet.EncodeTo(buf)
}

func TestPacketDecode(t *testing.T) {
	packet := new(Packet)
	buf := bytes.NewBuffer(nil)
	packet.DecodeFrom(buf)
}

func TestPacketSetFields(t *testing.T) {
	packet := new(Packet)
	packet.SetBitField("head", 4)
	packet.SetBitField("type", 8)
	packet.SetStuffBits(4)
}
