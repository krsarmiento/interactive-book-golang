package main

import (
	"fmt"
	"interactive-book-golang/interstory"
)

func main() {
	story, err := interstory.NewStory("story.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(story)
}