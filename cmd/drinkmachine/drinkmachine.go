package main

import (
	"flag"
	"github.com/neophenix/drinkmachine/internal/handlers"
	"github.com/neophenix/drinkmachine/internal/hw"
	"github.com/neophenix/drinkmachine/internal/models"
	"github.com/neophenix/drinkmachine/internal/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var port string
var dbFile string
var webroot string
var cacheTemplates bool

func main() {
	flag.StringVar(&port, "port", "80", "port to listen on")
	flag.StringVar(&dbFile, "db", "drinkmachine.db", "location of sqlite db")
	flag.StringVar(&webroot, "webroot", "web", "path of webroot (templates, static, etc)")
	flag.BoolVar(&cacheTemplates, "cache_templates", true, "cache templates or read from disk each time")
	flag.Parse()

	// Log out our "config"
	log.Printf("DB: %v\n", dbFile)
	log.Printf("Listening on port %v\n", port)

	// Channels for signal handline, sigs will get the sig, done is when we are done handling them
	sigs := make(chan os.Signal, 1)

	// Signal handling, since we don't exit we want to catch them and exit cleanly
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("%v\n", sig)
		models.Close()
		hw.ClosePumps()
		hw.CloseDisplay()
		log.Println("Exiting")
		os.Exit(0)
	}()

	template.CacheTemplates = cacheTemplates

	// Open the database
	err := models.Open(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	// setup GPIO + I2C
	hw.InitializePumps()
	hw.InitializeLCD()

	// tell the templates where to look for their files
	template.WebRoot = webroot

	// Static file server
	fs := http.FileServer(http.Dir(webroot + "/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/favicon.ico", fs)

	// The root handler does all the route checking and handoffs
	http.HandleFunc("/", handlers.RootHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
