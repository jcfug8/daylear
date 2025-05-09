package files

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const maxUploadSize = 2 * 1024 * 1024   // 2 MB
const maxInmemoryUploadSize = 10 * 1024 // 100K

func (s *Service) UploadRecipeImage(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the request body
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	var body io.Reader = r.Body
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		// Parse the multipart form
		err := r.ParseMultipartForm(maxInmemoryUploadSize)
		if err != nil {
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		// Get the file from the form
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		body = file
	}

	// Retrieve the auth-token cookie
	headers := r.Header["Authorization"]
	if len(headers) != 1 {
		// For any other error, return a bad request status
		http.Error(w, "No Authorization Token", http.StatusUnauthorized)
		return
	}

	authToken := strings.TrimPrefix(headers[0], "Bearer ")

	name, ok := mux.Vars(r)["name"]
	if !ok {
		http.Error(w, "No recipe name", http.StatusBadRequest)
		return
	}

	parent, id, err := s.recipeNamer.Parse(name)
	if err != nil {
		http.Error(w, "Invalid recipe name", http.StatusBadRequest)
		return
	}

	// Call the domain to check the token
	err = s.domain.AuthorizeRecipeParent(r.Context(), authToken, parent)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	imageURI, err := s.domain.UploadRecipeImage(r.Context(), parent, id, body)
	if err != nil {
		s.log.Error().Err(err).Msg("unable to upload recipe image")
		http.Error(w, "Interal Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(struct {
		ImageURI string `json:"image_uri"`
	}{ImageURI: imageURI})
	w.Write(res)
}
