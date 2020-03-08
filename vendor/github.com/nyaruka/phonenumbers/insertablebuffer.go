// This file was originally entirely copied from bytes/buffer.go
// in the go standard library, and was then modified from there.

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package phonenumbers

// Simple byte buffer for marshaling data.

import (
	"bytes"
	"errors"
	"io"
	"unicode/utf8"
)

// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type Builder struct {
	buf       []byte            // contents are the bytes buf[off : len(buf)]
	off       int               // read at &buf[off], write at &buf[len(buf)]
	runeBytes [utf8.UTFMax]byte // avoid allocation of slice on each WriteByte or Rune
	bootstrap [64]byte          // memory to hold first slice; helps small buffers (Printf) avoid allocation.
	lastRead  readOp            // last read operation, so that Unread* can work correctly.
}

// The readOp constants describe the last action performed on
// the buffer, so that UnreadRune and UnreadByte can
// check for invalid usage.
type readOp int

const (
	opInvalid  readOp = iota // Non-read operation.
	opReadRune               // Read rune.
	opRead                   // Any other read operation.
)

// Bytes returns a slice of the contents of the unread portion of the buffer;
// len(b.Bytes()) == b.Len().  If the caller changes the contents of the
// returned slice, the contents of the buffer will change provided there
// are no intervening method calls on the Buffer.
func (b *Builder) Bytes() []byte { return b.buf[b.off:] }

// String returns the contents of the unread portion of the buffer
// as a string.  If the Buffer is a nil pointer, it returns "<nil>".
func (b *Builder) String() string {
	if b == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(b.buf[b.off:])
}

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (b *Builder) Len() int { return len(b.buf) - b.off }

// Truncate discards all but the first n unread bytes from the buffer.
// It panics if n is negative or greater than the length of the buffer.
func (b *Builder) Truncate(n int) {
	b.lastRead = opInvalid
	switch {
	case n < 0 || n > b.Len():
		// Originall this was a panic, but I hate panicing (hehe),
		// so instead let's just ignore this issue.
		// panic("bytes.Buffer: truncation out of range")
		return
	case n == 0:
		// Reuse buffer space.
		b.off = 0
	}
	b.buf = b.buf[0 : b.off+n]
}

// Reset resets the buffer so it has no content.
// b.Reset() is the same as b.Truncate(0).
func (b *Builder) Reset() { b.Truncate(0) }

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Builder) grow(n int) int {
	m := b.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && b.off != 0 {
		b.Truncate(0)
	}
	if len(b.buf)+n > cap(b.buf) {
		var buf []byte
		if b.buf == nil && n <= len(b.bootstrap) {
			buf = b.bootstrap[0:]
		} else if m+n <= cap(b.buf)/2 {
			// We can slide things down instead of allocating a new
			// slice. We only need m+n <= cap(b.buf) to slide, but
			// we instead let capacity get twice as large so we
			// don't spend all our time copying.
			copy(b.buf[:], b.buf[b.off:])
			buf = b.buf[:m]
		} else {
			// not enough space anywhere
			buf = makeSlice(2*cap(b.buf) + n)
			copy(buf, b.buf[b.off:])
		}
		b.buf = buf
		b.off = 0
	}
	b.buf = b.buf[0 : b.off+m+n]
	return b.off + m
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Builder) Grow(n int) {
	if n < 0 {
		panic("bytes.Buffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[0:m]
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (b *Builder) Write(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	m := b.grow(len(p))
	return copy(b.buf[m:], p), nil
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
func (b *Builder) WriteString(s string) (n int, err error) {
	b.lastRead = opInvalid
	m := b.grow(len(s))
	return copy(b.buf[m:], s), nil
}

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with ErrTooLarge.
func (b *Builder) ReadFrom(r io.Reader) (n int64, err error) {
	b.lastRead = opInvalid
	// If buffer is empty, reset to recover space.
	if b.off >= len(b.buf) {
		b.Truncate(0)
	}
	for {
		if free := cap(b.buf) - len(b.buf); free < bytes.MinRead {
			// not enough space at end
			newBuf := b.buf
			if b.off+free < bytes.MinRead {
				// not enough space using beginning of buffer;
				// double buffer capacity
				newBuf = makeSlice(2*cap(b.buf) + bytes.MinRead)
			}
			copy(newBuf, b.buf[b.off:])
			b.buf = newBuf[:len(b.buf)-b.off]
			b.off = 0
		}
		m, e := r.Read(b.buf[len(b.buf):cap(b.buf)])
		b.buf = b.buf[0 : len(b.buf)+m]
		n += int64(m)
		if e == io.EOF {
			break
		}
		if e != nil {
			return n, e
		}
	}
	return n, nil // err is EOF, so return nil explicitly
}

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(bytes.ErrTooLarge)
		}
	}()
	return make([]byte, n)
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (b *Builder) WriteTo(w io.Writer) (n int64, err error) {
	b.lastRead = opInvalid
	if b.off < len(b.buf) {
		nBytes := b.Len()
		m, e := w.Write(b.buf[b.off:])
		if m > nBytes {
			panic("bytes.Buffer.WriteTo: invalid Write count")
		}
		b.off += m
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	// Buffer is now empty; reset.
	b.Truncate(0)
	return
}

// WriteByte appends the byte c to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match bufio.Writer's
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// ErrTooLarge.
func (b *Builder) WriteByte(c byte) error {
	b.lastRead = opInvalid
	m := b.grow(1)
	b.buf[m] = c
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to the
// buffer, returning its length and an error, which is always nil but is
// included to match bufio.Writer's WriteRune. The buffer is grown as needed;
// if it becomes too large, WriteRune will panic with ErrTooLarge.
func (b *Builder) WriteRune(r rune) (n int, err error) {
	if r < utf8.RuneSelf {
		b.WriteByte(byte(r))
		return 1, nil
	}
	n = utf8.EncodeRune(b.runeBytes[0:], r)
	b.Write(b.runeBytes[0:n])
	return n, nil
}

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained.  The return value n is the number of bytes read.  If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
func (b *Builder) Read(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	if b.off >= len(b.buf) {
		// Buffer is empty, reset to recover space.
		b.Truncate(0)
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}
	n = copy(p, b.buf[b.off:])
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (b *Builder) Next(n int) []byte {
	b.lastRead = opInvalid
	m := b.Len()
	if n > m {
		n = m
	}
	data := b.buf[b.off : b.off+n]
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return data
}

// ReadByte reads and returns the next byte from the buffer.
// If no byte is available, it returns error io.EOF.
func (b *Builder) ReadByte() (c byte, err error) {
	b.lastRead = opInvalid
	if b.off >= len(b.buf) {
		// Buffer is empty, reset to recover space.
		b.Truncate(0)
		return 0, io.EOF
	}
	c = b.buf[b.off]
	b.off++
	b.lastRead = opRead
	return c, nil
}

// ReadRune reads and returns the next UTF-8-encoded
// Unicode code point from the buffer.
// If no bytes are available, the error returned is io.EOF.
// If the bytes are an erroneous UTF-8 encoding, it
// consumes one byte and returns U+FFFD, 1.
func (b *Builder) ReadRune() (r rune, size int, err error) {
	b.lastRead = opInvalid
	if b.off >= len(b.buf) {
		// Buffer is empty, reset to recover space.
		b.Truncate(0)
		return 0, 0, io.EOF
	}
	b.lastRead = opReadRune
	c := b.buf[b.off]
	if c < utf8.RuneSelf {
		b.off++
		return rune(c), 1, nil
	}
	r, n := utf8.DecodeRune(b.buf[b.off:])
	b.off += n
	return r, n, nil
}

// UnreadRune unreads the last rune returned by ReadRune.
// If the most recent read or write operation on the buffer was
// not a ReadRune, UnreadRune returns an error.  (In this regard
// it is stricter than UnreadByte, which will unread the last byte
// from any read operation.)
func (b *Builder) UnreadRune() error {
	if b.lastRead != opReadRune {
		return errors.New("bytes.Buffer: UnreadRune: previous operation was not ReadRune")
	}
	b.lastRead = opInvalid
	if b.off > 0 {
		_, n := utf8.DecodeLastRune(b.buf[0:b.off])
		b.off -= n
	}
	return nil
}

// UnreadByte unreads the last byte returned by the most recent
// read operation.  If write has happened since the last read, UnreadByte
// returns an error.
func (b *Builder) UnreadByte() error {
	if b.lastRead != opReadRune && b.lastRead != opRead {
		return errors.New("bytes.Buffer: UnreadByte: previous operation was not a read")
	}
	b.lastRead = opInvalid
	if b.off > 0 {
		b.off--
	}
	return nil
}

// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
// If ReadBytes encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadBytes returns err != nil if and only if the returned data does not end in
// delim.
func (b *Builder) ReadBytes(delim byte) (line []byte, err error) {
	slice, err := b.readSlice(delim)
	// return a copy of slice. The buffer's backing array may
	// be overwritten by later calls.
	line = append(line, slice...)
	return
}

// readSlice is like ReadBytes but returns a reference to internal buffer data.
func (b *Builder) readSlice(delim byte) (line []byte, err error) {
	i := bytes.IndexByte(b.buf[b.off:], delim)
	end := b.off + i + 1
	if i < 0 {
		end = len(b.buf)
		err = io.EOF
	}
	line = b.buf[b.off:end]
	b.off = end
	b.lastRead = opRead
	return line, err
}

// ReadString reads until the first occurrence of delim in the input,
// returning a string containing the data up to and including the delimiter.
// If ReadString encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadString returns err != nil if and only if the returned data does not end
// in delim.
func (b *Builder) ReadString(delim byte) (line string, err error) {
	slice, err := b.readSlice(delim)
	return string(slice), err
}

// NewBuilder creates and initializes a new Buffer using buf as its initial
// contents.  It is intended to prepare a Buffer to read existing data.  It
// can also be used to size the internal buffer for writing. To do that,
// buf should have the desired capacity but a length of zero.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewBuilder(buf []byte) *Builder { return &Builder{buf: buf} }

// NewBuilderString creates and initializes a new Buffer using string s as its
// initial contents. It is intended to prepare a buffer to read an existing
// string.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewBuilderString(s string) *Builder {
	return &Builder{buf: []byte(s)}
}

// The whole reason we needed to copy this file was so we had access to the underlying slice

var ErrInvalidIndex = errors.New("insertablebuffer.Buf: invalid index")
var ErrFailedToGrow = errors.New("insertablebuffer.Buf: failed to grow buffer enough")

// Insert inserts the buffer at the desired position, growing the buffer
// as necessary. If i is less than zero, or greater than len(p), an error
// is returned.
func (b *Builder) Insert(i int, p []byte) (n int, err error) {
	b.lastRead = opInvalid
	if i < 0 || i > len(b.buf) {
		return -1, ErrInvalidIndex
	}

	m := b.grow(len(p))
	if len(b.buf)-m != len(p) {
		return -1, ErrFailedToGrow
	}
	copy(b.buf[i+len(p):], b.buf[i:len(b.buf)-len(p)])
	return copy(b.buf[i:i+len(p)], p), nil
}

// InsertString inserts the string at the desired position, growing the buffer
// as necessary. If i is less than zero, or greater than len(p), an error
// is returned.
func (b *Builder) InsertString(i int, p string) (n int, err error) {
	b.lastRead = opInvalid
	if i < 0 || i > len(b.buf) {
		return -1, ErrInvalidIndex
	}

	m := b.grow(len(p))
	if len(b.buf)-m != len(p) {
		return -1, ErrFailedToGrow
	}
	copy(b.buf[i+len(p):], b.buf[i:len(b.buf)-len(p)])
	return copy(b.buf[i:i+len(p)], p), nil
}

// ByteAt returns the byte at the given index. It returns
// ErrInvalidIndex if the given index is outside the buffer's
// length. There is no RuneAt currently, for I'm unsure
// what the desired behaviour should be (read rune by rune
// until we get to the ith rune, read the rune starting at i).
func (b *Builder) ByteAt(i int) (byte, error) {
	if i < 0 || i > len(b.buf) {
		return 0, ErrInvalidIndex
	}

	return b.buf[i], nil
}

func (b *Builder) ResetWith(buf []byte) (n int, err error) {
	b.Truncate(0)
	return b.Write(buf)
}

func (b *Builder) ResetWithString(s string) (n int, err error) {
	b.Truncate(0)
	return b.WriteString(s)
}
