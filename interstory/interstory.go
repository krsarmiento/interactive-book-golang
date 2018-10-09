package interstory

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"io"
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

func (sh StoryHandler) RenderArc(arcName string, w io.Writer) error {
	arc, exists := sh.Story.Arcs[arcName]
	if !exists {
		return errors.New("!! Arc does not exists for: " + arcName)
	}
	renderer := NewRenderer(sh.Type)
	tmpl := renderer.RenderTemplate()
	tmpl.Execute(w, arc)
	return nil
}

func (sh StoryHandler) RunConsole() {
	currentArc := "intro"
	oldArc := "intro"
	var tpl bytes.Buffer
	reader := bufio.NewReader(os.Stdin)
	for {
		tpl.Reset()
		err := sh.RenderArc(currentArc, &tpl)
		if err != nil {
			fmt.Println(err)
			currentArc = oldArc
			continue
		}
		arcText := tpl.String()
		fmt.Println("********************************************")
		fmt.Println(arcText)
		fmt.Println("********************************************")
		fmt.Print("Write the selection and press Enter >> ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSuffix(answer, "\n")
		oldArc = currentArc
		currentArc = answer
	}
}

type StoryRenderer interface {
	RenderTemplate() *template.Template
}

type HtmlRenderer struct {}

func (hr HtmlRenderer) RenderTemplate() *template.Template {
	htmlTemplate, err := ioutil.ReadFile("template.html")
	if err != nil {
		return nil
	}

	tmpl, err := template.New("html").Parse(string(htmlTemplate))
	if err != nil {
		panic(err)
	}
	return tmpl
}

type ConsoleRenderer struct {}

func (cr ConsoleRenderer) RenderTemplate() *template.Template {
	textTemplate, err := ioutil.ReadFile("template.txt")
	if err != nil {
		return nil
	}

	tmpl, err := template.New("text").Parse(string(textTemplate))
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
