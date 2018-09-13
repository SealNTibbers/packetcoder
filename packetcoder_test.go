package packetcoder

import (
	"bytes"
	"github.com/SealNTibbers/packetcoder/testutils"
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

func TestSchemeSetFields(t *testing.T) {
	scheme := NewBitScheme()
	scheme.SetBitField("head", 4)
	scheme.SetBitField("type", 8)
	scheme.SetStuffBits("fill", 4)
}

func TestSchemeBitSize(t *testing.T) {
	scheme := NewBitScheme()
	scheme.SetBitField("head", 4)
	scheme.SetBitField("type", 8)
	scheme.SetStuffBits("fill", 4)

	var size uint
	size = scheme.BitSize()
	testutils.ASSERT_UEQ(t, size, 16)

	scheme.SetBitField("crc", 8)

	size = scheme.BitSize()
	testutils.ASSERT_UEQ(t, size, 24)

	size = 0
	for _, field := range scheme.fields {
		size += field.bitSize()
	}

	testutils.ASSERT_UEQ(t, size, scheme.BitSize())
}

func TestSchemeFieldsOffset(t *testing.T) {
	scheme := NewBitScheme()
	scheme.SetBitField("head", 4)
	scheme.SetBitField("type", 8)
	scheme.SetStuffBits("fill", 4)
	scheme.SetBitField("crc", 8)

	offset, _ := scheme.OffsetOf("type")
	testutils.ASSERT_UEQ(t, offset, 4)

	offset, _ = scheme.OffsetOf("crc")
	testutils.ASSERT_UEQ(t, offset, 16)
}
