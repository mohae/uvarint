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
	{2287, 2, []byte{0xF8, 0xFF}},

	{2288, 3, []byte{0xF9, 0x00, 0x00}},
	{67823, 3, []byte{0xF9, 0xFF, 0xFF}},
	{67824, 4, []byte{0xFA, 0x01, 0x08, 0xF0}},
	{1<<24 - 1, 4, []byte{0xFA, 0xFF, 0xFF, 0xFF}},
	{1 << 24, 5, []byte{0xFB, 0x01, 0x00, 0x00, 0x00}},

	{1<<32 - 1, 5, []byte{0xFB, 0xFF, 0xFF, 0xFF, 0xFF}},
	{1 << 32, 6, []byte{0xFC, 0x01, 0x00, 0x00, 0x00, 0x00}},
	{1<<40 - 1, 6, []byte{0xFC, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{1 << 40, 7, []byte{0xFD, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1<<48 - 1, 7, []byte{0xFD, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},

	{1 << 48, 8, []byte{0xFE, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1<<56 - 1, 8, []byte{0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{1 << 56, 9, []byte{0xFF, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{1<<64 - 1, 9, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

func TestUvarint(t *testing.T) {
	for i, test := range tests {
		b := make([]byte, len(test.encoded))
		n := PutUvarint(b, test.decoded)
		if n != test.n {
			t.Errorf("encode %d: got %d want %d", i, n, test.n)
		}
		if !bytes.Equal(b, test.encoded) {
			t.Errorf("encode %d: got %v want %v", i, b[0:n], test.encoded)
		}
		v, n := Uvarint(test.encoded)
		if n != test.n {
			t.Errorf("decode %d: got %d want %d", i, n, test.n)
		}
		if v != test.decoded {
			t.Errorf("decode %d: got %d want %d", i, v, test.decoded)
		}
	}
}

// Benchmark using all the tests
func BenchmarkPutUvarintAll(b *testing.B) {
	buf := make([]byte, 9)
	var n int
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			n = PutUvarint(buf, test.decoded)
		}
	}
	_ = n
}

func BenchmarkUvarintAll(b *testing.B) {
	var res uint64
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			res, _ = Uvarint(test.encoded)
		}
	}
	_ = res
}

// Benchmark using a value < 241
func BenchmarkPutUvarintMin(b *testing.B) {
	buf := make([]byte, 9)
	var n int
	for i := 0; i < b.N; i++ {
		n = PutUvarint(buf, tests[2].decoded)
	}
	_ = n
}

func BenchmarkUvarintMin(b *testing.B) {
	var res uint64
	for i := 0; i < b.N; i++ {
		res, _ = Uvarint(tests[2].encoded)
	}
	_ = res
}

// Benchmark using a avlue > 1<<56
func BenchmarkPutUvarintMax(b *testing.B) {
	buf := make([]byte, 9)
	var n int
	for i := 0; i < b.N; i++ {
		n = PutUvarint(buf, tests[17].decoded)
	}
	_ = n
}

func BenchmarkUvarintMax(b *testing.B) {
	var res uint64
	for i := 0; i < b.N; i++ {
		res, _ = Uvarint(tests[17].encoded)
	}
	_ = res
}
