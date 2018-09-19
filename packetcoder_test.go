package packetcoder

import (
	"bytes"
	"github.com/SealNTibbers/packetcoder/testutils"
	"testing"
)

func TestPacketEncode(t *testing.T) {
	packet := NewPacket()
	scheme := NewBitScheme("test")
	scheme.AddBitField("head", 4)
	scheme.AddBitField("type", 8)
	scheme.AddStuffBits("fill", 4)
	scheme.AddBitField("crc", 8)
	packet.SetScheme(scheme)

	packet.WriteValue64("head", 5)
	packet.WriteValue64("type", 105)
	packet.WriteStuff("fill")
	packet.WriteValue64("crc", 99)
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
	scheme.AddBitField("head", 4)
	scheme.AddBitField("type", 8)
	scheme.AddStuffBits("fill", 4)
	scheme.AddBitField("crc", 8)
	packet.SetScheme(scheme)
	buf := bytes.NewBuffer([]byte{86, 144, 99, 86, 144, 140, 86, 144})
	// decode first packet
	_, err := packet.DecodeFrom(buf)

	value, _ := packet.ReadValue64("head")
	testutils.ASSERT_U64_EQ(t, value, 5)

	value, _ = packet.ReadValue64("type")
	testutils.ASSERT_U64_EQ(t, value, 105)

	value, _ = packet.ReadValue64("crc")
	testutils.ASSERT_U64_EQ(t, value, 99)

	// decode second packet
	_, err = packet.DecodeFrom(buf)

	value, _ = packet.ReadValue64("head")
	testutils.ASSERT_U64_EQ(t, value, 5)

	value, _ = packet.ReadValue64("type")
	testutils.ASSERT_U64_EQ(t, value, 105)

	value, _ = packet.ReadValue64("crc")
	testutils.ASSERT_U64_EQ(t, value, 140)

	//decode remainder and we should get error, because sizeInBits of packet more sizeInBits of buffer
	_, err = packet.DecodeFrom(buf)
	testutils.ASSERT_TRUE(t, err != nil)
}

func GetTestBitScheme() *BitScheme {
	scheme := NewBitScheme("testBits")
	scheme.AddBitField("head", 4)
	scheme.AddBitField("type", 8)
	scheme.AddByteField("data", 10)
	scheme.AddStuffBits("fill", 4)
	return scheme
}

func TestSchemeBitSize(t *testing.T) {
	scheme := GetTestBitScheme()

	var size uint
	size = scheme.BitSize()
	testutils.ASSERT_U_EQ(t, size, 96)

	scheme.AddBitField("crc", 8)

	size = scheme.BitSize()
	testutils.ASSERT_U_EQ(t, size, 104)

	size = 0
	for _, field := range scheme.fields {
		size += field.bitSize()
	}

	testutils.ASSERT_U_EQ(t, size, scheme.BitSize())
}

func TestBitAndByteSize(t *testing.T) {
	scheme := GetTestBitScheme()

	var size uint
	size, _ = scheme.BitSizeOf("data")
	testutils.ASSERT_U_EQ(t, size, 80)

	size, _ = scheme.ByteSizeOf("data")
	testutils.ASSERT_U_EQ(t, size, 10)
}

func TestSchemeFieldsOffset(t *testing.T) {
	scheme := GetTestBitScheme()
	scheme.AddBitField("crc", 8)

	offset, _ := scheme.OffsetOf("type")
	testutils.ASSERT_U_EQ(t, offset, 4)

	offset, _ = scheme.OffsetOf("data")
	testutils.ASSERT_U_EQ(t, offset, 12)

	offset, _ = scheme.OffsetOf("crc")
	testutils.ASSERT_U_EQ(t, offset, 96)
}

func TestPacketFieldValue64(t *testing.T) {
	packet := NewPacket()
	scheme := NewBitScheme("test")
	scheme.AddBitField("head", 4)
	scheme.AddBitField("type", 8)
	scheme.AddStuffBits("fill", 4)
	scheme.AddBitField("bigEndian", 32)
	scheme.AddBitFieldLittleEndian("littleEndian", 32)
	scheme.AddBitField("crc", 8)
	packet.SetScheme(scheme)

	packet.WriteValue64("head", 5)
	packet.WriteValue64("type", 105)
	packet.WriteStuff("fill")
	packet.WriteValue64("bigEndian", 0x01234567)
	packet.WriteValue64("littleEndian", 0x01234567)
	packet.WriteValue64("crc", 99)

	//test big endian
	bytes := packet.GetData().Bytes()
	testutils.ASSERT_BYTE_EQ(t, bytes[2], 0x01)
	testutils.ASSERT_BYTE_EQ(t, bytes[3], 0x23)
	testutils.ASSERT_BYTE_EQ(t, bytes[4], 0x45)
	testutils.ASSERT_BYTE_EQ(t, bytes[5], 0x67)

	//test little endian
	testutils.ASSERT_BYTE_EQ(t, bytes[6], 0x67)
	testutils.ASSERT_BYTE_EQ(t, bytes[7], 0x45)
	testutils.ASSERT_BYTE_EQ(t, bytes[8], 0x23)
	testutils.ASSERT_BYTE_EQ(t, bytes[9], 0x01)

	value, _ := packet.ReadValue64("head")
	testutils.ASSERT_U64_EQ(t, value, 5)

	value, _ = packet.ReadValue64("type")
	testutils.ASSERT_U64_EQ(t, value, 105)

	value, _ = packet.ReadValue64("bigEndian")
	testutils.ASSERT_U64_EQ(t, value, 0x01234567)

	value, _ = packet.ReadValue64("littleEndian")
	testutils.ASSERT_U64_EQ(t, value, 0x01234567)

	value, _ = packet.ReadValue64("crc")
	testutils.ASSERT_U64_EQ(t, value, 99)
}

func TestPacketFieldBytesValue(t *testing.T) {
	packet := NewPacket()
	scheme := GetTestBitScheme()
	scheme.AddByteFieldLittleEndian("littleEndian", 5)
	scheme.AddBitField("crc", 8)
	packet.SetScheme(scheme)

	packet.WriteValue64("head", 0x05)
	packet.WriteValue64("type", 0x69)
	packet.WriteBytes("data", []byte{0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78, 0x89, 0xAB})
	packet.WriteBytes("littleEndian", []byte{0x56, 0x67, 0x78, 0x89, 0xAB})
	packet.WriteStuff("fill")
	packet.WriteValue64("crc", 99)

	bytes := packet.GetData().Bytes()
	testutils.ASSERT_BYTE_EQ(t, bytes[0], 0x56)
	testutils.ASSERT_BYTE_EQ(t, bytes[1], 0x90)
	testutils.ASSERT_BYTE_EQ(t, bytes[2], 0x11)
	testutils.ASSERT_BYTE_EQ(t, bytes[3], 0x22)
	testutils.ASSERT_BYTE_EQ(t, bytes[10], 0x9A)

	testutils.ASSERT_BYTE_EQ(t, bytes[11], 0xBA)
	testutils.ASSERT_BYTE_EQ(t, bytes[12], 0xB8)
	testutils.ASSERT_BYTE_EQ(t, bytes[13], 0x97)
	testutils.ASSERT_BYTE_EQ(t, bytes[14], 0x86)

	value, _ := packet.ReadValue64("head")
	testutils.ASSERT_U64_EQ(t, value, 5)

	value, _ = packet.ReadValue64("type")
	testutils.ASSERT_U64_EQ(t, value, 105)

	bytesValue, _ := packet.ReadBytesValue("data")
	testutils.ASSERT_U_EQ(t, uint(bytesValue[0]), 0x01)
	testutils.ASSERT_U_EQ(t, uint(bytesValue[9]), 0xAB)

	value, _ = packet.ReadValue64("crc")
	testutils.ASSERT_U64_EQ(t, value, 99)
}
