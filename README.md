# httpfs

Serve file system access via http.

## flags

```
-fspath string
    URL path at which httpfs files are served. must start with slash (default "/fs/")
-fsroot string
    OS path at which files are stored and served. must start with slash (default "/tmp/")
```

## usage

```go
func servePlainHTTP() {
	http.HandleFunc("/fs/", httpfs.Handle)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
```

```go
func serveMartini() {
	m := martini.Classic()
	m.Get("/fs/**", httpfs.HandleGet)
	m.Put("/fs/**", httpfs.HandlePut)
	m.Delete("/fs/**", httpfs.HandleDelete)
	m.RunOnAddr(":10001")
}
```

```go
func serveGoji() {
	goji.Get("/fs/[a-zA-Z0-9._/-]+", httpfs.HandleGet)
	goji.Put("/fs/[a-zA-Z0-9._/-]+", httpfs.HandlePut)
	goji.Delete("/fs/[a-zA-Z0-9._/-]+", httpfs.HandleDelete)
	goji.ServeListener(bind.Socket(":10002"))
}
```

```go
func serveGorilla() {
	r := mux.NewRouter()
	r.HandleFunc("/fs/{path:.+}", httpfs.HandleGet).Methods("GET")
	r.HandleFunc("/fs/{path:.+}", httpfs.HandlePut).Methods("PUT")
	r.HandleFunc("/fs/{path:.+}", httpfs.HandleDelete).Methods("DELETE")
	http.ListenAndServe(":10003", r)
}
```

## client

By curl

```
$ curl -XPUT --data-binary "foobar" -v http://127.0.0.1:10000/fs/foo/bar
```

```
$ curl -XGET -v http://127.0.0.1:10000/fs/foo/bar
```

```
$ curl -XDELETE -v http://127.0.0.1:10000/fs/foo/bar
```

Or by go

```go
resp, err := http.Post("http://127.0.0.1:10000/fs/foo/bar", "", strings.NewReader("foobar"))
if err != nil {
	//
}
defer resp.Body.Close()
```

```go
req, err := http.NewRequest("PUT", "http://127.0.0.1:10000/fs/foo/bar", strings.NewReader("foobar"))
if err != nil {
	//
}
resp, err := http.DefaultClient.Do(req)
if err != nil {
	//
}
defer resp.Body.Close()
```

```go
resp, err := http.Get("http://127.0.0.1:10000/fs/foo/bar")
if err != nil {
	//
}
defer resp.Body.Close()

b, err := ioutil.ReadAll(resp.Body)
```

```go
req, err := http.NewRequest("DELETE", "http://127.0.0.1:10000/fs/foo/bar", nil)
if err != nil {
	//
}
resp, err := http.DefaultClient.Do(req)
if err != nil {
	//
}
defer resp.Body.Close()
```
