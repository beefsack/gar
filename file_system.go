package gar

import (
	"os"
	"path/filepath"
)

type FileSystem struct {
	base string
}

func NewFileSystemSource(base string) FileSystem {
	return FileSystem{
		base: base,
	}
}

func (fs FileSystem) Open(name string) (file File, ok bool, err error) {
	var f *os.File
	f, err = os.Open(filepath.Join(fs.base, name))
	if os.IsNotExist(err) {
		err = nil
		return
	} else if err != nil {
		return
	}
	if file.FileInfo, err = f.Stat(); err != nil {
		return
	}
	ok = true
	file.Content = f
	return
}

func (fs FileSystem) Files() ([]string, error) {
	files := []string{}
	if err := filepath.Walk(
		fs.base,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(fs.base, path)
			if err != nil {
				return err
			}
			files = append(files, rel)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return files, nil
}
