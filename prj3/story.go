package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
)

type Chapters struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type StoryHandler struct {
	chapters map[string]Chapters
	tmpl     *template.Template
}

func NewStoryHandler(inputFile string) (*StoryHandler, error) {

	jsonByte, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("can not open input file %v", err)
	}

	var chapters map[string]Chapters
	if err := json.Unmarshal(jsonByte, &chapters); err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles("layout.html")
	if err != nil {
		return nil, err
	}

	return &StoryHandler{
		chapters: chapters,
		tmpl:     tmpl,
	}, nil
}
