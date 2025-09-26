package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvideHTTPServer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x, err := ProvideHTTPServer(
			Config{
				HTTPSCertificateFile:    "",
				HTTPSCertificateKeyFile: "",
				StartupDeadline:         0,
				HTTPPort:                0,
				Debug:                   false,
			},
			nil,
			nil,
			nil,
		)

		assert.NotNil(t, x)
		assert.NoError(t, err)
	})
}
