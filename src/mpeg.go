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

type Version int

const (
	Version1 Version = 0
	Version2
	Version2_5
)

type Layer int

const (
	LayerI Layer = 0
	LayerII
	LayerIII
)

type ChannelMode int

const (
	Mono ChannelMode = 0
	DualMono
	Stereo
	JointStereo
)

type FrameHeader struct {
	version     Version     // MPEG audio version ID.
	layer       Layer       // MPEG layer description.
	protected   bool        // Protected by CRC.
	bitrate     int         // Bitrate in kbps.
	sampleRate  int         // Sample rate.
	padded      bool        // Frame is padded.
	channelMode ChannelMode // Mono, stereo, joint stereo, etc...
}
