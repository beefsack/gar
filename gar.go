package gar

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	FlagGzip = 1 << iota
	FlagExtractFileSystem
)

const (
	GarIdentifier      = "GAR"
	GarFileSizeOptSize = 8
	GarFlagsSize       = 1
	GarOptsSize        = GarFileSizeOptSize + GarFlagsSize
	GarFooterSize      = len(GarIdentifier) + GarOptsSize
)

var globalLoader *Loader

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

type File struct {
	FileInfo os.FileInfo
	Content  ReadSeekCloser
}

type GarInfo struct {
	Flags byte
	Size  int64
}

func (g GarInfo) IsFlag(flag byte) bool {
	return g.Flags&flag != 0
}

func (g *GarInfo) SetFlag(flag byte, value bool) {
	if value {
		g.Flags |= flag
	} else {
		g.Flags &= ^flag
	}
}

func (gi GarInfo) WriteTo(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, gi.Size); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, gi.Flags); err != nil {
		return err
	}
	if _, err := w.Write([]byte(GarIdentifier)); err != nil {
		return err
	}
	return nil
}

func Stat(r io.ReadSeeker) (is bool, info GarInfo, err error) {
	var (
		footer []byte
	)
	if _, err = r.Seek(int64(-GarFooterSize), 2); err != nil {
		err = fmt.Errorf("could not seek to start of gar footer, %v", err)
		return
	}
	if footer, err = ioutil.ReadAll(r); err != nil {
		err = fmt.Errorf("could not read gar footer, %v", err)
		return
	}
	if len(footer) != GarFooterSize {
		err = fmt.Errorf("footer size did not match expected %d", GarFooterSize)
		return
	}
	is = string(footer[GarFooterSize-len(GarIdentifier):]) == GarIdentifier
	if !is {
		return
	}
	br := bytes.NewReader(footer)
	if err = binary.Read(br, binary.LittleEndian, &info.Size); err != nil {
		err = fmt.Errorf("could not read size value from footer, %v", err)
	}
	if err = binary.Read(br, binary.LittleEndian, &info.Flags); err != nil {
		err = fmt.Errorf("could not read gzip flag from footer, %v", err)
	}
	return
}

func NewGarSource(r io.ReadSeeker) (tar *Tar, err error) {
	g, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	return NewTarSource(g), nil
}

func GlobalLoader() (*Loader, error) {
	if globalLoader == nil {
		l, err := DefaultLoader()
		if err != nil {
			return nil, fmt.Errorf("could not get default loader, %v", err)
		}
		SetGlobalLoader(l)
	}
	return globalLoader, nil
}

func SetGlobalLoader(loader *Loader) {
	globalLoader = loader
}

func Open(name string) (file File, ok bool, err error) {
	l, err := GlobalLoader()
	if err != nil {
		return
	}
	return l.Open(name)
}

func Files() (files []string, err error) {
	l, err := GlobalLoader()
	if err != nil {
		return
	}
	return l.Files()
}
