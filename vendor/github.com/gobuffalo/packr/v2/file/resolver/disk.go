package resolver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gobuffalo/packr/v2/file"
	"github.com/gobuffalo/packr/v2/internal/takeon/github.com/karrick/godirwalk"
	"github.com/gobuffalo/packr/v2/plog"
)

var _ Resolver = &Disk{}

type Disk struct {
	Root string
}

func (d Disk) String() string {
	return String(&d)
}

func (d *Disk) Resolve(box string, name string) (file.File, error) {
	var err error
	path := OsPath(name)
	if !filepath.IsAbs(path) {
		path, err = ResolvePathInBase(OsPath(d.Root), path)
		if err != nil {
			return nil, err
		}
	}

	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, os.ErrNotExist
	}
	if bb, err := ioutil.ReadFile(path); err == nil {
		return file.NewFile(OsPath(name), bb)
	}
	return nil, os.ErrNotExist
}

// ResolvePathInBase returns a path that is guaranteed to be inside of the base directory or an error
func ResolvePathInBase(base, path string) (string, error) {
	// Determine the absolute file path of the base directory
	d, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}

	// Return the base directory if no file was requested
	if path == "/" || path == "\\" {
		return d, nil
	}

	// Resolve the absolute file path after combining the key with base
	p, err := filepath.Abs(filepath.Join(d, path))
	if err != nil {
		return "", err
	}

	// Verify that the resolved path is inside of the base directory
	if !strings.HasPrefix(p, d+string(filepath.Separator)) {
		return "", os.ErrNotExist
	}
	return p, nil
}

var _ file.FileMappable = &Disk{}

func (d *Disk) FileMap() map[string]file.File {
	moot := &sync.Mutex{}
	m := map[string]file.File{}
	root := OsPath(d.Root)
	if _, err := os.Stat(root); err != nil {
		return m
	}
	callback := func(path string, de *godirwalk.Dirent) error {
		if _, err := os.Stat(root); err != nil {
			return nil
		}
		if !de.IsRegular() {
			return nil
		}
		moot.Lock()
		name := strings.TrimPrefix(path, root+string(filepath.Separator))
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		m[name], err = file.NewFile(name, b)
		if err != nil {
			return err
		}
		moot.Unlock()
		return nil
	}
	err := godirwalk.Walk(root, &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback:            callback,
	})
	if err != nil {
		plog.Logger.Errorf("[%s] error walking %v", root, err)
	}
	return m
}
