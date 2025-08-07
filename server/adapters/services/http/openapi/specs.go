package openapi

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/jcfug8/daylear/server/openapi"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type spec struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (s *Service) GetSpecs(w http.ResponseWriter, r *http.Request) {
	specs := []spec{}
	blacklist := []string{
		"api/namer",
		"api/types",
	}

	// Walk through the embedded filesystem to find all swagger.json files
	err := fs.WalkDir(openapi.FS, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, blacklisted := range blacklist {
			if strings.Contains(filePath, blacklisted) {
				return nil
			}
		}

		// Check if the file is a swagger.json file
		if !d.IsDir() && strings.HasSuffix(filePath, ".swagger.json") {
			// Extract the name from the path (remove .swagger.json extension)
			baseName := path.Base(strings.TrimSuffix(filePath, ".swagger.json"))
			// Replace underscores with spaces and capitalize each word
			name := cases.Title(language.English).String(strings.ReplaceAll(baseName, "_", " "))
			// Create the URL for the spec
			url := "/openapi/" + filePath

			specs = append(specs, spec{
				Name: name,
				Url:  url,
			})
		}
		return nil
	})

	if err != nil {
		http.Error(w, "Failed to read specs", http.StatusInternalServerError)
		return
	}

	// Set content type and return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(specs)
}
