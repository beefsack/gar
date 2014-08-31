package gar

import (
	"fmt"
	"os"
	"path/filepath"
)

func DefaultLoader() (*Loader, error) {
	l := NewLoader()
	// Use gar file first if executable is one.
	f, err := os.Open(os.Args[0])
	if err != nil {
		return nil, fmt.Errorf("could not open executable, %v", err)
	}
	is, _, err := Stat(f)
	if err != nil {
		return nil, fmt.Errorf("could not stat executable, %v", err)
	}
	if is {
		// Executable is a gar archive.
		r, err := NewReader(f)
		if err != nil {
			return nil, err
		}
		l.AddSource(NewTarSource(r))
	} else if garPath := os.Getenv("GAR_PATH"); garPath != "" {
		// Executable is not a gar archive but a custom gar path provided.
		l.AddSource(NewFileSystemSource(garPath))
	} else {
		// Fallback to use directory of executable and pwd.
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("could not get working directory, %v", err)
		}
		l.AddSource(NewFileSystemSource(filepath.Dir(os.Args[0])))
		l.AddSource(NewFileSystemSource(wd))
	}
	return l, nil
}
