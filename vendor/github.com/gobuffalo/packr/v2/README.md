# Packr (v2)

[![GoDoc](https://godoc.org/github.com/gobuffalo/packr/v2?status.svg)](https://godoc.org/github.com/gobuffalo/packr/v2)

Packr is a simple solution for bundling static assets inside of Go binaries. Most importantly it does it in a way that is friendly to developers while they are developing.

## Intro Video

To get an idea of the what and why of Packr, please enjoy this short video: [https://vimeo.com/219863271](https://vimeo.com/219863271).

## Library Installation

```text
$ go get -u github.com/gobuffalo/packr/v2/...
```

## Binary Installation

```text
$ go get -u github.com/gobuffalo/packr/v2/packr2
```

## New File Format FAQs

In version `v2.0.0` the file format changed and is not backward compatible with the `packr-v1.x` library.

#### Can `packr-v1.x` read the new format?

No, it can not. Because of the way the new file format works porting it to `packr-v1.x` would be difficult. PRs are welcome though. :)

#### Can `packr-v2.x` read `packr-v1.x` files?

Yes it can, but that ability will eventually be phased out. Because of that we recommend moving to the new format.

#### Can `packr-v2.x` generate `packr-v1.x` files?

Yes it can, but that ability will eventually be phased out. Because of that we recommend moving to the new format.

The `--legacy` command is available on all commands that generate `-packr.go` files.

```bash
$ packr2 --legacy
```

## Usage

### In Code

The first step in using Packr is to create a new box. A box represents a folder on disk. Once you have a box you can get `string` or `[]byte` representations of the file.

```go
// set up a new box by giving it a name and an optional (relative) path to a folder on disk:
box := packr.New("My Box", "./templates")

// Get the string representation of a file, or an error if it doesn't exist:
html, err := box.FindString("index.html")

// Get the []byte representation of a file, or an error if it doesn't exist:
html, err := box.Find("index.html")
```

### What is a Box?

A box represents a folder, and any sub-folders, on disk that you want to have access to in your binary. When compiling a binary using the `packr2` CLI the contents of the folder will be converted into Go files that can be compiled inside of a "standard" go binary. Inside of the compiled binary the files will be read from memory. When working locally the files will be read directly off of disk. This is a seamless switch that doesn't require any special attention on your part.

#### Example

Assume the follow directory structure:

```
├── main.go
└── templates
    ├── admin
    │   └── index.html
    └── index.html
```

The following program will read the `./templates/admin/index.html` file and print it out.

```go
package main

import (
  "fmt"

  "github.com/gobuffalo/packr/v2"
)

func main() {
  box := packr.New("myBox", "./templates")

  s, err := box.FindString("admin/index.html")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(s)
}
```

### Development Made Easy

In order to get static files into a Go binary, those files must first be converted to Go code. To do that, Packr, ships with a few tools to help build binaries. See below.

During development, however, it is painful to have to keep running a tool to compile those files.

Packr uses the following resolution rules when looking for a file:

1. Look for the file in-memory (inside a Go binary)
1. Look for the file on disk (during development)

Because Packr knows how to fall through to the file system, developers don't need to worry about constantly compiling their static files into a binary. They can work unimpeded.

Packr takes file resolution a step further. When declaring a new box you use a relative path, `./templates`. When Packr receives this call it calculates out the absolute path to that directory. By doing this it means you can be guaranteed that Packr can find your files correctly, even if you're not running in the directory that the box was created in. This helps with the problem of testing, where Go changes the `pwd` for each package, making relative paths difficult to work with. This is not a problem when using Packr.

---

## Usage with HTTP

A box implements the [`http.FileSystem`](https://golang.org/pkg/net/http/#FileSystem) interface, meaning it can be used to serve static files.

```go
package main

import (
	"net/http"

	"github.com/gobuffalo/packr/v2"
)

func main() {
	box := packr.New("someBoxName", "./templates")

	http.Handle("/", http.FileServer(box))
	http.ListenAndServe(":3000", nil)
}
```

---

## Building a Binary

Before you build your Go binary, run the `packr2` command first. It will look for all the boxes in your code and then generate `.go` files that pack the static files into bytes that can be bundled into the Go binary.

```
$ packr2
```

Then run your `go build command` like normal.

*NOTE*: It is not recommended to check-in these generated `-packr.go` files. They can be large, and can easily become out of date if not careful. It is recommended that you always run `packr2 clean` after running the `packr2` tool.

#### Cleaning Up

When you're done it is recommended that you run the `packr2 clean` command. This will remove all of the generated files that Packr created for you.

```
$ packr2 clean
```

Why do you want to do this? Packr first looks to the information stored in these generated files, if the information isn't there it looks to disk. This makes it easy to work with in development.

---

## Debugging

The `packr2` command passes all arguments down to the underlying `go` command, this includes the `-v` flag to print out `go build` information. Packr looks for the `-v` flag, and will turn on its own verbose logging. This is very useful for trying to understand what the `packr` command is doing when it is run.

---

## FAQ

### Compilation Errors with Go Templates

Q: I have a program with Go template files, those files are named `foo.go` and look like the following:

```
// Copyright {{.Year}} {{.Author}}. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package {{.Project}}
```

When I run `packr2` I get errors like:

```
expected 'IDENT', found '{'
```

A: Packr works by searching your `.go` files for [`github.com/gobuffalo/packr/v2#New`](https://godoc.org/github.com/gobuffalo/packr/v2#New) or [`github.com/gobuffalo/packr/v2#NewBox`](https://godoc.org/github.com/gobuffalo/packr/v2#NewBox) calls. Because those files aren't "proper" Go files, Packr can't parse them to find the box declarations. To fix this you need to tell Packr to ignore those files when searching for boxes. A couple solutions to this problem are:

* Name the files something else. The `.tmpl` extension is the idiomatic way of naming these types of files.
* Rename the folder containing these files to start with an `_`, for example `_templates`. Packr, like Go, will ignore folders starting with the `_` character when searching for boxes.

### Dynamic Box Paths

Q: I need to set the path of a box using a variable, but `packr.New("foo", myVar)` doesn't work correctly.

A: Packr attempts to "automagically" set it's resolution directory when using [`github.com/gobuffalo/packr/v2#New`](https://godoc.org/github.com/gobuffalo/packr/v2#New), however, for dynamic paths you need to set it manually:

```go
box := packr.New("foo", "|")
box.ResolutionDir = myVar
```

### I don't want to pack files, but still use the Packr interface.

Q: I want to write code that using the Packr tools, but doesn't actually pack the files into my binary. How can I do that?

A: Using [`github.com/gobuffalo/packr/v2#Folder`](https://godoc.org/github.com/gobuffalo/packr/v2#Folder) gives you back a `*packr.Box` that can be used as normal, but is excluded by the Packr tool when compiling.

### Packr Finds No Boxes

Q: I run `packr2 -v` but it doesn't find my boxes:

```
DEBU[2019-03-18T18:48:52+01:00] *parser.Parser#NewFromRoots found prospects=0
DEBU[2019-03-18T18:48:52+01:00] found 0 boxes
```

A: Packr works by parsing `.go` files to find [`github.com/gobuffalo/packr/v2#Box`](https://godoc.org/github.com/gobuffalo/packr/v2#Box) and [`github.com/gobuffalo/packr/v2#NewBox`](https://godoc.org/github.com/gobuffalo/packr/v2#NewBox) declarations. If there aren't any `.go` in the folder that `packr2` is run in it can not find those declarations. To fix this problem run the `packr2` command in the directory containing your `.go` files.

### Box Interfaces

Q: I want to be able to easily test my applications by passing in mock boxes. How do I do that?

A: Packr boxes and files conform to the interfaces found at [`github.com/gobuffalo/packd`](https://godoc.org/github.com/gobuffalo/packd). Change your application to use those interfaces instead of the concrete Packr types.

```go
// using concrete type
func myFunc(box *packr.Box) {}

// using interfaces
func myFunc(box packd.Box) {}
```
