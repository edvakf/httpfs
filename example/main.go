package main

// go run main.go

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/edvakf/httpfs"
	"github.com/go-martini/martini"
	"github.com/gorilla/mux"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/bind"
)

func servePlainHTTP() {
	http.HandleFunc("/fs/", httpfs.Handle)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func serveMartini() {
	m := martini.Classic()
	m.Get("/fs/**", httpfs.HandleGet)
	m.Put("/fs/**", httpfs.HandlePut)
	m.Delete("/fs/**", httpfs.HandleDelete)
	m.RunOnAddr(":10001")
}

func serveGoji() {
	goji.Get("/fs/[a-zA-Z0-9._/-]+", httpfs.HandleGet)
	goji.Put("/fs/[a-zA-Z0-9._/-]+", httpfs.HandlePut)
	goji.Delete("/fs/[a-zA-Z0-9._/-]+", httpfs.HandleDelete)
	goji.ServeListener(bind.Socket(":10002"))
}

func serveGorilla() {
	r := mux.NewRouter()
	r.HandleFunc("/{path:.+}", httpfs.HandleGet).Methods("GET")
	r.HandleFunc("/{path:.+}", httpfs.HandlePut).Methods("PUT")
	r.HandleFunc("/{path:.+}", httpfs.HandleDelete).Methods("DELETE")
	http.ListenAndServe(":10003", r)
}

func main() {
	flag.Parse()
	go servePlainHTTP()
	go serveMartini()
	go serveGoji()
	go serveGorilla()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for {
		switch <-ch {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			os.Exit(0)
		}
	}
}
