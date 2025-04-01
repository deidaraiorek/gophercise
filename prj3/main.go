package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	handler, err := NewStoryHandler("gopher.json")
	if err != nil {
		log.Fatalf("Failed to initialize story: %v", err)
	}

	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", handler))

}

func (handler *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chapter := filepath.Base(r.URL.Path)
	if chapter == "" || chapter == "/" || chapter == "\\" {
		chapter = "intro"
	}
	handler.display(chapter, w, r)
}

func (handler *StoryHandler) display(chapter string, w http.ResponseWriter, r *http.Request) {
	story, err := handler.chapters[chapter]
	if !err {
		http.NotFound(w, r)
		log.Printf("Chapter not found: %s", chapter)
		return
	}

	if err := handler.tmpl.Execute(w, story); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}
