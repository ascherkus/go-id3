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
	"encoding/binary"
	"bufio"
	"fmt"
)

type v23Parser struct {
	reader *bufio.Reader
}

type v23Frame struct {
	id    string
	size  int32
	flags uint16
}

func NewID3v23Parser(reader *bufio.Reader) id3Parser {
	parser := new(v23Parser)
	parser.reader = reader
	return parser
}

func (parser *v23Parser) HasFrame() bool {
	data, err := parser.reader.Peek(4)
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

func (parser *v23Parser) ReadFrame(file *File) {
	frame := new(v23Frame)

	id := readBytes(parser.reader, 4)
	frame.id = string(id)

	// ID3 v2.3 doesn't use sync-safe frame sizes.
	binary.Read(parser.reader, binary.BigEndian, &frame.size)
	binary.Read(parser.reader, binary.BigEndian, &frame.flags)

	switch frame.id {
	case "TALB":
		file.Album = readString(parser.reader, int(frame.size))
	case "TRCK":
		file.Track = readString(parser.reader, int(frame.size))
	case "TPE1":
		file.Artist = readString(parser.reader, int(frame.size))
	case "TCON":
		// TODO: handle ID3v1 numerical genres.
		file.Genre = readString(parser.reader, int(frame.size))
	case "TIT2":
		file.Name = readString(parser.reader, int(frame.size))
	case "TYER":
		file.Year = readString(parser.reader, int(frame.size))
	case "TPOS":
		file.Disc = readString(parser.reader, int(frame.size))
	case "TLEN":
		file.Length = readString(parser.reader, int(frame.size))

	// Skip these frame types.
	case "APIC": // Picture
		fallthrough
	case "COMM": // Comment
		fallthrough
	case "GEOB": // General encapsulated object
		fallthrough
	case "PRIV": // Private
		fallthrough
	case "MCDI": // Music CD identifier
		fallthrough
	case "NCON": // ???
		fallthrough
	case "POPM": // Popularimeter
		fallthrough
	case "RGAD": // ???
		fallthrough
	case "RVAD": // Relative volume adjustment
		fallthrough
	case "SYLT": // Synchronized lyrics/text transcription
		fallthrough
	case "TBPM": // Beats per minute
		fallthrough
	case "TCMP": // ???
		fallthrough
	case "TCOM": // Composer
		fallthrough
	case "TCOP": // Copyright message
		fallthrough
	case "TDAT": // Date of recording (DDMM)
		fallthrough
	case "TDEN": // ???
		fallthrough
	case "TDTG": // ???
		fallthrough
	case "TENC": // Encoded by
		fallthrough
	case "TFLT": // File type
		fallthrough
	case "TIT1": // Content group description
		fallthrough
	case "TIT3": // Subtitle/Description refinement
		fallthrough
	case "TLAN": // Language
		fallthrough
	case "TMED": // Media type
		fallthrough
	case "TOPE": // Original artist
		fallthrough
	case "TORY": // Original release year
		fallthrough
	case "TPE2": // Band/Orchestra/Accompaniment
		fallthrough
	case "TPE3": // Conductor
		fallthrough
	case "TPUB": // Publisher
		fallthrough
	case "TSIZ": // Size
		fallthrough
	case "TSRC": // International recording code
		fallthrough
	case "TSSE": // Software/Hardware encoder settings
		fallthrough
	case "TXXX": // User defined text
		fallthrough
	case "USER": // Terms of use
		fallthrough
	case "UFID": // Unique file identifier
		fallthrough
	case "USLT": // Unsynchronized lyrics/text transcription
		fallthrough
	case "XSOP": // ???
		fallthrough
	case "WCOM": // Commercial information
		fallthrough
	case "WOAF": // Official audio file webpage
		fallthrough
	case "WXXX": // User defined URL
		skipBytes(parser.reader, int(frame.size))
		break
	default:
		panic(fmt.Sprintf("Unrecognized frame ID: %s, size=%d\n", frame.id, frame.size))
	}
}
