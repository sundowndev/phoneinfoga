package packr

import (
	"sort"
	"strings"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2/file"
	"github.com/gobuffalo/packr/v2/file/resolver"
	"github.com/gobuffalo/packr/v2/plog"
)

// WalkFunc is used to walk a box
type WalkFunc = packd.WalkFunc

// Walk will traverse the box and call the WalkFunc for each file in the box/folder.
func (b *Box) Walk(wf WalkFunc) error {
	m := map[string]file.File{}

	dr := b.DefaultResolver
	if dr == nil {
		cd := resolver.OsPath(b.ResolutionDir)
		dr = &resolver.Disk{Root: cd}
	}
	if fm, ok := dr.(file.FileMappable); ok {
		for n, f := range fm.FileMap() {
			m[n] = f
		}
	}
	var err error
	b.resolvers.Range(func(n string, r resolver.Resolver) bool {
		var f file.File
		f, err = r.Resolve("", n)
		if err != nil {
			return false
		}
		keep := true
		for k := range m {
			if strings.EqualFold(k, n) {
				keep = false
			}
		}
		if keep {
			m[n] = f
		}
		return true
	})
	if err != nil {
		return err
	}

	var keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		osPath := resolver.OsPath(k)
		plog.Debug(b, "Walk", "path", k, "osPath", osPath)
		if err := wf(osPath, m[k]); err != nil {
			return err
		}
	}
	return nil
}

// WalkPrefix will call box.Walk and call the WalkFunc when it finds paths that have a matching prefix
func (b *Box) WalkPrefix(prefix string, wf WalkFunc) error {
	ipref := resolver.OsPath(prefix)
	return b.Walk(func(path string, f File) error {
		ipath := resolver.OsPath(path)
		if strings.HasPrefix(ipath, ipref) {
			if err := wf(path, f); err != nil {
				return err
			}
		}
		return nil
	})
}
