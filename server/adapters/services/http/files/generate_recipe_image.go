package files

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/rs/zerolog/log"
)

func (s *Service) GenerateRecipeImage(w http.ResponseWriter, r *http.Request) {
	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	name, ok := mux.Vars(r)["name"]
	if !ok {
		http.Error(w, "No recipe name", http.StatusBadRequest)
		return
	}

	mRecipe := model.Recipe{}
	_, err = s.recipeNamer.Parse(name, &mRecipe)
	if err != nil {
		http.Error(w, "Invalid recipe name", http.StatusBadRequest)
		return
	}

	log.Info().Msgf("Generating recipe image for %d", mRecipe.Id)
	file, err := s.domain.GenerateRecipeImage(r.Context(), authAccount, mRecipe.Parent, mRecipe.Id)
	if err != nil {
		http.Error(w, "Failed to generate recipe image", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(file.ContentLength, 10))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)
}
