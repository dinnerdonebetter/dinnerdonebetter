package frontend

import (
	"gitlab.com/prixfixe/prixfixe/internal/routing"
)

// SetupRoutes sets up the routes.
func (s *service) SetupRoutes(router routing.Router) {
	router = router.WithMiddleware(s.authService.UserAttributionMiddleware)

	staticFileServer, err := s.StaticDir("/frontend")
	if err != nil {
		s.logger.Error(err, "establishing static file server")
	}
	router.Get("/*", staticFileServer)
}
