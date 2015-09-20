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
	goji.Get("/fs/*", httpfs.HandleGet)
	goji.Put("/fs/*", httpfs.HandlePut)
	goji.Delete("/fs/*", httpfs.HandleDelete)
	goji.ServeListener(bind.Socket(":10002"))
}

func main() {
	flag.Parse()
	go servePlainHTTP()
	go serveMartini()
	go serveGoji()

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