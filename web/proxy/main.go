package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var (
	appMux    = http.NewServeMux()
	vendorMux = http.NewServeMux()
	apiMux    *httputil.ReverseProxy

	devMode = flag.Bool("dev", true, "Activate dev mode")
)

const (
	API_PREFIX = "/api"
)

func main() {
	flag.Parse()

	apiUrl, err := url.Parse("http://boot2docker:3000")
	if err != nil {
		log.Fatal("Bad api url", err)
	}

	apiMux = httputil.NewSingleHostReverseProxy(apiUrl)
	// apiMux.Director = func(r *http.Request) {
	// 	fmt.Printf("Got req %v", r)
	// }

	if *devMode {
		fmt.Printf("dev mode\n")
		appMux.Handle("/", http.FileServer(http.Dir("../src/")))
		vendorMux.Handle("/", http.FileServer(http.Dir("../")))
	} else {
		fmt.Printf("build mode\n")
		appMux.Handle("/", http.FileServer(http.Dir("../build/")))
	}

	http.HandleFunc("/", router)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func router(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.RequestURI, API_PREFIX) {
		fmt.Printf("[API ] %v\n", r.RequestURI)
		r.RequestURI = r.RequestURI[len(API_PREFIX):]
		r.URL.Path = r.URL.Path[len(API_PREFIX):]
		apiMux.ServeHTTP(w, r)
		return
	}

	fmt.Printf("[FILE] %v\n", r.RequestURI)
	if *devMode {
		if strings.HasPrefix(r.RequestURI, "/vendor") || strings.HasPrefix(r.RequestURI, "/build") {
			vendorMux.ServeHTTP(w, r)
			return
		}
	}
	appMux.ServeHTTP(w, r)
}
