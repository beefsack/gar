package gar

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

var FilesNoHiddenOrGo = regexp.MustCompile("")

type Writer struct {
	w io.Writer
	t *tar.Writer
	b *bytes.Buffer
}

func NewWriter(w io.Writer) *Writer {
	wr := &Writer{
		w: w,
		b: bytes.NewBuffer([]byte{}),
	}
	wr.t = tar.NewWriter(wr.b)
	return wr
}

func (w *Writer) WriteFileAtPath(source, target string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	return w.WriteFileWithFileInfo(target, fi, file)
}

func (w *Writer) WriteFile(hdr *tar.Header, contents io.Reader) error {
	if hdr.Size == 0 {
		// Calculate the size in case the file isn't empty but it hasn't been set.
		b, err := ioutil.ReadAll(contents)
		if err != nil {
			return err
		}
		hdr.Size = int64(len(b))
		contents = bytes.NewReader(b)
	}
	if err := w.t.WriteHeader(hdr); err != nil {
		return err
	}
	_, err := io.Copy(w.t, contents)
	return err
}

func (w *Writer) WriteFileWithName(name string, contents io.Reader) error {
	return w.WriteFile(&tar.Header{
		Name: name,
	}, contents)
}

func (w *Writer) WriteFileWithFileInfo(name string, fi os.FileInfo, contents io.Reader) error {
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	hdr.Name = name
	return w.WriteFile(hdr, contents)
}

func (w *Writer) Close(flags byte) error {
	if err := w.t.Close(); err != nil {
		return err
	}
	wr := w.w
	isGzip := flags&FlagGzip != 0
	if isGzip {
		wr = gzip.NewWriter(wr)
	}
	n, err := io.Copy(wr, bytes.NewReader(w.b.Bytes()))
	if err != nil {
		return err
	}
	gi := GarInfo{
		Flags: flags,
		Size:  n,
	}
	return gi.WriteTo(wr)
}
