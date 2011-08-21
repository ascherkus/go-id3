// Copyright 2011 Andrew Scherkus
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package id3

import (
	"bufio"
	"fmt"
	"utf8"
	"utf16"
	"strings"
)

var skipBuffer []byte = make([]byte, 1024*4)


// TODO: this is pretty inefficient but needed since Go handles extended ISO/IEC 8859-1 characters (>= 0x80)
// differently depending on whether you call string() with a []byte or an []int.
func expand(data []byte) []int {
	p := make([]int, len(data))
	for i, b := range data {
		p[i] = int(b)
	}
	return p
}

func toUTF16(data []byte) []uint16 {
	if len(data) < 2 {
		panic("Sequence is too short too contain a UTF-16 BOM")
	}
	if len(data)%2 > 0 {
		// TODO: if this is UTF-16 BE then this is likely encoded wrong
		data = append(data, 0)
	}

	var shift0, shift1 uint
	if data[0] == 0xFF && data[1] == 0xFE {
		// UTF-16 LE
		shift0 = 0
		shift1 = 8
	} else if data[0] == 0xFE && data[1] == 0xFF {
		// UTF-16 BE
		shift0 = 8
		shift1 = 0
		panic("UTF-16 BE found!")
	} else {
		panic(fmt.Sprintf("Unrecognized UTF-16 BOM: 0x%02X%02X", data[0], data[1]))
	}

	s := make([]uint16, 0, len(data)/2)
	for i := 2; i < len(data); i += 2 {
		s = append(s, uint16(data[i])<<shift0|uint16(data[i+1])<<shift1)
	}
	return s
}

// Sizes are stored big endian but with the first bit set to 0 and always ignored.
func parseSize(data []byte) int32 {
	size := int32(0)
	for i, b := range data {
		if b&0x80 > 0 {
			fmt.Println("Size byte had non-zero first bit")
		}

		shift := uint32(len(data)-i-1) * 7
		size |= int32(b&0x7f) << shift
	}
	return size
}

func parseString(data []byte) string {
	var s string
	switch data[0] {
	case 0: // ISO-8859-1 text.
		s = string(expand(data[1:]))
		break
	case 1: // UTF-16 with BOM.
		s = string(utf16.Decode(toUTF16(data[1:])))
		break
	case 2: // UTF-16BE without BOM.
		panic("Unsupported text encoding UTF-16BE.")
	case 3: // UTF-8 text.
		s = utf8.NewString(string(data[1:])).String()
		break
	default:
		// No encoding, assume ISO-8859-1 text.
		s = string(expand(data))
	}
	return strings.TrimRight(s, "\u0000")
}

func readBytes(reader *bufio.Reader, c int) []byte {
	b := make([]byte, c)
	pos := 0
	for pos < c {
		i, err := reader.Read(b[pos:])
		pos += i
		if err != nil {
			panic(err)
		}
	}
	return b
}

func readString(reader *bufio.Reader, c int) string {
	return parseString(readBytes(reader, c))
}

func skipBytes(reader *bufio.Reader, c int) {
	pos := 0
	for pos < c {
		end := c - pos
		if end > len(skipBuffer) {
			end = len(skipBuffer)
		}

		i, err := reader.Read(skipBuffer[0:end])
		pos += i
		if err != nil {
			panic(err)
		}
	}
}
