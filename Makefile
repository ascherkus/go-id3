include $(GOROOT)/src/Make.inc

TARG=id3
GOFILES=\
  src/id3.go\
  src/id3v22.go\
  src/id3v23.go\
  src/id3v24.go\
  src/mp3.go\
  src/mpeg.go\

include $(GOROOT)/src/Make.pkg
