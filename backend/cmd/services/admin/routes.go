package main

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/platform/routing"

	ghttp "maragu.dev/gomponents/http"
)

func (s *AdminFrontendServer) setupRoutes(router routing.Router) {
	r := router.WithMiddleware(s.authMiddleware)

	r.Get("/", ghttp.Adapt(s.homeRoute))

	r.Get(fmt.Sprintf("/users/{%s}", userIDURLParamKey), ghttp.Adapt(s.UserPage))
	r.Get("/users", ghttp.Adapt(s.UsersList))
	r.Get("/api/users/search", ghttp.Adapt(s.UsersSearch))

	router.Get("/login", ghttp.Adapt(s.LoginPage))
	router.Post("/login/submit", ghttp.Adapt(s.LoginSubmission))
}
