package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

// define a struct for story arcs
type storyArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func parseStory(filename string) (map[string]storyArc, error) {
	// TODO: read json file and return a map
	// the map contain title arc as key
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var story []storyArc
	if err := json.Unmarshal(content, &story); err != nil {
		return nil, err
	}

	storyMap := make(map[string]storyArc)

	for _, arc := range story {
		storyMap[arc.Title] = arc
	}

	return storyMap, nil
}

func executeTemplate(w http.ResponseWriter, tmpl *template.Template, story storyArc) {
	// TODO: execute html and send to the client
	err := tmpl.Execute(w, story)
	if err != nil {
		log.Fatal("could not execute template")
	}
}

// manage the starter request and display the first arc
func startHandler(w http.ResponseWriter, r *http.Request, storyMap map[string]storyArc) {
	story := storyMap["intro"]
	tmpl, err := template.New("").Parse(defaultHandlerTmpl)
	if err != nil {
		log.Fatal("could not parse template")
	}
	executeTemplate(w, tmpl, story)
}

// manage the requestes for netx arcs and display the corresponding arc
func storyHandler() {

}

func main() {
	storyMap, err := parseStory("gopher.json")
	if err != nil {
		log.Fatal(err)
	}

	// http.HandleFunc("/", startHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startHandler(w, r, storyMap)
	})
	http.ListenAndServe(":8080", nil)
}

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </body>
</html>`
