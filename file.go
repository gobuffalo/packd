package packd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
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
	if offset == 0 && whence == io.SeekStart {
		f.Reader = bytes.NewReader(f.original)
		return 0, nil
	}
	return -1, errors.New("unsuported Seek operation")
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

func (s *virtualFile) Read(p []byte) (int, error) {
	if s.Reader == nil {
		s.Reader = bytes.NewReader(s.original)
	}
	i, err := s.Reader.Read(p)

	if i == 0 || err == io.EOF {
		s.Reader = bytes.NewReader(s.original)
	}
	return i, err
}

func (s *virtualFile) Write(p []byte) (int, error) {
	bb := &bytes.Buffer{}
	i, err := bb.Write(p)
	if err != nil {
		return i, errors.WithStack(err)
	}
	s.original = bb.Bytes()
	s.info = fileInfo{
		Path:    s.name,
		size:    int64(bb.Len()),
		modTime: time.Now(),
	}
	return i, nil
}

// NewDir returns a new "virtual" file
func NewFile(name string, r io.Reader) (File, error) {
	return buildFile(name, r)
}

// NewDir returns a new "virtual" directory
func NewDir(name string) (File, error) {
	v, err := buildFile(name, nil)
	if err != nil {
		return v, errors.WithStack(err)
	}
	v.info.isDir = true
	return v, nil
}

func buildFile(name string, r io.Reader) (*virtualFile, error) {
	bb := &bytes.Buffer{}
	if r != nil {
		io.Copy(bb, r)
	}
	return &virtualFile{
		name:     name,
		original: bb.Bytes(),
		info: fileInfo{
			Path:    name,
			size:    int64(bb.Len()),
			modTime: time.Now(),
		},
	}, nil
}
