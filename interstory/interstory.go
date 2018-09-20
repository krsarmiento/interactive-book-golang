package interstory

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
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
}

func (sh StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	path := strings.TrimLeft(r.URL.Path, "/")
	arc, exists := sh.Story.Arcs[path]
	if exists {
		fmt.Println(arc.Title)
		fmt.Fprintln(w, arc.Title)
	}
}

func NewStory(file string) (*Story, error) {
	introArcName := "intro"
	storyFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	
	var arcs map[string]Arc
	json.Unmarshal(storyFile, &arcs)

	if _, introExists := arcs[introArcName]; !introExists {
		return nil, errors.New("There's no intro for this story")
	}

	story := &Story{Arcs: arcs}
	return story, nil
}