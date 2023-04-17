package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

// Chapter rappresents a single chapter of the story
type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

// Story represents the entire story
// it is a map: the key is the title chapter and the value is the chapter
type Story map[string]Chapter

// load the json file
func loadStory(filename string) (Story, error) {
	// open file inb reading mode
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close() // REMEMBER to close it

	// decoding the json content in a Chapter map
	var story Story
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&story); err != nil { // dont FORGET pointers
		return nil, err
	}
	return story, nil
}

// handle HTTP request for a specific chapter
func chapterHandler(story Story) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// remove space and "/"  from HTTP request
		path := strings.TrimSpace(r.URL.Path)
		if path == "" || path == "/" {
			path = "/intro"
		} else {
			path = path[1:]
		}

		// if the request chapter is present in the map, return it using html template
		if chapter, ok := story[path]; ok {
			t := template.Must(template.ParseFiles("template.html")) // USE .Must
			if err := t.Execute(w, chapter); err != nil {
				log.Printf("template error: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		} else {
			http.NotFound(w, r)
		}
	}
}

func main() {
	// load story
	story, err := loadStory("gopher.json")
	if err != nil {
		log.Fatal(err)
	}

	// setting up the handler
	http.HandleFunc("/", chapterHandler(story))

	port := "8080"
	log.Printf("listening on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
