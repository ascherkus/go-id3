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

package main

import (
	"fmt"
	"id3"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s [path to mp3s]\n", os.Args[0])
		return
	}

	for _, s := range os.Args[1:] {
		var fd, err = os.Open(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open %s: %s\n", s, err)
			return
		}
		file := id3.Read(fd)
		if file == nil {
			fmt.Fprintf(os.Stderr, "Could not read ID3 information from %s\n", s)
		} else {
			fmt.Println(s)
			fmt.Printf("Header\t%s\n", file.Header)
			fmt.Printf("Name\t%s\n", file.Name)
			fmt.Printf("Artist\t%s\n", file.Artist)
			fmt.Printf("Album\t%s\n", file.Album)
			fmt.Printf("Year\t%s\n", file.Year)
			fmt.Printf("Track\t%s\n", file.Track)
			fmt.Printf("Disc\t%s\n", file.Disc)
			fmt.Printf("Genre\t%s\n", file.Genre)
			fmt.Printf("Length\t%s\n", file.Length)
			fmt.Println()
		}
		fd.Close()
	}
}
