package interstory

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Arc struct {
	Title   string
	Story   []string
	Options []map[string]string
}

func NewStory(file string) (map[string]Arc, error) {
	introArcName := "intro"
	storyFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var story map[string]Arc
	json.Unmarshal(storyFile, &story)

	if _, introExists := story[introArcName]; !introExists {
		return nil, errors.New("There's no intro for this story")
	}

	return story, nil
}