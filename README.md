# httpfs

Serve file system access via http.

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
	goji.Get("/fs/*", httpfs.HandleGet)
	goji.Put("/fs/*", httpfs.HandlePut)
	goji.Delete("/fs/*", httpfs.HandleDelete)
	goji.ServeListener(bind.Socket(":10002"))
}
```
