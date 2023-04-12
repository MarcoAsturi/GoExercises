package main

import (
	"fmt"
	"log"
	"net/http"
)

// manages http request
// if the request path is present in the map, it will redirect to the repective page
// otherwise it will pass a fallback request
type pathRedirectHandler struct {
	pathMap  map[string]string // mapping of paths to URLs
	fallback http.Handler
}

// implements ServeHTTP method to check if the request is present in the map, imlpementing http.Handler interface
// if the request path is present, it will redirect to the corresponding url
func (prh *pathRedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if redirectTo, ok := prh.pathMap[r.URL.Path]; ok {
		http.Redirect(w, r, redirectTo, http.StatusFound)
		return
	}

	prh.fallback.ServeHTTP(w, r)
}

// creates a new pathRedirectHandler from the checked map
func newPathRedirectHandler(pathMap map[string]string, fallback http.Handler) *pathRedirectHandler {
	return &pathRedirectHandler{
		pathMap:  pathMap,
		fallback: fallback,
	}
}

// handler that manage the redirect map
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.Handler {
	pathRedirect := newPathRedirectHandler(pathsToUrls, fallback)
	return pathRedirect
}

func main() {
	// creating server mux
	mux := defaultMux()

	// creation of redirect map
	pathsToUrls := map[string]string{
		"/dogs":  "https://www.tsmean.com/articles/mustache/the-ultimate-mustache-tutorial/",
		"/cats":  "https://courses.calhoun.io/lessons/les_goph_05",
		"/birds": "https://stackoverflow.com/questions/46311740/parsing-html-with-go",
	}

	// creation of the handler that manage the redirect map
	// mapHandler := MapHandler(pathsToUrls, mux)
	http.Handle("/", MapHandler(pathsToUrls, mux))

	// mux.Handle("/", mapHandler)

	// registartion of each path of the map - USELESS?
	// for path, redirectTo := range pathsToUrls {
	// 	mux.Handle(path, http.RedirectHandler(redirectTo, http.StatusFound))
	// }

	// starting server
	fmt.Println(("Server starting on 8080..."))
	// just a check for display the handler
	fmt.Printf("Handler: %v\n", mux)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func for creation of default handler
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", emptyPage)
	return mux
}

// handler for default page
func emptyPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Write a short url next to \"localhost:8080\"\nFor example \"localhost:8080/dogs\"")
}
