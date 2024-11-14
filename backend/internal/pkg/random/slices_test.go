package random

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElement(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleArray := []string{
			"The",
			"FitnessGram™",
			"Pacer",
			"Test",
			"is",
			"a",
			"multistage",
			"aerobic",
			"capacity",
			"test",
			"that",
			"progressively",
			"gets",
			"more",
			"difficult",
			"as",
			"it",
			"continues.",
			"The",
			"20",
			"meter",
			"pacer",
			"test",
			"will",
			"begin",
			"in",
			"30",
			"seconds.",
			"Line",
			"up",
			"at",
			"the",
			"start.",
			"The",
			"running",
			"speed",
			"starts",
			"slowly,",
			"but",
			"gets",
			"faster",
			"each",
			"minute",
			"after",
			"you",
			"hear",
			"this",
			"signal.",
		}

		assert.True(t, slices.Contains(exampleArray, Element(exampleArray)))
	})
}
