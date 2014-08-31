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
	GarIdentifier      = "GAR"
	GarOptFileSizeSize = 8
	GarOptGzipSize     = 1
	GarOptSize         = GarOptFileSizeSize + GarOptGzipSize
	GarHeaderSize      = len(GarIdentifier) + GarOptSize
)

var globalLoader *Loader

type File struct {
	FileInfo os.FileInfo
	Content  io.ReadSeeker
}

type GarInfo struct {
	Gzip bool
	Size int64
}

func (gi GarInfo) WriteTo(w io.Writer) error {
	var gzipRaw byte
	if err := binary.Write(w, binary.LittleEndian, gi.Size); err != nil {
		return err
	}
	if gi.Gzip {
		gzipRaw = 1
	}
	if err := binary.Write(w, binary.LittleEndian, gzipRaw); err != nil {
		return err
	}
	if _, err := w.Write([]byte(GarIdentifier)); err != nil {
		return err
	}
	return nil
}

func Stat(r io.ReadSeeker) (is bool, info GarInfo, err error) {
	var (
		header  []byte
		gzipRaw byte
	)
	if _, err = r.Seek(int64(-GarHeaderSize), 2); err != nil {
		err = fmt.Errorf("could not seek to start of gar header, %v", err)
		return
	}
	if header, err = ioutil.ReadAll(r); err != nil {
		err = fmt.Errorf("could not read gar header, %v", err)
		return
	}
	if len(header) != GarHeaderSize {
		err = fmt.Errorf("header size did not match expected %d", GarHeaderSize)
		return
	}
	is = string(header[GarHeaderSize-len(GarIdentifier):]) == GarIdentifier
	if !is {
		return
	}
	br := bytes.NewReader(header)
	if err = binary.Read(br, binary.LittleEndian, &info.Size); err != nil {
		err = fmt.Errorf("could not read size value from header, %v", err)
	}
	if err = binary.Read(br, binary.LittleEndian, &gzipRaw); err != nil {
		err = fmt.Errorf("could not read gzip flag from header, %v", err)
	}
	info.Gzip = gzipRaw != 0
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
