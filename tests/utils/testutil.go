package testutils

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/pkg/types"
)

func init() {
	fake.Seed(time.Now().UnixNano())
}

// errArbitrary is an arbitrary error.
var errArbitrary = errors.New("blah")

// BrokenSessionContextDataFetcher is a deliberately broken sessionContextDataFetcher.
func BrokenSessionContextDataFetcher(*http.Request) (*types.SessionContextData, error) {
	return nil, errArbitrary
}

// BuildArbitraryImage builds an image with a bunch of colors in it.
func BuildArbitraryImage(widthAndHeight int) image.Image {
	img := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: widthAndHeight, Y: widthAndHeight}})

	// Set color for each pixel.
	for x := 0; x < widthAndHeight; x++ {
		for y := 0; y < widthAndHeight; y++ {
			img.Set(x, y, color.RGBA{R: uint8(x % math.MaxUint8), G: uint8(y % math.MaxUint8), B: uint8(x + y%math.MaxUint8), A: math.MaxUint8})
		}
	}

	return img
}

// BuildArbitraryImagePNGBytes builds an image with a bunch of colors in it.
func BuildArbitraryImagePNGBytes(widthAndHeight int) []byte {
	var b bytes.Buffer
	if err := png.Encode(&b, BuildArbitraryImage(widthAndHeight)); err != nil {
		panic(err)
	}

	return b.Bytes()
}

// BuildTestRequest builds an arbitrary *http.Request.
func BuildTestRequest(t *testing.T) *http.Request {
	t.Helper()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodOptions,
		"https://prixfixe.verygoodsoftwarenotvirus.ru",
		http.NoBody,
	)

	require.NotNil(t, req)
	assert.NoError(t, err)

	return req
}
