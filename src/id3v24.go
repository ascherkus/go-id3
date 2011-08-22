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

type v24Parser struct {
	reader *bufio.Reader
}

type v24Frame struct {
	id                       string
	size                     int32
	tagAlterPreservation     bool
	fileAlterPreservation    bool
	readOnly                 bool
	containsGroupInformation bool
	compressed               bool
	encrypted                bool
	unsynchronized           bool
	dataLengthIndicator      bool
}

func newID3v24Parser(reader *bufio.Reader) id3Parser {
	parser := new(v24Parser)
	parser.reader = reader
	return parser
}

func (parser *v24Parser) HasFrame() bool {
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

func (parser *v24Parser) ReadFrame(file *File) {
	frame := new(v24Frame)

	id := readBytes(parser.reader, 4)
	frame.id = string(id)

	// ID3 v2.4 uses sync-safe frame sizes.
	size := readBytes(parser.reader, 4)
	frame.size = parseSize(size)

	flags := readBytes(parser.reader, 2)
	frame.tagAlterPreservation = flags[0]&1<<6 != 0
	frame.fileAlterPreservation = flags[0]&1<<5 != 0
	frame.readOnly = flags[0]&1<<4 != 0
	frame.containsGroupInformation = flags[1]&1<<6 != 0
	frame.compressed = flags[1]&1<<3 != 0
	frame.encrypted = flags[1]&1<<2 != 0
	frame.unsynchronized = flags[1]&1<<1 != 0
	frame.dataLengthIndicator = flags[1]&1<<0 != 0

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
	case "TDRC":
		// TODO: implement timestamp parsing
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
	case "MCDI": // Music CD identifier
		fallthrough
	case "PRIV": // Private
		fallthrough
	case "RVA2": // Relative volume adjustment 2 (v2.4 only)
		fallthrough
	case "TCMP": // ???
		fallthrough
	case "TCOM": // Composer
		fallthrough
	case "TCOP": // Copyright message
		fallthrough
	case "TDEN": // Encoding time (v2.4 only)
		fallthrough
	case "TDTG": // Tagging time (v2.4 only)
		fallthrough
	case "TENC": // Encoded by
		fallthrough
	case "TLAN": // Language
		fallthrough
	case "TPE2": // Band/Orchestra/Accompaniment
		fallthrough
	case "TPUB": // Publisher
		fallthrough
	case "TSOP": // Performer sort order (v2.4 only)
		fallthrough
	case "TSRC": // International recording code
		fallthrough
	case "TSSE": // Software/Hardware encoder settings
		fallthrough
	case "TXXX": // User defined text
		fallthrough
	case "UFID": // Unique file identifier
		fallthrough
	case "WXXX": // User defined URL
		skipBytes(parser.reader, int(frame.size))
		break
	default:
		panic(fmt.Sprintf("Unrecognized frame ID: %s, size=%d\n", frame.id, frame.size))
	}
}
