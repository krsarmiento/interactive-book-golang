package main

import (
	"interactive-book-golang/interstory"
)

func main() {
	fileName := "story.json"
	presentationType := "txt"
	interstory.Run(fileName, presentationType)
}