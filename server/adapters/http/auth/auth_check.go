package auth

import (
	"net/http"

	"github.com/jcfug8/daylear/server/adapters/http/headers"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *Service) AuthCheck(w http.ResponseWriter, r *http.Request) {
	authToken, err := headers.GetAuthToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Call the domain to check the token
	user, err := s.domain.ParseToken(r.Context(), authToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	name, err := s.userNamer.Format(user.Id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pbUser := &pb.User{
		Name: name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(protojson.Format(pbUser)))
}
