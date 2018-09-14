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
					{"name": "data", "size": 16}
				]
	}
]`

func TestReadingSchemeFromJSON(t *testing.T) {
	var schemes []*BitScheme
	schemes = ReadSchemesFromString(sampleJSON)
	testutils.ASSERT_EQ(t, len(schemes), 2)
	testutils.ASSERT_UEQ(t, schemes[0].BitSize(), 24)
	testutils.ASSERT_UEQ(t, schemes[1].BitSize(), 32)
	size, err := schemes[0].SizeOf("head")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "head")
	}
	testutils.ASSERT_UEQ(t, size, 4)

	size, err = schemes[0].SizeOf("type")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "type")
	}
	testutils.ASSERT_UEQ(t, size, 8)

	size, err = schemes[0].SizeOf("fill")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "fill")
	}
	testutils.ASSERT_UEQ(t, size, 4)

	size, err = schemes[0].SizeOf("crc")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "crc")
	}
	testutils.ASSERT_UEQ(t, size, 8)

	size, err = schemes[1].SizeOf("head")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "head")
	}
	testutils.ASSERT_UEQ(t, size, 4)

	size, err = schemes[1].SizeOf("fill")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "fill")
	}
	testutils.ASSERT_UEQ(t, size, 4)

	size, err = schemes[1].SizeOf("type")
	if err != nil {
		t.Fatalf("This field: |%s| not  found in packet.", "type")
	}
	testutils.ASSERT_UEQ(t, size, 8)

	size, err = schemes[1].SizeOf("data")
	if err != nil {
		t.Fatalf("This field: |%s| not found in packet.", "data")
	}
	testutils.ASSERT_UEQ(t, size, 16)
}
