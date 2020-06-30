package frontend

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/spf13/afero"
)

// Routes returns a map of route to HandlerFunc for the parent router to set.
// this keeps routing logic in the frontend service and not in the server itself.
func (s *Service) Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		// "/login":    s.LoginPage,
		// "/register": s.RegistrationPage,
	}
}

func (s *Service) buildStaticFileServer(fileDir string) (*afero.HttpFs, error) {
	var afs afero.Fs
	if s.config.CacheStaticFiles {
		afs = afero.NewMemMapFs()
		files, err := ioutil.ReadDir(fileDir)
		if err != nil {
			return nil, fmt.Errorf("reading directory for frontend files: %w", err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			fp := filepath.Join(fileDir, file.Name())
			f, err := afs.Create(fp)
			if err != nil {
				return nil, fmt.Errorf("creating static file in memory: %w", err)
			}

			bs, err := ioutil.ReadFile(fp)
			if err != nil {
				return nil, fmt.Errorf("reading static file from directory: %w", err)
			}

			if _, err = f.Write(bs); err != nil {
				return nil, fmt.Errorf("loading static file into memory: %w", err)
			}

			if err = f.Close(); err != nil {
				s.logger.Error(err, "closing file while setting up static dir")
			}
		}
		afs = afero.NewReadOnlyFs(afs)
	} else {
		afs = afero.NewOsFs()
	}

	return afero.NewHttpFs(afs), nil
}

var (
	// Here is where you should put route regexes that need to be ignored by the static file server.
	// For instance, if you allow someone to see an event in the frontend via a URL that contains dynamic.
	// information, such as `/event/123`, you would want to put something like this below:
	// 		eventsFrontendPathRegex = regexp.MustCompile(`/event/\d+`)

	// validInstrumentsFrontendPathRegex matches URLs against our frontend router's specification for specific valid instrument routes.
	validInstrumentsFrontendPathRegex = regexp.MustCompile(`/valid_instruments/\d+`)
	// validIngredientsFrontendPathRegex matches URLs against our frontend router's specification for specific valid ingredient routes.
	validIngredientsFrontendPathRegex = regexp.MustCompile(`/valid_ingredients/\d+`)
	// validIngredientTagsFrontendPathRegex matches URLs against our frontend router's specification for specific valid ingredient tag routes.
	validIngredientTagsFrontendPathRegex = regexp.MustCompile(`/valid_ingredient_tags/\d+`)
	// ingredientTagMappingsFrontendPathRegex matches URLs against our frontend router's specification for specific ingredient tag mapping routes.
	ingredientTagMappingsFrontendPathRegex = regexp.MustCompile(`/ingredient_tag_mappings/\d+`)
	// validPreparationsFrontendPathRegex matches URLs against our frontend router's specification for specific valid preparation routes.
	validPreparationsFrontendPathRegex = regexp.MustCompile(`/valid_preparations/\d+`)
	// requiredPreparationInstrumentsFrontendPathRegex matches URLs against our frontend router's specification for specific required preparation instrument routes.
	requiredPreparationInstrumentsFrontendPathRegex = regexp.MustCompile(`/required_preparation_instruments/\d+`)
	// validIngredientPreparationsFrontendPathRegex matches URLs against our frontend router's specification for specific valid ingredient preparation routes.
	validIngredientPreparationsFrontendPathRegex = regexp.MustCompile(`/valid_ingredient_preparations/\d+`)
	// recipesFrontendPathRegex matches URLs against our frontend router's specification for specific recipe routes.
	recipesFrontendPathRegex = regexp.MustCompile(`/recipes/\d+`)
	// recipeTagsFrontendPathRegex matches URLs against our frontend router's specification for specific recipe tag routes.
	recipeTagsFrontendPathRegex = regexp.MustCompile(`/recipe_tags/\d+`)
	// recipeStepsFrontendPathRegex matches URLs against our frontend router's specification for specific recipe step routes.
	recipeStepsFrontendPathRegex = regexp.MustCompile(`/recipe_steps/\d+`)
	// recipeStepPreparationsFrontendPathRegex matches URLs against our frontend router's specification for specific recipe step preparation routes.
	recipeStepPreparationsFrontendPathRegex = regexp.MustCompile(`/recipe_step_preparations/\d+`)
	// recipeStepIngredientsFrontendPathRegex matches URLs against our frontend router's specification for specific recipe step ingredient routes.
	recipeStepIngredientsFrontendPathRegex = regexp.MustCompile(`/recipe_step_ingredients/\d+`)
	// recipeIterationsFrontendPathRegex matches URLs against our frontend router's specification for specific recipe iteration routes.
	recipeIterationsFrontendPathRegex = regexp.MustCompile(`/recipe_iterations/\d+`)
	// recipeIterationStepsFrontendPathRegex matches URLs against our frontend router's specification for specific recipe iteration step routes.
	recipeIterationStepsFrontendPathRegex = regexp.MustCompile(`/recipe_iteration_steps/\d+`)
	// iterationMediasFrontendPathRegex matches URLs against our frontend router's specification for specific iteration media routes.
	iterationMediasFrontendPathRegex = regexp.MustCompile(`/iteration_medias/\d+`)
	// invitationsFrontendPathRegex matches URLs against our frontend router's specification for specific invitation routes.
	invitationsFrontendPathRegex = regexp.MustCompile(`/invitations/\d+`)
	// reportsFrontendPathRegex matches URLs against our frontend router's specification for specific report routes.
	reportsFrontendPathRegex = regexp.MustCompile(`/reports/\d+`)
)

// StaticDir builds a static directory handler.
func (s *Service) StaticDir(staticFilesDirectory string) (http.HandlerFunc, error) {
	fileDir, err := filepath.Abs(staticFilesDirectory)
	if err != nil {
		return nil, fmt.Errorf("determining absolute path of static files directory: %w", err)
	}

	httpFs, err := s.buildStaticFileServer(fileDir)
	if err != nil {
		return nil, fmt.Errorf("establishing static server filesystem: %w", err)
	}

	s.logger.WithValue("static_dir", fileDir).Debug("setting static file server")
	fs := http.StripPrefix("/", http.FileServer(httpFs.Dir(fileDir)))

	return func(res http.ResponseWriter, req *http.Request) {
		rl := s.logger.WithRequest(req)
		rl.Debug("static file requested")
		switch req.URL.Path {
		// list your frontend history routes here.
		case "/register",
			"/login",
			"/admin",
			"/admin/dashboard",
			"/admin/enumerations/valid_instruments",
			"/admin/enumerations/valid_instruments/new",
			"/admin/enumerations/valid_ingredients",
			"/admin/enumerations/valid_ingredients/new",
			"/admin/enumerations/valid_ingredient_tags",
			"/admin/enumerations/valid_ingredient_tags/new",
			"/admin/enumerations/valid_ingredient_preparations",
			"/admin/enumerations/valid_ingredient_preparations/new",
			"/admin/enumerations/valid_preparations",
			"/admin/enumerations/valid_preparations/new",
			"/ingredient_tag_mappings",
			"/ingredient_tag_mappings/new",
			"/required_preparation_instruments",
			"/required_preparation_instruments/new",
			"/recipes",
			"/recipes/new",
			"/recipe_tags",
			"/recipe_iterations",
			"/recipe_iterations/new",
			"/iteration_medias",
			"/iteration_medias/new",
			"/invitations",
			"/invitations/new",
			"/reports",
			"/reports/new",
			"/password/new":
			rl.Debug("rerouting")
			req.URL.Path = "/"
		}
		if validInstrumentsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting valid instrument request")
			req.URL.Path = "/"
		}
		if validIngredientsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting valid ingredient request")
			req.URL.Path = "/"
		}
		if validIngredientTagsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting valid ingredient tag request")
			req.URL.Path = "/"
		}
		if ingredientTagMappingsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting ingredient tag mapping request")
			req.URL.Path = "/"
		}
		if validPreparationsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting valid preparation request")
			req.URL.Path = "/"
		}
		if requiredPreparationInstrumentsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting required preparation instrument request")
			req.URL.Path = "/"
		}
		if validIngredientPreparationsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting valid ingredient preparation request")
			req.URL.Path = "/"
		}
		if recipesFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting recipe request")
			req.URL.Path = "/"
		}
		if recipeTagsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting recipe tag request")
			req.URL.Path = "/"
		}
		if recipeStepsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting recipe step request")
			req.URL.Path = "/"
		}
		if recipeStepPreparationsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting recipe step preparation request")
			req.URL.Path = "/"
		}
		if recipeStepIngredientsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting recipe step ingredient request")
			req.URL.Path = "/"
		}
		if recipeIterationsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting recipe iteration request")
			req.URL.Path = "/"
		}
		if recipeIterationStepsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting recipe iteration step request")
			req.URL.Path = "/"
		}
		if iterationMediasFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting iteration media request")
			req.URL.Path = "/"
		}
		if invitationsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting invitation request")
			req.URL.Path = "/"
		}
		if reportsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting report request")
			req.URL.Path = "/"
		}

		fs.ServeHTTP(res, req)
	}, nil
}
