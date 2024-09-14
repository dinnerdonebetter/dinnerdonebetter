package main

import (
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestSchemaFromInstance(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		schema := SchemaFromInstance(&types.User{})
		assert.NotNil(t, schema)
	})
}
