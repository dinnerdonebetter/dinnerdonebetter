package models

import (
	"context"
	"net/http"
)

type (
	// Instrument represents an instrument
	Instrument struct {
		ID          uint64  `json:"id"`
		Name        string  `json:"name"`
		Variant     string  `json:"variant"`
		Description string  `json:"description"`
		Icon        string  `json:"icon"`
		CreatedOn   uint64  `json:"created_on"`
		UpdatedOn   *uint64 `json:"updated_on"`
		ArchivedOn  *uint64 `json:"archived_on"`
	}

	// InstrumentList represents a list of instruments
	InstrumentList struct {
		Pagination
		Instruments []Instrument `json:"instruments"`
	}

	// InstrumentCreationInput represents what a user could set as input for creating instruments
	InstrumentCreationInput struct {
		Name        string `json:"name"`
		Variant     string `json:"variant"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	// InstrumentUpdateInput represents what a user could set as input for updating instruments
	InstrumentUpdateInput struct {
		Name        string `json:"name"`
		Variant     string `json:"variant"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	// InstrumentDataManager describes a structure capable of storing instruments permanently
	InstrumentDataManager interface {
		GetInstrument(ctx context.Context, instrumentID, userID uint64) (*Instrument, error)
		GetInstrumentCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllInstrumentsCount(ctx context.Context) (uint64, error)
		GetInstruments(ctx context.Context, filter *QueryFilter, userID uint64) (*InstrumentList, error)
		GetAllInstrumentsForUser(ctx context.Context, userID uint64) ([]Instrument, error)
		CreateInstrument(ctx context.Context, input *InstrumentCreationInput) (*Instrument, error)
		UpdateInstrument(ctx context.Context, updated *Instrument) error
		ArchiveInstrument(ctx context.Context, id, userID uint64) error
	}

	// InstrumentDataServer describes a structure capable of serving traffic related to instruments
	InstrumentDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an InstrumentInput with an instrument
func (x *Instrument) Update(input *InstrumentUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Variant != "" && input.Variant != x.Variant {
		x.Variant = input.Variant
	}

	if input.Description != "" && input.Description != x.Description {
		x.Description = input.Description
	}

	if input.Icon != "" && input.Icon != x.Icon {
		x.Icon = input.Icon
	}
}

// ToInput creates a InstrumentUpdateInput struct for an instrument
func (x *Instrument) ToInput() *InstrumentUpdateInput {
	return &InstrumentUpdateInput{
		Name:        x.Name,
		Variant:     x.Variant,
		Description: x.Description,
		Icon:        x.Icon,
	}
}
