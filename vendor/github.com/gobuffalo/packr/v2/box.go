package packr

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2/file"
	"github.com/gobuffalo/packr/v2/file/resolver"
	"github.com/gobuffalo/packr/v2/internal/takeon/github.com/markbates/oncer"
	"github.com/gobuffalo/packr/v2/plog"
)

var _ packd.Box = &Box{}
var _ packd.HTTPBox = &Box{}
var _ packd.Addable = &Box{}
var _ packd.Walkable = &Box{}
var _ packd.Finder = &Box{}

// Box represent a folder on a disk you want to
// have access to in the built Go binary.
type Box struct {
	Path            string            `json:"path"`
	Name            string            `json:"name"`
	ResolutionDir   string            `json:"resolution_dir"`
	DefaultResolver resolver.Resolver `json:"default_resolver"`
	resolvers       resolversMap
	dirs            dirsMap
}

// NewBox returns a Box that can be used to
// retrieve files from either disk or the embedded
// binary.
// Deprecated: Use New instead.
func NewBox(path string) *Box {
	oncer.Deprecate(0, "packr.NewBox", "Use packr.New instead.")
	return New(path, path)
}

// New returns a new Box with the name of the box
// and the path of the box.
func New(name string, path string) *Box {
	plog.Debug("packr", "New", "name", name, "path", path)
	b, _ := findBox(name)
	if b != nil {
		return b
	}

	b = construct(name, path)
	plog.Debug(b, "New", "Box", b, "ResolutionDir", b.ResolutionDir)
	b, err := placeBox(b)
	if err != nil {
		panic(err)
	}

	return b
}

// Folder returns a Box that will NOT be packed.
// This is useful for writing tests or tools that
// need to work with a folder at runtime.
func Folder(path string) *Box {
	return New(path, path)
}

// SetResolver allows for the use of a custom resolver for
// the specified file
func (b *Box) SetResolver(file string, res resolver.Resolver) {
	d := filepath.Dir(file)
	b.dirs.Store(d, true)
	plog.Debug(b, "SetResolver", "file", file, "resolver", fmt.Sprintf("%T", res))
	b.resolvers.Store(resolver.Key(file), res)
}

// AddString converts t to a byteslice and delegates to AddBytes to add to b.data
func (b *Box) AddString(path string, t string) error {
	return b.AddBytes(path, []byte(t))
}

// AddBytes sets t in b.data by the given path
func (b *Box) AddBytes(path string, t []byte) error {
	m := map[string]file.File{}
	f, err := file.NewFile(path, t)
	if err != nil {
		return err
	}
	m[resolver.Key(path)] = f
	res := resolver.NewInMemory(m)
	b.SetResolver(path, res)
	return nil
}

// FindString returns either the string of the requested
// file or an error if it can not be found.
func (b *Box) FindString(name string) (string, error) {
	bb, err := b.Find(name)
	return string(bb), err
}

// Find returns either the byte slice of the requested
// file or an error if it can not be found.
func (b *Box) Find(name string) ([]byte, error) {
	f, err := b.Resolve(name)
	if err != nil {
		return []byte(""), err
	}
	bb := &bytes.Buffer{}
	io.Copy(bb, f)
	return bb.Bytes(), nil
}

// Has returns true if the resource exists in the box
func (b *Box) Has(name string) bool {
	_, err := b.Find(name)
	return err == nil
}

// HasDir returns true if the directory exists in the box
func (b *Box) HasDir(name string) bool {
	oncer.Do("packr2/box/HasDir"+b.Name, func() {
		for _, f := range b.List() {
			for d := filepath.Dir(f); d != "."; d = filepath.Dir(d) {
				b.dirs.Store(d, true)
			}
		}
	})
	if name == "/" {
		return b.Has("index.html")
	}
	_, ok := b.dirs.Load(name)
	return ok
}

// Open returns a File using the http.File interface
func (b *Box) Open(name string) (http.File, error) {
	plog.Debug(b, "Open", "name", name)
	f, err := b.Resolve(name)
	if err != nil {
		if len(filepath.Ext(name)) == 0 {
			return b.openWoExt(name)
		}
		return f, err
	}
	f, err = file.NewFileR(name, f)
	plog.Debug(b, "Open", "name", f.Name(), "file", f.Name())
	return f, err
}

func (b *Box) openWoExt(name string) (http.File, error) {
	if !b.HasDir(name) {
		id := path.Join(name, "index.html")
		if b.Has(id) {
			return b.Open(id)
		}
		return nil, os.ErrNotExist
	}
	d, err := file.NewDir(name)
	plog.Debug(b, "Open", "name", name, "dir", d)
	return d, err
}

// List shows "What's in the box?"
func (b *Box) List() []string {
	var keys []string

	b.Walk(func(path string, info File) error {
		if info == nil {
			return nil
		}
		finfo, _ := info.FileInfo()
		if !finfo.IsDir() {
			keys = append(keys, path)
		}
		return nil
	})
	sort.Strings(keys)
	return keys
}

// Resolve will attempt to find the file in the box,
// returning an error if the find can not be found.
func (b *Box) Resolve(key string) (file.File, error) {
	key = strings.TrimPrefix(key, "/")

	var r resolver.Resolver

	b.resolvers.Range(func(k string, vr resolver.Resolver) bool {
		lk := strings.ToLower(resolver.Key(k))
		lkey := strings.ToLower(resolver.Key(key))
		if lk == lkey {
			r = vr
			return false
		}
		return true
	})

	if r == nil {
		r = b.DefaultResolver
		if r == nil {
			r = resolver.DefaultResolver
			if r == nil {
				return nil, fmt.Errorf("resolver.DefaultResolver is nil")
			}
		}
	}
	plog.Debug(r, "Resolve", "box", b.Name, "key", key)

	f, err := r.Resolve(b.Name, key)
	if err != nil {
		z, err := resolver.ResolvePathInBase(resolver.OsPath(b.ResolutionDir), filepath.FromSlash(path.Clean("/"+resolver.OsPath(key))))
		if err != nil {
			plog.Debug(r, "Resolve", "box", b.Name, "key", key, "err", err)
			return f, err
		}

		f, err = r.Resolve(b.Name, z)
		if err != nil {
			plog.Debug(r, "Resolve", "box", b.Name, "key", z, "err", err)
			return f, err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return f, err
		}
		f, err = file.NewFile(key, b)
		if err != nil {
			return f, err
		}
	}
	plog.Debug(r, "Resolve", "box", b.Name, "key", key, "file", f.Name())
	return f, nil
}
