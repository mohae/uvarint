package varint

import (
	"bytes"
	"testing"
)

var tests = []struct {
	decoded uint64
	n       int
	encoded []byte
}{
	{0, 1, []byte{0x00}},
	{1, 1, []byte{0x01}},
	{240, 1, []byte{0xF0}},
	{241, 2, []byte{0xF1, 0x01}},
	{248, 2, []byte{0xF1, 0x08}},

	{249, 2, []byte{0xF1, 0x09}},
	{2287, 2, []byte{0xF8, 0xFF}},
	{2288, 3, []byte{0xF9, 0x00, 0x00}},
	{67823, 3, []byte{0xF9, 0xFF, 0xFF}},
	{67824, 4, []byte{0xFA, 0x01, 0x08, 0xF0}},

	{16777215, 4, []byte{0xFA, 0xFF, 0xFF, 0xFF}},
	{16777216, 5, []byte{0xFB, 0x01, 0x00, 0x00, 0x00}},
	{4294967295, 5, []byte{0xFB, 0xFF, 0xFF, 0xFF, 0xFF}},
	{4294967296, 6, []byte{0xFC, 0x01, 0x00, 0x00, 0x00, 0x00}},
	{1099511627775, 6, []byte{0xFC, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},

	{1099511627776, 7, []byte{0xFD, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{281474976710655, 7, []byte{0xFD, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{281474976710656, 8, []byte{0xFE, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{72057594037927935, 8, []byte{0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{90000000000000000, 9, []byte{0xFF, 0x01, 0x3F, 0xBE, 0x85, 0xED, 0xC9, 0x00, 0x00}},
}

func TestVarInt(t *testing.T) {
	for i, test := range tests {
		b := make([]byte, len(test.encoded))
		n := Encode(b, test.decoded)
		if n != test.n {
			t.Errorf("encode %d: got %d want %d", i, n, test.n)
		}
		if !bytes.Equal(b, test.encoded) {
			t.Errorf("encode %d: got %v want %v", i, b[0:n], test.encoded)
		}
		v, n := Decode(test.encoded)
		if n != test.n {
			t.Errorf("decode %d: got %d want %d", i, n, test.n)
		}
		if v != test.decoded {
			t.Errorf("decode %d: got %d want %d", i, v, test.decoded)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	buf := make([]byte, 9)
	var n int
	b.SetBytes(8)
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			n = Encode(buf, test.decoded)
		}
	}
	_ = n
}

func BenchmarkDecode(b *testing.B) {
	var res uint64
	b.SetBytes(8)
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			res, _ = Decode(test.encoded)
		}
	}
	_ = res
}
