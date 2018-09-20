package main

import (
	"fmt"
	"interactive-book-golang/interstory"
	"net/http"
)

func main() {
	port := "7000"
	story, err := interstory.NewStory("story.json")
	if err != nil {
		panic(err)
	}

	storyHandler := interstory.StoryHandler{Story: story}
	mux := http.NewServeMux()
	mux.Handle("/", storyHandler)
	fmt.Println("Listening at port " + port + "...")
	http.ListenAndServe(":" + port, mux)
}