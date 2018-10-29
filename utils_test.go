package packetcoder

import (
	"github.com/SealNTibbers/packetcoder/testutils"
	"testing"
)

var testArrayOfJSONPackets = `[
	{
		"name": "packet1",
		"fields":
                [
					{"name": "head", "size": 4},
					{"name": "type", "size": 8},
					{"name": "fill", "size": 4},
					{"name": "crc", "size": 8}
				]
	},
	{
		"name": "packet2",
		"fields":
                [
					{"name": "head", "size": 4},
					{"name": "fill", "size": 4},
					{"name": "type", "size": 8},
					{"name": "data", "size": 16, "littleEndian": true}
				]
	}
]`

var testJSONPacket = `
	{
		"name": "packet1",
		"fields":
                [
					{"name": "head", "size": 4},
					{"name": "type", "size": 8},
					{"name": "fill", "size": 4},
					{"name": "crc", "size": 8}
				]
	}`

func TestReadingSchemeFromJSON(t *testing.T) {
	var scheme *BitScheme
	scheme = ReadSchemeFromString(testJSONPacket)
	size, err := scheme.BitSizeOf("head")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "head")
	}
	testutils.ASSERT_U_EQ(t, size, 4)

	size, err = scheme.BitSizeOf("type")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "type")
	}
	testutils.ASSERT_U_EQ(t, size, 8)

	size, err = scheme.BitSizeOf("fill")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "fill")
	}
	testutils.ASSERT_U_EQ(t, size, 4)

	size, err = scheme.BitSizeOf("crc")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "crc")
	}
	testutils.ASSERT_U_EQ(t, size, 8)
}

func TestReadingSchemesFromJSONArray(t *testing.T) {
	var schemes map[string]*BitScheme
	schemes = ReadSchemesFromString(testArrayOfJSONPackets)
	testutils.ASSERT_EQ(t, len(schemes), 2)
	testutils.ASSERT_U_EQ(t, schemes["packet1"].BitSize(), 24)
	testutils.ASSERT_U_EQ(t, schemes["packet2"].BitSize(), 32)
	size, err := schemes["packet1"].BitSizeOf("head")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "head")
	}
	testutils.ASSERT_U_EQ(t, size, 4)

	size, err = schemes["packet1"].BitSizeOf("type")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "type")
	}
	testutils.ASSERT_U_EQ(t, size, 8)

	size, err = schemes["packet1"].BitSizeOf("fill")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "fill")
	}
	testutils.ASSERT_U_EQ(t, size, 4)

	size, err = schemes["packet1"].BitSizeOf("crc")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "crc")
	}
	testutils.ASSERT_U_EQ(t, size, 8)

	size, err = schemes["packet2"].BitSizeOf("head")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "head")
	}
	testutils.ASSERT_U_EQ(t, size, 4)

	size, err = schemes["packet2"].BitSizeOf("fill")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "fill")
	}
	testutils.ASSERT_U_EQ(t, size, 4)

	size, err = schemes["packet2"].BitSizeOf("type")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "type")
	}
	testutils.ASSERT_U_EQ(t, size, 8)

	size, err = schemes["packet2"].BitSizeOf("data")
	if err != nil {
		t.Fatalf("This field: |%s| not found in packet.", "data")
	}
	testutils.ASSERT_U_EQ(t, size, 16)
}
