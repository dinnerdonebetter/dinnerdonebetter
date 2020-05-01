package frontend

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
)

func buildRequest(t *testing.T) *http.Request {
	t.Helper()

	req, err := http.NewRequest(
		http.MethodGet,
		"https://verygoodsoftwarenotvirus.ru",
		nil,
	)

	require.NotNil(t, req)
	assert.NoError(t, err)
	return req
}

func TestService_StaticDir(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}

		cwd, err := os.Getwd()
		require.NoError(t, err)

		hf, err := s.StaticDir(cwd)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/http_routes_test.go"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/login"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend valid instruments routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/valid_instruments/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend valid ingredients routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/valid_ingredients/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend valid ingredient tags routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/valid_ingredient_tags/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend ingredient tag mappings routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/ingredient_tag_mappings/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend valid preparations routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/valid_preparations/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend required preparation instruments routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/required_preparation_instruments/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend valid ingredient preparations routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/valid_ingredient_preparations/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend recipes routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/recipes/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend recipe tags routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/recipe_tags/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend recipe steps routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/recipe_steps/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend recipe step preparations routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/recipe_step_preparations/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend recipe step ingredients routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/recipe_step_ingredients/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend recipe iterations routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/recipe_iterations/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend recipe iteration steps routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/recipe_iteration_steps/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend iteration medias routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/iteration_medias/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend invitations routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/invitations/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend reports routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/reports/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func TestService_Routes(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		assert.NotNil(t, (&Service{}).Routes())
	})
}

func TestService_buildStaticFileServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		s := &Service{
			config: config.FrontendSettings{
				CacheStaticFiles: true,
			},
		}
		cwd, err := os.Getwd()
		require.NoError(t, err)

		actual, err := s.buildStaticFileServer(cwd)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})
}
