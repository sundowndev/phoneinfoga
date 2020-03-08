package packd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

var _ File = &virtualFile{}
var _ io.Reader = &virtualFile{}
var _ io.Writer = &virtualFile{}
var _ fmt.Stringer = &virtualFile{}

type virtualFile struct {
	io.Reader
	name     string
	info     fileInfo
	original []byte
}

func (f virtualFile) Name() string {
	return f.name
}

func (f *virtualFile) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.(*bytes.Reader).Seek(offset, whence)
}

func (f virtualFile) FileInfo() (os.FileInfo, error) {
	return f.info, nil
}

func (f *virtualFile) Close() error {
	return nil
}

func (f virtualFile) Readdir(count int) ([]os.FileInfo, error) {
	return []os.FileInfo{f.info}, nil
}

func (f virtualFile) Stat() (os.FileInfo, error) {
	return f.info, nil
}

func (f virtualFile) String() string {
	return string(f.original)
}

// Read reads the next len(p) bytes from the virtualFile and
// rewind read offset to 0 when it met EOF.
func (f *virtualFile) Read(p []byte) (int, error) {
	i, err := f.Reader.Read(p)

	if i == 0 || err == io.EOF {
		f.Seek(0, io.SeekStart)
	}
	return i, err
}

// Write copies byte slice p to content of virtualFile.
func (f *virtualFile) Write(p []byte) (int, error) {
	return f.write(p)
}

// write copies byte slice or data from io.Reader to content of the
// virtualFile and update related information of the virtualFile.
func (f *virtualFile) write(d interface{}) (c int, err error) {
	bb := &bytes.Buffer{}
	switch d.(type) {
	case []byte:
		c, err = bb.Write(d.([]byte))
	case io.Reader:
		if d != nil {
			i64, e := io.Copy(bb, d.(io.Reader))
			c = int(i64)
			err = e
		}
	default:
		err = fmt.Errorf("unknown type of argument")
	}

	if err != nil {
		return c, err
	}

	f.info.size = int64(c)
	f.info.modTime = time.Now()
	f.original = bb.Bytes()
	f.Reader = bytes.NewReader(f.original)
	return c, nil
}

// NewFile returns a new "virtual" file
func NewFile(name string, r io.Reader) (File, error) {
	return buildFile(name, r)
}

// NewDir returns a new "virtual" directory
func NewDir(name string) (File, error) {
	v, err := buildFile(name, nil)
	if err != nil {
		return v, err
	}
	v.info.isDir = true
	return v, nil
}

func buildFile(name string, r io.Reader) (*virtualFile, error) {
	vf := &virtualFile{
		name: name,
		info: fileInfo{
			Path:    name,
			modTime: time.Now(),
		},
	}

	var err error
	if r != nil {
		_, err = vf.write(r)
	} else {
		_, err = vf.write([]byte{}) // for safety
	}
	return vf, err
}
