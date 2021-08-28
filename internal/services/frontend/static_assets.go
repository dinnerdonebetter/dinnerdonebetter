package frontend

import (
	// for embedding assets.
	_ "embed"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

func (s *service) buildStaticFileServer() *afero.HttpFs {
	return afero.NewHttpFs(afero.NewOsFs())
}

var (
	// Here is where you should put route regexes that need to be ignored by the static file server.
	// For instance, if you allow someone to see an event in the frontend via a URL that contains dynamic.
	// information, such as `/event/123`, you would want to put something like this below:
	// 		eventsFrontendPathRegex = regexp.MustCompile(`/event/\d+`)

	// validIngredientAdminPathRegex matches URLs against our frontend router's specification for specific valid ingredient routes.
	validIngredientAdminPathRegex = regexp.MustCompile(`/admin/valid_ingredients/\d+`)

	// validInstrumentAdminPathRegex matches URLs against our frontend router's specification for specific valid instrument routes.
	validInstrumentAdminPathRegex = regexp.MustCompile(`/admin/valid_instruments/\d+`)

	// validPreparationAdminPathRegex matches URLs against our frontend router's specification for specific valid preparation routes.
	validPreparationAdminPathRegex = regexp.MustCompile(`/admin/valid_preparations/\d+`)

	// recipeAppPathRegex matches URLs against our frontend router's specification for specific recipe routes.
	recipeAppPathRegex = regexp.MustCompile(`/recipes/\d+`)
)

// StaticDir builds a static directory handler.
func (s *service) StaticDir(staticFilesDirectory string) (http.HandlerFunc, error) {
	fileDir, err := filepath.Abs(staticFilesDirectory)
	if err != nil {
		return nil, fmt.Errorf("determining absolute path of static files directory: %w", err)
	}

	httpFs := s.buildStaticFileServer()

	s.logger.WithValue("static_dir", fileDir).Debug("setting static file server")
	fs := http.StripPrefix("/", http.FileServer(httpFs.Dir(fileDir)))

	return func(res http.ResponseWriter, req *http.Request) {
		logger := s.logger.WithRequest(req)
		logger.Debug("static file requested")

		sessCtxData, sessCtxErr := s.sessionContextDataFetcher(req)
		if sessCtxErr != nil {
			logger.Error(sessCtxErr, "fetching session context data")
		}

		if strings.HasPrefix(req.URL.Path, "/admin") && sessCtxData != nil && !sessCtxData.ServiceRolePermissionChecker().IsServiceAdmin() {
			res.WriteHeader(http.StatusUnauthorized)
			http.Redirect(res, req, "/login", http.StatusUnauthorized)
			return
		}

		switch req.URL.Path {
		// list your frontend history routes here.
		case "/register",
			"/login",
			"/home",
			"/recipes",
			"/plans",
			"/household",
			"/admin",
			"/admin/dashboard",
			"/admin/users",
			"/admin/recipes/new",
			"/admin/valid_ingredients",
			"/admin/valid_ingredients/new",
			"/admin/valid_instruments",
			"/admin/valid_instruments/new",
			"/admin/valid_preparations",
			"/admin/valid_preparations/new",
			"/admin/recipes":
			logger.Debug("rerouting")
			req.URL.Path = "/"
		}

		if validIngredientAdminPathRegex.MatchString(req.URL.Path) {
			logger.Debug("rerouting admin valid ingredient request")
			req.URL.Path = "/"
		}

		if validInstrumentAdminPathRegex.MatchString(req.URL.Path) {
			logger.Debug("rerouting admin valid instrument request")
			req.URL.Path = "/"
		}

		if validPreparationAdminPathRegex.MatchString(req.URL.Path) {
			logger.Debug("rerouting admin valid preparation request")
			req.URL.Path = "/"
		}

		if recipeAppPathRegex.MatchString(req.URL.Path) {
			logger.Debug("rerouting app individual recipe request")
			req.URL.Path = "/"
		}

		logger.WithValue("destination", req.URL.Path).Debug("heading to frontend path")

		fs.ServeHTTP(res, req)
	}, nil
}
