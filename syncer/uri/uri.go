package uri

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"time"
)

var protocolRegistry map[string]reflect.Type

type Visitor interface {
	Visit(Uri) error
}

//Parser parse a uri string to interface Uri
type Parser interface {
	Parse(string) (Uri, error)
}

func Register(scheme string, typ interface{}) error {
	_, dup := protocolRegistry[scheme]
	if dup {
		return DupError{scheme}
	}
	protocolRegistry[scheme] = reflect.TypeOf(typ).Elem()
	return nil
}

func init() {
	if protocolRegistry == nil {
		protocolRegistry = make(map[string]reflect.Type, 4)
	}
	Register("local", (*UriLocal)(nil))

}

type Uri interface {
	Scheme() string
	Host() string
	Path() string
	Uri() string
	Abs() string
	Parent() (Uri, error)

	Create(isDir bool, fi os.FileMode) error
	OpenRead() (io.ReadCloser, error)
	OpenWrite() (io.WriteCloser, error)
	Remove() error
	Walk(v Visitor) chan error

	IsDir() bool
	Exist() bool
	IsAbs() bool

	Mode() os.FileMode
	ModTime() time.Time

	// Read(p []byte)(n int,err error)
	// Write(p []byte)(n int,err error)
	// Close()error

	setHost(h string)
	setPath(p string)
	setScheme(s string)
}

func Parse(u string) (Uri, error) {

	urlp, err := url.Parse(u)
	if err != nil {
		return nil, ParseError{u, err.Error()}
	}

	UriType, exist := protocolRegistry[urlp.Scheme]
	if !exist {
		return nil, ProtocolError{urlp.Scheme, "protocol not supported."}
	}

	UriVal := reflect.New(UriType)
	i := UriVal.Interface()
	Urip, ok := i.(Uri)
	if !ok {
		return nil, ProtocolError{urlp.Scheme, "protocol not fully supported."}
	}
	Urip.setScheme(urlp.Scheme)
	Urip.setHost(urlp.Host)
	Urip.setPath(urlp.Path)

	return Urip, nil

}

type ParseError struct {
	Uri     string
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("Error parsing %s on %s: %s", e.Uri, runtime.GOOS, e.Message)
}

type ProtocolError struct {
	Protocol string
	Message  string
}

func (e ProtocolError) Error() string {
	return fmt.Sprintf("When handling protocol %s: %s", e.Protocol, e.Message)
}

type OpenError struct {
	Uri     string
	Message string
}

func (e OpenError) Error() string {
	return fmt.Sprintf("Open %s: %s", e.Uri, e.Message)
}

type DupError struct {
	Typ string
}

func (e DupError) Error() string {
	return fmt.Sprintln("Duplicate uri type: %s", e.Typ)
}
