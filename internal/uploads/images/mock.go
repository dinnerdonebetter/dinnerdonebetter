package images

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
)

var _ ImageUploadProcessor = (*MockImageUploadProcessor)(nil)

// MockImageUploadProcessor is a mock ImageUploadProcessor.
type MockImageUploadProcessor struct {
	mock.Mock
}

// ProcessFile satisfies the ImageUploadProcessor interface.
func (m *MockImageUploadProcessor) ProcessFile(ctx context.Context, req *http.Request, filename string) (*Upload, error) {
	args := m.Called(ctx, req, filename)

	return args.Get(0).(*Upload), args.Error(1)
}

func (m *MockImageUploadProcessor) ProcessFiles(ctx context.Context, req *http.Request, filenamePrefix string) ([]*Upload, error) {
	args := m.Called(ctx, req, filenamePrefix)

	return args.Get(0).([]*Upload), args.Error(1)
}
