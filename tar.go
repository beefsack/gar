package gar

import (
	"archive/tar"
	"bytes"
	"io"
)

type Tar struct {
	tree   map[string]File
	reader io.Reader
}

func NewTarSource(r io.Reader) *Tar {
	return &Tar{
		reader: r,
	}
}

func (t *Tar) LoadIfRequired() error {
	if t.tree == nil {
		return t.Load()
	}
	return nil
}

func (t *Tar) Load() error {
	t.tree = map[string]File{}
	tr := tar.NewReader(t.reader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		f := File{
			FileInfo: hdr.FileInfo(),
		}
		contents := bytes.NewBuffer([]byte{})
		if _, err = io.Copy(contents, tr); err != nil {
			return err
		}
		f.Content = bytes.NewReader(contents.Bytes())
		t.tree[hdr.Name] = f
	}
	return nil
}

func (t *Tar) Open(name string) (file File, ok bool, err error) {
	if err = t.LoadIfRequired(); err != nil {
		return
	}
	file, ok = t.tree[name]
	return
}

func (t *Tar) Files() ([]string, error) {
	if err := t.LoadIfRequired(); err != nil {
		return nil, err
	}
	files := make([]string, len(t.tree))
	ctr := 0
	for f, _ := range t.tree {
		files[ctr] = f
		ctr += 1
	}
	return files, nil
}
