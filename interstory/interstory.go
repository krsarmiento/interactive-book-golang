package interstory

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
	"html/template"
	"os"
	"bufio"
	"bytes"
)

type Arc struct {
	Title   string
	Story   []string
	Options []map[string]string
}

type Story struct {
	Arcs map[string]Arc
}

type StoryHandler struct {
	Story *Story
	Type  string
}

func (sh StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	sh.RenderArc(path, w)
}

func (sh StoryHandler) RenderArc(arcName string, w http.ResponseWriter) {
	arc, exists := sh.Story.Arcs[arcName]
	if !exists {
		fmt.Println("!! Arc does not exists for: " + arcName)
		return
	}
	renderer := NewRenderer(sh.Type)
	tmpl := renderer.RenderTemplate()
	tmpl.Execute(w, arc)
}

func (sh StoryHandler) RenderArcBytes(arcName string, w *bytes.Buffer) {
	arc, exists := sh.Story.Arcs[arcName]
	if !exists {
		fmt.Println("!! Arc does not exists for: " + arcName)
		return
	}
	renderer := NewRenderer(sh.Type)
	tmpl := renderer.RenderTemplate()
	tmpl.Execute(w, arc)
}

func (sh StoryHandler) RunConsole() {
	currentArc := "intro"
	for {

		var tpl bytes.Buffer
		reader := bufio.NewReader(os.Stdin)
		sh.RenderArcBytes(currentArc, &tpl)
		result := tpl.String()
		fmt.Println(result)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSuffix(answer, "\n")
		fmt.Println(answer)
	}
}

type StoryRenderer interface {
	RenderTemplate() *template.Template
}

type HtmlRenderer struct {}

func (hr HtmlRenderer) RenderTemplate() *template.Template {
	htmlTemplate := `
	<html>
		<h2>{{.Title}}</h2>
		{{range .Story}}
			<p>{{.}}</p>
		{{end}}

		{{if .Options}}
			<h3>Options</h3>
			{{range .Options}}
				<p>
					{{index . "text"}} <a href="{{index . "arc"}}">Go!</a>
				</p>
			{{end}}
		{{else}}
			<a href="/intro">Start over</a>
		{{end}}
	</html>
	`
	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
	return tmpl
}

type ConsoleRenderer struct {}

func (cr ConsoleRenderer) RenderTemplate() *template.Template {
	textTemplate := `{{.Title}}
		{{range .Story}}
			{{.}}\n
		{{end}}
		{{if .Options}}
			Options\n
			{{range .Options}}
				{{index . "text"}} <a href="{{index . "arc"}}">Go!</a>\n
			{{end}}
		{{else}}
			<a href="/intro">Start over</a>\n
		{{end}}
	</html>
	`
	tmpl, err := template.New("text").Parse(textTemplate)
	if err != nil {
		fmt.Println(err)
	}
	return tmpl
}

func NewStoryHandler(fileName string, presentationType string) (*StoryHandler, error) {
	introArcName := "intro"
	storyFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	
	var arcs map[string]Arc
	json.Unmarshal(storyFile, &arcs)

	if _, introExists := arcs[introArcName]; !introExists {
		return nil, errors.New("There's no intro for this story")
	}

	story := &Story{Arcs: arcs}
	storyHandler := &StoryHandler{Story: story, Type: presentationType}
	return storyHandler, nil
}

func NewRenderer(renderType string) StoryRenderer {
	if (renderType == "html") {
		return &HtmlRenderer{}
	}
	if (renderType == "console") {
		return &ConsoleRenderer{}
	}
	return nil
}

func Run(fileName string, presentationType string) {
	storyHandler, err := NewStoryHandler(fileName, presentationType)
	if err != nil {
		panic(err)
	}

	if presentationType == "html" {
		port := "7000"
		mux := http.NewServeMux()
		mux.Handle("/", storyHandler)
		fmt.Println("Listening at port " + port + "...")
		http.ListenAndServe(":" + port, mux)
	}
	if presentationType == "console" {
		storyHandler.RunConsole()
	}
	
}
