package gar

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestReader_Read(t *testing.T) {
	// Files used for testing.
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling licence."},
	}
	// Create some random file content.
	buf := bytes.NewBuffer([]byte("blah blah blah blah"))
	tarBuf := bytes.NewBuffer([]byte{})
	// Write files as tar to buffer.
	tw := tar.NewWriter(tarBuf)
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	n, err := tarBuf.WriteTo(buf)
	if err != nil {
		t.Fatal(err)
	}
	// Write the gar header
	gi := GarInfo{
		Gzip: false,
		Size: n,
	}
	if err := gi.WriteTo(buf); err != nil {
		t.Fatal(err)
	}
	// Now read the gar back out.
	r, err := NewReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(r)
	ctr := 0
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		if ctr > len(files) {
			t.Fatal("ctr incremented higher than the number of files")
		}
		if hdr.Name != files[ctr].Name {
			t.Fatalf("Expected name %s but got %s", files[ctr].Name, hdr.Name)
		}
		b, err := ioutil.ReadAll(tr)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != files[ctr].Body {
			t.Fatalf("Expected body %s but got %s", files[ctr].Name, string(b))
		}
		ctr += 1
	}
}
