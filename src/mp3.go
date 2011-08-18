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
	"io"
	"fmt"
)

type File struct {
	ID3v2  *ID3v2Header
	Name   string
	Artist string
	Album  string
	Year   string
	Track  string
	Disc   string
	Genre  string
	Length string
}

func Read(reader io.Reader) *File {
	file := new(File)
	bufReader := bufio.NewReader(reader)
	if !isID3Tag(bufReader) {
		return nil
	}

	file.ID3v2 = parseID3v2Header(bufReader)
	limitReader := bufio.NewReader(io.LimitReader(bufReader, int64(file.ID3v2.size)))
	var parser ID3Parser
	if file.ID3v2.version == 2 {
		parser = NewID3v22Parser(limitReader)
	} else if file.ID3v2.version == 3 {
		parser = NewID3v23Parser(limitReader)
	} else if file.ID3v2.version == 4 {
		parser = NewID3v24Parser(limitReader)
	} else {
		panic(fmt.Sprintf("Unrecognized ID3v2 version: %d", file.ID3v2.version))
	}

	for parser.HasFrame() {
		parser.ReadFrame(file)
	}

	return file
}
