package panicking

import (
	"log"

	"github.com/stretchr/testify/mock"
)

// Panicker panics
type Panicker interface {
	Panic(msg string, fmtVars ...interface{})
}

// StandardPanicker is for use in production situations where a panic is appropriate.
// Such situations include:
// 		1. Inability to generate secrets due to some issue accessing /dev/urandom.
//		2. Inability to connect to a database after sufficient retry attempts.
// Something similarly disastrous should occur before StandardPanicker is invoked.
type StandardPanicker struct{}

// Panic implements our Panicker interface
func (s *StandardPanicker) Panic(msg string, fmtVars ...interface{}) {
	log.Panicf(msg, fmtVars...)
}

// MockPanicker is for testing that panics occur in the appropriate situations.
type MockPanicker struct {
	mock.Mock
}

// Panic implements our Panicker interface
func (m *MockPanicker) Panic(msg string, fmtVars ...interface{}) {
	m.Called(msg, fmtVars)
}
