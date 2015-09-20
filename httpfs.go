package httpfs

import (
	"errors"
	"flag"
	"io"

	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var PathPrefix string
var Root string

func init() {
	flag.StringVar(&PathPrefix, "fspath", "/fs/", "URL path at which httpfs files are served. must start with slash")
	flag.StringVar(&Root, "fsroot", "/tmp/", "OS path at which files are stored and served. must start with slash")
}

func getPath(r *http.Request) (string, error) {
	path := Root + strings.TrimPrefix(r.URL.Path, PathPrefix)
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(abs, Root) {
		return "", errors.New("directory traversal attack!")
	}
	return abs, nil
}

// `func(w http.ResponseWriter, r *http.Request)` is a function type
// that can be passed to http.HandleFunc's second argument
//
// usage:
// http.HandleFunc("/fs/", httpfs.Handle)
//
func Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT", "POST":
		HandlePut(w, r)
	case "DELETE":
		HandleDelete(w, r)
	case "GET":
		HandleGet(w, r)
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

// to be used for PUT handler, but can be used as POST handler too
// it's easier to use as a POST handler to be used by go
//
// curl example:
// curl -XPOST --data-binary "hoge" -v http://127.0.0.1:8080/fs/foo
// curl -XPUT --data-binary "hoge" -v http://127.0.0.1:8080/fs/foo
//
// go example:
// resp, err := http.Post("http://localhost:8080/fs/foo", "", strings.NewReader("abc"))
// or
// req, err := http.NewRequest("PUT", "http://127.0.0.1:8080/fs/foo", strings.NewReader("abc"))
// resp, err := http.DefaultClient.Do(req)
func HandlePut(w http.ResponseWriter, r *http.Request) {
	wrapErrHandler(w, r, ErrHandlePut)
}

func ErrHandlePut(w http.ResponseWriter, r *http.Request) error {
	path, err := getPath(r)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		return err
	}

	return nil
}

// curl example:
// curl -XDELETE -v http://127.0.0.1:8080/fs/foo
//
// go example:
// req, err := http.NewRequest("DELETE", "http://127.0.0.1:8080/fs/foo", nil)
// resp, err := http.DefaultClient.Do(req)
func HandleDelete(w http.ResponseWriter, r *http.Request) {
	wrapErrHandler(w, r, ErrHandleDelete)
}

func ErrHandleDelete(w http.ResponseWriter, r *http.Request) error {
	path, err := getPath(r)
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "", http.StatusNotFound)
			return nil
		}
		return err
	}

	if info.IsDir() {
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
		return nil
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

// curl example:
// curl -XGET -v http://127.0.0.1:8080/fs/foo
//
// go example:
// resp, err := http.Get("http://localhost:8080/fs/foo")
func HandleGet(w http.ResponseWriter, r *http.Request) {
	wrapErrHandler(w, r, ErrHandleGet)
}

func ErrHandleGet(w http.ResponseWriter, r *http.Request) error {
	path, err := getPath(r)
	if err != nil {
		return err
	}
	http.ServeFile(w, r, path)
	return nil
}

// alias for http.HandleFunc's handler or http.Handler's ServeHTTP method
type handler func(http.ResponseWriter, *http.Request)

// similar to above but returns an error
// it can be converted with handleErr
type errHandler func(http.ResponseWriter, *http.Request) error

// converts errHandler to handler
// if an error is returned, return a response with InternalServerError
func wrapErrHandler(w http.ResponseWriter, r *http.Request, handler errHandler) {
	err := handler(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
