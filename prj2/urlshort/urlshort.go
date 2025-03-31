package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

type URL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var urls []URL
	if err := yaml.Unmarshal(yml, &urls); err != nil {
		return nil, err
	}
	mappedURLs := buildMap(urls)
	return MapHandler(mappedURLs, fallback), nil
}

func buildMap(paths []URL) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, path := range paths {
		pathsToUrls[path.Path] = path.URL
	}
	return pathsToUrls
}

func JSONHandler(jsonByte []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var urls []URL
	if err := json.Unmarshal(jsonByte, &urls); err != nil {
		return nil, err
	}
	mappedURLs := buildMap(urls)
	return MapHandler(mappedURLs, fallback), nil
}
