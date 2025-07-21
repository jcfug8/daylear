package files

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
)

func (s *Service) UploadUserImage(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the request body
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	var body io.Reader = r.Body
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := r.ParseMultipartForm(maxInmemoryUploadSize)
		if err != nil {
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}
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

	name, ok := mux.Vars(r)["name"]
	if !ok {
		http.Error(w, "No user name", http.StatusBadRequest)
		return
	}

	mUser := model.User{}
	_, err = s.userNamer.Parse(name, &mUser)
	if err != nil {
		http.Error(w, "Invalid user name", http.StatusBadRequest)
		return
	}

	imageURI, err := s.domain.UploadUserImage(r.Context(), authAccount, mUser.Id, body)
	if err != nil {
		s.log.Error().Err(err).Msg("unable to upload user image")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(struct {
		ImageURI string `json:"image_uri"`
	}{ImageURI: imageURI})
	w.Write(res)
}
