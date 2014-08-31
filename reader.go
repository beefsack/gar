package gar

import (
	"fmt"
	"io"
)

type Reader struct {
	GarInfo GarInfo
	r       io.ReadSeeker
	buf     []byte
	n       int
	eof     bool
}

func NewReader(r io.ReadSeeker) (*Reader, error) {
	gr := &Reader{
		r:   r,
		buf: make([]byte, 2048+GarHeaderSize*2),
	}
	is, info, err := Stat(r)
	if err != nil {
		return nil, err
	}
	if !is {
		return nil, fmt.Errorf("stream is not a valid gar archive")
	}
	gr.GarInfo = info
	if _, err = r.Seek(-(int64(GarHeaderSize) + info.Size), 2); err != nil {
		return gr, fmt.Errorf("could not seek to start of tar, %v", err)
	}
	return gr, nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	l := len(p)
	if len(p) == 0 {
		return 0, nil
	}
	toRead := r.n
	if toRead > GarHeaderSize {
		toRead -= GarHeaderSize // Always keep the header size in the buffer.
	}
	if toRead > l {
		toRead = l // We won't copy more than cap(p).
	}
	if toRead > 0 {
		copy(p, r.buf[:toRead])     // Copy part of the buffer into p.
		copy(r.buf, r.buf[toRead:]) // Flush that part of the buffer.
		r.n -= toRead
	}
	if !r.eof {
		// If we're not at true EOF yet, read a bit more into the buffer.
		newN, err := r.r.Read(r.buf[r.n:])
		r.n += newN
		if err == io.EOF {
			r.eof = true
		} else if err != nil {
			return toRead, err
		}
	} else if r.n <= GarHeaderSize {
		// We now have only the header left in the buffer, so this becomes EOF.
		return toRead, io.EOF
	}
	if toRead == 0 {
		// toRead will be 0 on first pass, call again.
		return r.Read(p)
	}
	return toRead, nil
}
