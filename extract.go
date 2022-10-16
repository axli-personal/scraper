package main

import (
	"log"
	"strings"
)

func extractAuthor(text string) string {
	array := strings.Split(text, "/")
	if len(array) == 0 {
		log.Print("Fail to extract author:", text)
		return "#"
	}
	return array[0]
}

func extractBookId(path string) string {
	array := strings.Split(path, "/")
	if len(array) <= 2 {
		log.Print("Fail to extract book ID:", path)
		return "0"
	}
	return array[2]
}
