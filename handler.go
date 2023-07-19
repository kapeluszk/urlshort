package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
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
	pathUrls, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}

	pathsToUrls := buildMapYaml(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseJson(jsn)
	if err != nil {
		return nil, err
	}

	pathsToUrls := buildMapJson(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseJson(data []byte) ([]pathUrlJson, error) {
	var pathUrls []pathUrlJson
	err := json.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func parseYaml(data []byte) ([]pathUrlYaml, error) {
	var pathUrls []pathUrlYaml
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func buildMapYaml(pathUrls []pathUrlYaml) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.Url
	}
	return pathToUrls
}

func buildMapJson(pathUrls []pathUrlJson) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.Url
	}
	return pathToUrls
}

type pathUrlYaml struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

type pathUrlJson struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}
