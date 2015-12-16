// Package varint encodes and decodes SQLite4 variable length integers
// which is an encoding of 64-bit unsigned integers into 1-9 bytes.
// See: http://www.sqlite.org/src4/doc/trunk/www/varint.wiki
package varint

// Decode decodes the received varint encoded byte slice; returning
// the value and the amount of bytes used.
//
// From: http://www.sqlite.org/src4/doc/trunk/www/varint.wiki
// Decode
//
// If A0 is between 0 and 240 inclusive, then the result is the value of A0.
// If A0 is between 241 and 248 inclusive, then the result is
//    240+256*(A0-241)+A1.
// If A0 is 249 then the result is 2288+256*A1+A2.
// If A0 is 250 then the result is A1..A3 as a 3-byte big-ending integer.
// If A0 is 251 then the result is A1..A4 as a 4-byte big-ending integer.
// If A0 is 252 then the result is A1..A5 as a 5-byte big-ending integer.
// If A0 is 253 then the result is A1..A6 as a 6-byte big-ending integer.
// If A0 is 254 then the result is A1..A7 as a 7-byte big-ending integer.
// If A0 is 255 then the result is A1..A8 as a 8-byte big-ending integer.
func Decode(buf []byte) (uint64, int) {
	// check the first byte
	if buf[0] >= 0 && buf[0] <= 0xF0 {
		return uint64(buf[0]), 1
	}
	if buf[0] >= 0xF1 && buf[0] <= 0xF8 {
		return 240 + 256*(uint64(buf[0])-241)+ uint64(buf[1]), 2
	}
	if buf[0] == 0xF9 {
		return 2288+256*uint64(buf[1]) + uint64(buf[2]), 3
	}
	if buf[0] == 0xFA {
		return bigEndianToUint64(buf[1:4]), 4
	}
	if buf[0] == 0xFB {
		return bigEndianToUint64(buf[1:5]), 5
	}
	if buf[0] == 0xFC {
		return bigEndianToUint64(buf[1:6]), 6
	}
	if buf[0] == 0xFD {
		return bigEndianToUint64(buf[1:7]), 7
	}
	if buf[0] == 0xFE {
		return bigEndianToUint64(buf[1:8]), 8
	}
	if buf[0] == 0xFF {
		return bigEndianToUint64(buf[1:9]), 9
	}
	// panic here
	panic("decode: invalid varint")
}


// Encode encodes the received uint64 into varint using the minimum
// necessary bytes.  The number of bytes written is returned.
//
// From: http://www.sqlite.org/src4/doc/trunk/www/varint.wiki
// Encode
//
// Let the input value be V.
//
// If V<=240 then output a single by A0 equal to V.
// If V<=2287 then output A0 as (V-240)/256 + 241 and A1 as (V-240)%256.
// If V<=67823 then output A0 as 249, A1 as (V-2288)/256, and A2
//    as (V-2288)%256.
// If V<=16777215 then output A0 as 250 and A1 through A3 as a big-endian
//    3-byte integer.
// If V<=4294967295 then output A0 as 251 and A1..A4 as a big-ending
//    4-byte integer.
// If V<=1099511627775 then output A0 as 252 and A1..A5 as a big-ending
//    5-byte integer.
// If V<=281474976710655 then output A0 as 253 and A1..A6 as a big-ending
//    6-byte integer.
// If V<=72057594037927935 then output A0 as 254 and A1..A7 as a big-ending
//    7-byte integer.
// Otherwise then output A0 as 255 and A1..A8 as a big-ending 8-byte integer.
func Encode(buf []byte, x uint64) int {
	if x <= 240 {
		buf[0] = byte(x)
		return 1
	}
	if x <= 2287 {
		buf[0] = byte((x-240)/256+241)
		buf[1] = byte((x-240)%256)
		return 2
	}
	if x <= 67823 {
		buf[0] = 0xF9
		buf[1] = byte((x-2288)/256)
		buf[2] = byte((x-2288)%256)
		return 3
	}
	if x <= 16777215 {
		buf[0] = 0xFA
		uint64ToBigEndian(buf[1:], x, 3)
		return 4
	}
	if x <= 4294967295 {
		buf[0] = 0xFB
		uint64ToBigEndian(buf[1:], x, 4)
		return 5
	}
	if x <= 1099511627775 {
		buf[0] = 0xFC
		uint64ToBigEndian(buf[1:], x, 5)
		return 6
	}
	if x <= 281474976710655 {
		buf[0] = 0xFD
		uint64ToBigEndian(buf[1:], x, 6)
		return 7
	}
	if x <= 72057594037927935 {
		buf[0] = 0xFE
		uint64ToBigEndian(buf[1:], x, 7)
		return 8
	}
	buf[0] = 0xFF
	uint64ToBigEndian(buf[1:], x, 8)
	return 9
}

// uint64ToBigEndian fills buf with x as a n byte big-endian integer
func uint64ToBigEndian(buf []byte, x uint64, n int) {
	for i := 0; i < n; i++ {
		buf[i] = byte(x >> uint(8*(n-(i+1))))
	}
}

// fromBigEndian takeks the contents of the buff and returns
func bigEndianToUint64(buf []byte) uint64 {
	var x uint64
	for i := 0; i < len(buf); i++ {
		x |= uint64(buf[i]) << uint(8*(len(buf)-(i+1)))
	}
	return x
}
