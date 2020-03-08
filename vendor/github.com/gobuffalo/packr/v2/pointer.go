package packr

import (
	"github.com/gobuffalo/packr/v2/file"
	"github.com/gobuffalo/packr/v2/file/resolver"
	"github.com/gobuffalo/packr/v2/plog"
)

// Pointer is a resolvr which resolves
// a file from a different box.
type Pointer struct {
	ForwardBox  string
	ForwardPath string
}

var _ resolver.Resolver = Pointer{}

// Resolve attempts to find the file in the specific box
// with the specified key
func (p Pointer) Resolve(box string, path string) (file.File, error) {
	plog.Debug(p, "Resolve", "box", box, "path", path, "forward-box", p.ForwardBox, "forward-path", p.ForwardPath)
	b, err := findBox(p.ForwardBox)
	if err != nil {
		return nil, err
	}
	f, err := b.Resolve(p.ForwardPath)
	if err != nil {
		return f, err
	}
	plog.Debug(p, "Resolve", "box", box, "path", path, "file", f)
	return file.NewFileR(path, f)
}
