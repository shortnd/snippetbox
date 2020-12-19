package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr      string
	StaticDir string
}

func main() {
	config := new(Config)
	flag.StringVar(&config.Addr, "addr", "4000", "HTTP network address")
	flag.StringVar(&config.StaticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir(config.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	server := http.Server{
		Addr: ":" + config.Addr,
		ErrorLog: errorLog,
		Handler: mux,
	}
	infoLog.Printf("Starting server on http://localhost:%s", config.Addr)
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
