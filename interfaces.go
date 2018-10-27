package packd

import (
	"io"
	"net/http"
	"os"
)

type WalkFunc func(string, File) error

type Box interface {
	HTTPBox
	Lister
	Addable
	Finder
	Walkable
	Haser
}

type Haser interface {
	Has(string) bool
}

type Walkable interface {
	Walk(wf WalkFunc) error
}

type Finder interface {
	Find(string) ([]byte, error)
	FindString(name string) (string, error)
}

type HTTPBox interface {
	Open(name string) (http.File, error)
}

type Lister interface {
	List() []string
}

type Addable interface {
	AddString(path string, t string) error
	AddBytes(path string, t []byte) error
}

type File interface {
	io.ReadCloser
	io.Writer
	FileInfo() (os.FileInfo, error)
	Readdir(count int) ([]os.FileInfo, error)
	Seek(offset int64, whence int) (int64, error)
	Stat() (os.FileInfo, error)
}

type LegacyBox interface {
	String(name string) string
	MustString(name string) (string, error)
	Bytes(name string) []byte
	MustBytes(name string) ([]byte, error)
}
