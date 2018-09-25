package main

import (
	"fmt"
	"interactive-book-golang/interstory"
	"net/http"
)

func main() {
	port := "7000"
	fileName := "story.json"
	presentationType := "html"
	storyHandler, err := interstory.NewStoryHandler(fileName, presentationType)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", storyHandler)
	fmt.Println("Listening at port " + port + "...")
	http.ListenAndServe(":" + port, mux)
}