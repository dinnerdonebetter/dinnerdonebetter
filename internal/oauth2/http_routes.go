package oauth2

import (
	"net/http"
)

// AuthorizeHandler is our oauth2 auth route.
func (s *Service) AuthorizeHandler(res http.ResponseWriter, req *http.Request) {
	if err := s.oauth2Server.HandleAuthorizeRequest(res, req); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func (s *Service) TokenHandler(res http.ResponseWriter, req *http.Request) {
	if err := s.oauth2Server.HandleTokenRequest(res, req); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
