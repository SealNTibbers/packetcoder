package packetcoder

import (
	"github.com/SealNTibbers/packetcoder/testutils"
	"testing"
)

var sampleJSON = `[
	{
		"name": "packet1",
		"fields":
                [
					{"name": "head", "sizeInBits": 4},
					{"name": "type", "sizeInBits": 8},
					{"name": "fill", "sizeInBits": 4},
					{"name": "crc", "sizeInBits": 8}
				]
	},
	{
		"name": "packet2",
		"fields":
                [
					{"name": "head", "sizeInBits": 4},
					{"name": "fill", "sizeInBits": 4},
					{"name": "type", "sizeInBits": 8},
					{"name": "data", "sizeInBits": 16, "littleEndian": true}
				]
	}
]`

func TestReadingSchemeFromJSON(t *testing.T) {
	var schemes map[string]*BitScheme
	schemes = ReadSchemesFromString(sampleJSON)
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
