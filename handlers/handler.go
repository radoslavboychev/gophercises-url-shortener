package urlshort

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/yaml.v3"
)

type YAMLFile struct {
	Path string
	URL  string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := pathsToUrls[r.URL.Path]
		if url != "" {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
}

func buildMap(yamlFile []YAMLFile) (builtMap map[string]string) {
	builtMap = make(map[string]string)

	for _, yaml := range yamlFile {
		builtMap[yaml.Path] = yaml.URL
	}

	return builtMap
}

func parseYAML(yamlData []byte) (yamlFile []YAMLFile, err error) {
	err = yaml.Unmarshal(yamlData, &yamlFile)
	if err != nil {
		return nil, err
	}
	return yamlFile, nil
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var y YAMLFile
	if err := yaml.Unmarshal(yml, &y); err != nil {
		log.Fatal(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if y.URL != "" {
			http.Redirect(w, r, y.Path, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}), nil
}

func parseJSON(jsonData []byte) (file []YAMLFile, err error) {
	err = json.Unmarshal(jsonData, &file)
	if err != nil {
		return []YAMLFile{}, err
	}
	return file, nil
}

func JSONHandler(jsonData []byte, fallback http.Handler) (jsonHandler http.HandlerFunc, err error) {
	parsedJSON, err := parseJSON(jsonData)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	jsonHandler = MapHandler(pathMap, fallback)
	return
}
