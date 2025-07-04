package files

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
)

func (s *Service) OCRRecipe(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the request body
	// r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

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

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	recipe, err := s.domain.OCRRecipe(r.Context(), authAccount, body)
	if err != nil {
		s.log.Error().Err(err).Msg("unable to upload recipe image")
		http.Error(w, "Interal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(struct {
		Recipe model.Recipe `json:"recipe"`
	}{Recipe: recipe})

	w.Write(res)
}
