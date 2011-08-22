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
)

type v22Parser struct {
	reader *bufio.Reader
}

type v22Frame struct {
	id   string
	size int32
}

func newID3v22Parser(reader *bufio.Reader) id3Parser {
	parser := new(v22Parser)
	parser.reader = reader
	return parser
}

func (parser *v22Parser) HasFrame() bool {
	data, err := parser.reader.Peek(3)
	if err != nil {
		return false
	}

	for _, c := range data {
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

func (parser *v22Parser) ReadFrame(file *File) {
	frame := new(v22Frame)

	id := readBytes(parser.reader, 3)
	frame.id = string(id)

	size := readBytes(parser.reader, 3)
	frame.size = int32(size[0])<<16 | int32(size[1])<<8 | int32(size[2])

	switch frame.id {
	case "TAL":
		file.Album = readString(parser.reader, int(frame.size))
	case "TRK":
		file.Track = readString(parser.reader, int(frame.size))
	case "TP1":
		file.Artist = readString(parser.reader, int(frame.size))
	case "TT2":
		file.Name = readString(parser.reader, int(frame.size))
	case "TYE":
		file.Year = readString(parser.reader, int(frame.size))
	case "TPA":
		file.Disc = readString(parser.reader, int(frame.size))
	case "TCO":
		file.Genre = readString(parser.reader, int(frame.size))

	// Skip these frame types.
	case "CM1": // ???
		fallthrough
	case "COM": // Comment
		fallthrough
	case "IPP": // ???
		fallthrough
	case "PCS": // ???
		fallthrough
	case "PIC": // Picture
		fallthrough
	case "TCM": // Composer
		fallthrough
	case "TCT": // ???
		fallthrough
	case "TDR": // ???
		fallthrough
	case "TEN": // Encoded by
		fallthrough
	case "TID": // ???
		fallthrough
	case "TP2": // Band/Orchestra/Accompaniment
		fallthrough
	case "TT1": // Content group description
		fallthrough
	case "UFI": // Unique file identifier
		fallthrough
	case "ULT": // Unsychronized lyrics/text transcription
		fallthrough
	case "WFD": // ???
		skipBytes(parser.reader, int(frame.size))
		break
	default:
		panic(fmt.Sprintf("Unrecognized frame ID: %s, size=%d\n", frame.id, frame.size))
	}
}
