package files

import (
	"io"
	"net/http"
	"strings"

	"github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *Service) OCRRecipe(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the request body
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	var files []io.Reader
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		// Parse the multipart form
		err := r.ParseMultipartForm(maxInmemoryUploadSize)
		if err != nil {
			http.Error(w, "unable to parse multipart form", http.StatusBadRequest)
			return
		}

		// Get the file from the form
		fileHeaders := r.MultipartForm.File["files"]

		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, "Error reading file", http.StatusBadRequest)
				return
			}
			files = append(files, file)
			defer file.Close()
		}
	} else {
		files = append(files, r.Body)
	}

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	recipe, err := s.domain.OCRRecipe(r.Context(), authAccount, files)
	if err != nil {
		s.log.Error().Err(err).Msg("unable to upload recipe image")
		http.Error(w, "Interal Error", http.StatusInternalServerError)
		return
	}

	pbRecipe, err := convert.RecipeToProto(s.recipeNamer, s.recipeAccessNamer, recipe)
	if err != nil {
		s.log.Error().Err(err).Msg("unable to convert recipe to proto")
		http.Error(w, "Interal Error", http.StatusInternalServerError)
		return
	}

	scrapeRecipeResponse := &pb.ScrapeRecipeResponse{
		Recipe: pbRecipe,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonRecipe, _ := protojson.Marshal(scrapeRecipeResponse)
	w.Write(jsonRecipe)
}
