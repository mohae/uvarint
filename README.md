# varint
[![Build Status](https://travis-ci.org/mohae/varint.png)](https://travis-ci.org/mohae/varint)  

Varint is SQLite4's variable-length integers: an encoding of 64-bit unsigned integers into between 1 and 9 bytes.

The Encode and Decode functions also document the rules for encoding/decoding as documented on http://www.sqlite.org/src4/doc/trunk/www/varint.wiki.  The rules are copied directly from that page.

For more information about SQLite4's varint, please refer to the above link.

When encoding, the passed byte slice must be have a len of _n_ bytes where _n_ is the number of bytes that the number will encode to, up to a max of 9 bytes.  The number to encode must be of type `uint64`.  The number of bytes encoded will be returned.

Decode takes a byte slice and returns the decoded number as a `uint64` and the number of bytes read.

## Usage

    package main

    import "github.com/mohae/varint"

    func main() {
            buf := make([]byte, 9)
            n := varint.Encode(buf, 42)
            fmt.Printf("Encoded %d bytes: %#v\n", buf)

            v, n := varintDecode(buf)
            fmt.Printf("Decoded %d bytes: %d\n", n, v)
    }

Output:

    Encoded 1 bytes: 0x2a
    Decoded 1 bytes: 42
