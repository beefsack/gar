package gar

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"
)

func TestWriter(t *testing.T) {
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
	w := NewWriter(buf)
	// Write files as tar to buffer.
	for _, file := range files {
		if err := w.WriteFileWithName(
			file.Name,
			bytes.NewReader([]byte(file.Body)),
		); err != nil {
			t.Fatal(err)
		}
	}
	w.Close(CompressionNone)

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
		ctr += 1
	}
}
