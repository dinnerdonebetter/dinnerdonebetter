package dataprivacy

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/dataprivacy"
)

func (r *repository) FetchUserDataCollection(ctx context.Context, userID string) (*dataprivacy.UserDataCollectionResponse, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	// TODO
	x := &dataprivacy.UserDataCollectionResponse{}

	return x, nil
}

func (r *repository) DeleteUser(ctx context.Context, userID string) error {
	// TODO
	return nil
}
