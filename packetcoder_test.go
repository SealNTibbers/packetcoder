package packetcoder

import (
	"bytes"
	"github.com/SealNTibbers/packetcoder/testutils"
	"testing"
)

func TestPacketEncode(t *testing.T) {
	packet := NewPacket()
	scheme := NewBitScheme("test")
	scheme.SetBitField("head", 4)
	scheme.SetBitField("type", 8)
	scheme.SetStuffBits("fill", 4)
	scheme.SetBitField("crc", 8)
	packet.SetScheme(scheme)

	packet.WriteValue("head", 5)
	packet.WriteValue("type", 105)
	packet.WriteStuff("fill")
	packet.WriteValue("crc", 99)
	buf := bytes.NewBuffer(nil)
	packet.EncodeTo(buf)
	testutils.ASSERT_EQ(t, len(buf.Bytes()), 3)
	testutils.ASSERT_U_EQ(t, uint(buf.Bytes()[0]), 86)
	testutils.ASSERT_U_EQ(t, uint(buf.Bytes()[1]), 144)
	testutils.ASSERT_U_EQ(t, uint(buf.Bytes()[2]), 99)
}

func TestPacketDecode(t *testing.T) {
	packet := new(Packet)
	scheme := NewBitScheme("test")
	scheme.SetBitField("head", 4)
	scheme.SetBitField("type", 8)
	scheme.SetStuffBits("fill", 4)
	scheme.SetBitField("crc", 8)
	packet.SetScheme(scheme)
	buf := bytes.NewBuffer([]byte{86, 144, 99, 86, 144, 140, 86, 144})
	// decode first packet
	_, err := packet.DecodeFrom(buf)

	value, _ := packet.ReadValue("head")
	testutils.ASSERT_U64_EQ(t, value, 5)

	value, _ = packet.ReadValue("type")
	testutils.ASSERT_U64_EQ(t, value, 105)

	value, _ = packet.ReadValue("crc")
	testutils.ASSERT_U64_EQ(t, value, 99)

	// decode second packet
	_, err = packet.DecodeFrom(buf)

	value, _ = packet.ReadValue("head")
	testutils.ASSERT_U64_EQ(t, value, 5)

	value, _ = packet.ReadValue("type")
	testutils.ASSERT_U64_EQ(t, value, 105)

	value, _ = packet.ReadValue("crc")
	testutils.ASSERT_U64_EQ(t, value, 140)

	//decode remainder and we should get error, because size of packet more size of buffer
	_, err = packet.DecodeFrom(buf)
	testutils.ASSERT_TRUE(t, err != nil)
}

func TestSchemeSetFields(t *testing.T) {
	scheme := NewBitScheme("test")
	scheme.SetBitField("head", 4)
	scheme.SetBitField("type", 8)
	scheme.SetStuffBits("fill", 4)
}

func TestSchemeBitSize(t *testing.T) {
	scheme := NewBitScheme("test")
	scheme.SetBitField("head", 4)
	scheme.SetBitField("type", 8)
	scheme.SetStuffBits("fill", 4)

	var size uint
	size = scheme.BitSize()
	testutils.ASSERT_U_EQ(t, size, 16)

	scheme.SetBitField("crc", 8)

	size = scheme.BitSize()
	testutils.ASSERT_U_EQ(t, size, 24)

	size = 0
	for _, field := range scheme.fields {
		size += field.bitSize()
	}

	testutils.ASSERT_U_EQ(t, size, scheme.BitSize())
}

func TestSchemeFieldsOffset(t *testing.T) {
	scheme := NewBitScheme("test")
	scheme.SetBitField("head", 4)
	scheme.SetBitField("type", 8)
	scheme.SetStuffBits("fill", 4)
	scheme.SetBitField("crc", 8)

	offset, _ := scheme.OffsetOf("type")
	testutils.ASSERT_U_EQ(t, offset, 4)

	offset, _ = scheme.OffsetOf("crc")
	testutils.ASSERT_U_EQ(t, offset, 16)
}

func TestPacketFieldValues(t *testing.T) {
	packet := NewPacket()
	scheme := NewBitScheme("test")
	scheme.SetBitField("head", 4)
	scheme.SetBitField("type", 8)
	scheme.SetStuffBits("fill", 4)
	scheme.SetBitField("bigEndian", 32)
	scheme.SetBitFieldLittleEndian("littleEndian", 32)
	scheme.SetBitField("crc", 8)
	packet.SetScheme(scheme)

	packet.WriteValue("head", 5)
	packet.WriteValue("type", 105)
	packet.WriteStuff("fill")
	packet.WriteValue("bigEndian", 0x01234567)
	packet.WriteValue("littleEndian", 0x01234567)
	packet.WriteValue("crc", 99)

	/*field, _ := packet.Scheme.GetField("bigEndian")
	if field.IsLittleEndian() {
		t.Fatalf("this field should be big endian for now")
	}
	field.SetLittleEndian(true)
	if !field.IsLittleEndian() {
		t.Fatalf("this field should be little endian now")
	}*/

	//test big endian
	bytes := packet.GetData().Bytes()
	testutils.ASSERT_U_EQ(t, uint(bytes[2]), 0x01)
	testutils.ASSERT_U_EQ(t, uint(bytes[3]), 0x23)
	testutils.ASSERT_U_EQ(t, uint(bytes[4]), 0x45)
	testutils.ASSERT_U_EQ(t, uint(bytes[5]), 0x67)

	//test little endian
	testutils.ASSERT_U_EQ(t, uint(bytes[6]), 0x67)
	testutils.ASSERT_U_EQ(t, uint(bytes[7]), 0x45)
	testutils.ASSERT_U_EQ(t, uint(bytes[8]), 0x23)
	testutils.ASSERT_U_EQ(t, uint(bytes[9]), 0x01)

	value, _ := packet.ReadValue("head")
	testutils.ASSERT_U64_EQ(t, value, 5)

	value, _ = packet.ReadValue("type")
	testutils.ASSERT_U64_EQ(t, value, 105)

	value, _ = packet.ReadValue("bigEndian")
	testutils.ASSERT_U64_EQ(t, value, 0x01234567)

	value, _ = packet.ReadValue("littleEndian")
	testutils.ASSERT_U64_EQ(t, value, 0x01234567)

	value, _ = packet.ReadValue("crc")
	testutils.ASSERT_U64_EQ(t, value, 99)
}
