package main

import (
	"fmt"
	"html/template"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/chrj/placetaytay/taytay"

	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
)

var cacheLock = &sync.Mutex{}

var t = template.New("")

func main() {

	var err error

	if err = taytay.Index("taytay"); err != nil {
		log.Fatalf("couldn't index pictures: %v", err)
	}

	if t, err = template.ParseGlob(
		filepath.Join("taytay", "*.html"),
	); err != nil {
		log.Fatalf("couldn't parse templates: %v", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", AboutHandler)
	r.PathPrefix("/.well-known/").Handler(http.StripPrefix("/.well-known/", http.FileServer(http.Dir("taytay/.well-known"))))
	r.HandleFunc("/random", RandomHandler)
	r.HandleFunc("/{width:[0-9]+}x{height:[0-9]+}", ImageHandler)

	log.Fatal(http.ListenAndServe("127.0.0.1:8060", r))

}

func AboutHandler(rw http.ResponseWriter, req *http.Request) {

	f, err := os.Open(filepath.Join("taytay", "index.html"))

	if err != nil {
		http.Error(rw, "Couldn't read index.html", http.StatusInternalServerError)
		log.Printf("error: %v", err)
		return
	}

	defer f.Close()

	io.Copy(rw, f)

}

func RandomHandler(rw http.ResponseWriter, req *http.Request) {

	sizes := [10]string{}
	for i := 0; i < 10; i++ {

		w, h := 1, 1000
		ratio := 0.0

		for ratio < 0.5 || ratio > 1.5 {
			w, h = (rand.Intn(100)+5)*10, (rand.Intn(100)+5)*10
			ratio = float64(w) / float64(h)
		}

		sizes[i] = fmt.Sprintf("%dx%d", w, h)

	}

	if err := t.ExecuteTemplate(rw, "random.html", sizes); err != nil {
		http.Error(rw, "Couldn't render random.html", http.StatusInternalServerError)
		log.Printf("error: %v", err)
		return
	}

}

func ImageHandler(rw http.ResponseWriter, req *http.Request) {

	var (
		width, height int
	)

	v := mux.Vars(req)

	if i, err := strconv.ParseInt(v["width"], 10, strconv.IntSize); err != nil || i <= 0 || i > 2048 {
		http.Error(rw, "Couldn't parse width", http.StatusBadRequest)
		return
	} else {
		width = int(i)
	}

	if i, err := strconv.ParseInt(v["height"], 10, strconv.IntSize); err != nil || i <= 0 || i > 2048 {
		http.Error(rw, "Couldn't parse height", http.StatusBadRequest)
		return
	} else {
		height = int(i)
	}

	rw.Header().Set("Content-Type", "image/png")

	cacheLock.Lock()

	cacheFilename := filepath.Join(
		"taytay",
		"cache",
		fmt.Sprintf("%dx%d.png", width, height),
	)

	_, err := os.Stat(cacheFilename)

	cacheLock.Unlock()

	if err == nil {
		http.ServeFile(rw, req, cacheFilename)
		return
	}

	ts := taytay.Closest(width, height)

	im := imaging.Fill(
		ts.Image(),
		width,
		height,
		imaging.Center,
		imaging.MitchellNetravali,
	)

	cacheLock.Lock()
	defer cacheLock.Unlock()

	f, err := os.OpenFile(cacheFilename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		http.Error(rw, "Couldn't open cache file", http.StatusInternalServerError)
		log.Printf("error: %v", err)
		return
	}

	defer f.Close()

	if err := png.Encode(f, im); err != nil {
		http.Error(rw, "Couldn't encode image", http.StatusInternalServerError)
		log.Printf("error: %v", err)
		return
	}

	f.Close()

	log.Printf("cached: %s", cacheFilename)

	http.ServeFile(rw, req, cacheFilename)

}
