package gar

import "bytes"

type ByteReaderCloser struct {
	*bytes.Reader
}

func NewByteReaderCloser(b []byte) *ByteReaderCloser {
	return &ByteReaderCloser{bytes.NewReader(b)}
}

func (brc *ByteReaderCloser) Close() error {
	return nil
}
