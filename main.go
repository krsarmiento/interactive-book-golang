package main

import (
	"interactive-book-golang/interstory"
)

func main() {
	fileName := "story.json"
	presentationType := "console"
	interstory.Run(fileName, presentationType)
}