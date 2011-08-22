ID3 Parsing For Go
==================

Andrew Scherkus
August 17, 2011


Introduction
------------

Simple ID3 parsing library for go based on the specs at www.id3.org.

It doesn't handle everything but at least gets the imporant bits like artist,
album, track, etc...


Usage
-----
Pass in a suitable io.Reader and away you go!

    fd, _ := os.Open("foo.mp3")
    defer fd.Close()
    file := id3.Read(fd)
    if file != nil {
            fmt.Println(file)
    }

Currently go-id3 will panic if it encounters a frame it does not know about.
This is temporary while I try to fix some bugs.


Examples
--------
An example tag reading program can be found in the examples directory.

    cd examples
    gomake
    ./tagreader path/to/file.mp3 [...]
