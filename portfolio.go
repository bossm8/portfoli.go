package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"bossm.ch/portfolio/handler"
)

const (
	templateDir = "templates"
	staticDir   = "static"
)

var baseTpl = filepath.Join("templates", "html", "base.html")

func main() {

	addr := flag.String(
		"srv.address",
		"127.0.0.1",
		"Listen address for the protfolio server",
	)
	port := flag.Int(
		"srv.port",
		8080,
		"Listen port for the portfolio server",
	)
	flag.Parse()

	startServer(fmt.Sprintf("%s:%d", *addr, *port))

}

func startServer(addr string) {

	fs := http.FileServer(http.Dir(staticDir))

	_http := &handler.RegexHandler{}

	_http.Handle("/static/", http.StripPrefix("/static", fs))
	_http.Handle("/favicon.ico", fs)
	_http.HandleFunc(".*", serveParamless)

	err := http.ListenAndServe(addr, _http)
	if err != nil {
		log.Fatal(err)
	}

}

func serveParamless(w http.ResponseWriter, r *http.Request) {

	var tpl *template.Template
	var err error

	htmlFile := "index"
	if r.URL.Path != "/" {
		htmlFile = filepath.Clean(r.URL.Path)
	}
	htmlTpl := filepath.Join(templateDir, "html", htmlFile+".html")
	if tpl, err = template.ParseFiles(baseTpl, htmlTpl); err != nil {
		log.Fatal(err)
	}
	if err = tpl.ExecuteTemplate(w, "base", nil); err != nil {
		log.Fatal(err)
	}

}
