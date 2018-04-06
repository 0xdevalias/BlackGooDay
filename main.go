package main

import (
	"io/ioutil"
	"log"
)

// Ref: https://github.com/russross/blackfriday/blob/master/markdown.go#L160-L200

func main() {
	md, err := ioutil.ReadFile("example.md")
	if err != nil {
		log.Fatalf("failed to read markdown file: %v", err)
	}

	docBytes := RunnyBlackGoo(md)

	err = ioutil.WriteFile("example.docx", docBytes,0644)
	if err != nil {
		log.Fatalf("failed to write output document: %v", err)
	}
}
