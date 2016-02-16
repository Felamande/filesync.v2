package uri

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UriLocal struct {
	scheme string
	path   string
	host   string
	uri    string
	abs    string
}

func (u *UriLocal) Host() string {
	return u.host

}

func (u *UriLocal) Uri() string {
	if u.uri != "" {
		return u.uri
	}
	u.uri = u.scheme + "://" + u.host + u.path
	return u.uri

}

func (u *UriLocal) Abs() string {
	if u.abs != "" {
		return u.abs
	}
	u.abs = u.host + u.path
	return u.abs
}

func (u *UriLocal) IsAbs() bool {
	return filepath.IsAbs(u.Abs())
}

func (u *UriLocal) Scheme() string {
	return u.scheme
}

func (u *UriLocal) Mode() os.FileMode {
	fi, err := os.Stat(u.host + u.path)
	if err != nil {
		return os.ModePerm
	}
	return fi.Mode()
}

func (u *UriLocal) Exist() bool {
	fi, _ := os.Stat(u.Abs())
	if fi != nil {
		return true
	}
	return false
}

func (u *UriLocal) ModTime() time.Time {
	fi, _ := os.Stat(u.host + u.path)
	if fi == nil {
		return time.Date(1970, time.January, 1, 0, 0, 0, 1, time.Local)
	}
	return fi.ModTime()
}

func (u *UriLocal) Create(IsDir bool, m os.FileMode) (err error) {

	if u.Exist() {
		return nil
	}

	if IsDir {
		err = os.Mkdir(u.Abs(), m)
		if err != nil {
			return
		}
	} else {
		var fd *os.File
		fd, err = os.OpenFile(u.Abs(), os.O_CREATE, m)
		defer fd.Close()
		if err != nil {
			return
		}
	}
	return nil
}

func (u *UriLocal) OpenRead() (io.ReadCloser, error) {
	if !filepath.IsAbs(u.Abs()) {
		return nil, OpenError{u.Uri(), "is not an absolute path."}
	}
	return os.OpenFile(u.Abs(), os.O_RDONLY, u.Mode())

}
func (u *UriLocal) OpenWrite() (io.WriteCloser, error) {
	if !filepath.IsAbs(u.Abs()) {
		return nil, OpenError{u.Uri(), "is not an absolute path."}
	}

	return os.OpenFile(u.Abs(), os.O_WRONLY, u.Mode())

}

func (u *UriLocal) Remove() error {
	return os.Remove(u.Abs())
}

func (u *UriLocal) Walk(v Visitor) (c chan error) {
	c = make(chan error, 4)
	go func(vv Visitor) {
		if !u.IsDir() {
			c <- errors.New("walk " + u.Abs() + ": is not a directory")
			return
		}
		walkFunc := func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				c <- err
				return nil
			}
			path = strings.Replace(path, "\\", "/", -1)
			urip, err := Parse(u.Scheme() + "://" + path)
			if err != nil {
				c <- err
				return nil
			}
			err = v.Visit(urip) //Visit add subtree to wacther
			if err != nil {
				c <- err
			}
			return nil
		}
		err := filepath.Walk(u.Abs(), walkFunc)

		if err != nil {
			c <- err
		}

		close(c)
	}(v)
	return c
}

func (u *UriLocal) Path() string {

	return u.path
}

func (u *UriLocal) Parent() (Uri, error) {
	p := filepath.Dir(u.Abs())
	p = strings.Replace(p, "\\", "/", -1)
	return Parse(u.Scheme() + "://" + p)
}

func (u *UriLocal) IsDir() bool {
	fi, _ := os.Stat(u.Abs())
	if fi == nil {
		return false
		//TODO
	}

	return fi.IsDir()
}

func (u *UriLocal) setHost(h string) {
	u.host = h
}
func (u *UriLocal) setPath(p string) {
	u.path = p
}
func (u *UriLocal) setScheme(s string) {
	u.scheme = s
}
