package mockencoding

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

// ClientEncoder is a mock ClientEncoder.
type ClientEncoder struct {
	mock.Mock
}

// ContentType satisfies the ClientEncoder interface.
func (m *ClientEncoder) ContentType() string {
	return m.Called().String(0)
}

// Unmarshal satisfies the ClientEncoder interface.
func (m *ClientEncoder) Unmarshal(ctx context.Context, data []byte, v any) error {
	return m.Called(ctx, data, v).Error(0)
}

// Encode satisfies the ClientEncoder interface.
func (m *ClientEncoder) Encode(ctx context.Context, dest io.Writer, v any) error {
	return m.Called(ctx, dest, v).Error(0)
}

// EncodeReader satisfies the ClientEncoder interface.
func (m *ClientEncoder) EncodeReader(ctx context.Context, data any) (io.Reader, error) {
	returnValues := m.Called(ctx, data)

	return returnValues.Get(0).(io.Reader), returnValues.Error(1)
}
