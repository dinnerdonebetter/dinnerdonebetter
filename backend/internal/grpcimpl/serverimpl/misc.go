package serverimpl

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return &emptypb.Empty{}, nil
}
