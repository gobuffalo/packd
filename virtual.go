package packd

import (
	"bytes"
	"io"
	"os"
	"time"
)

var virtualFileModTime = time.Now()
var _ File = virtualFile{}

type virtualFile struct {
	*bytes.Buffer
	Name string
	info fileInfo
}

func (f virtualFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f virtualFile) FileInfo() (os.FileInfo, error) {
	return f.info, nil
}

func (f virtualFile) Close() error {
	return nil
}

func (f virtualFile) Readdir(count int) ([]os.FileInfo, error) {
	return []os.FileInfo{f.info}, nil
}

func (f virtualFile) Stat() (os.FileInfo, error) {
	return f.info, nil
}

// NewDir returns a new "virtual" file
func NewFile(name string, r io.Reader) (File, error) {
	bb := &bytes.Buffer{}
	io.Copy(bb, r)
	return virtualFile{
		Buffer: bb,
		Name:   name,
		info: fileInfo{
			Path:     name,
			Contents: bb.Bytes(),
			size:     int64(bb.Len()),
			modTime:  time.Now(),
		},
	}, nil
}

// NewDir returns a new "virtual" directory
func NewDir(name string) (File, error) {
	bb := &bytes.Buffer{}
	return virtualFile{
		Buffer: bb,
		Name:   name,
		info: fileInfo{
			Path:     name,
			Contents: bb.Bytes(),
			size:     int64(bb.Len()),
			modTime:  time.Now(),
			isDir:    true,
		},
	}, nil
}

type fileInfo struct {
	Path     string
	Contents []byte
	size     int64
	modTime  time.Time
	isDir    bool
}

func (f fileInfo) Name() string {
	return f.Path
}

func (f fileInfo) Size() int64 {
	return f.size
}

func (f fileInfo) Mode() os.FileMode {
	return 0444
}

func (f fileInfo) ModTime() time.Time {
	return f.modTime
}

func (f fileInfo) IsDir() bool {
	return f.isDir
}

func (f fileInfo) Sys() interface{} {
	return nil
}
