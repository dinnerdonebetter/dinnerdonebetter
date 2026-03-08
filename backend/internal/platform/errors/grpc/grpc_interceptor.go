package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryErrorEncodingInterceptor returns a unary interceptor that encodes handler
// errors into gRPC status details for wire transmission.
// Handlers should return errors (optionally wrapped); the interceptor will
// derive the gRPC code via MapToGRPC and attach the encoded error to details.
func UnaryErrorEncodingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		// If already a status.Error, pass through (handler already set the code).
		if st, ok := status.FromError(err); ok {
			code := MapToGRPC(err, st.Code())
			if code != st.Code() {
				return nil, status.Error(code, st.Message())
			}
			return nil, err
		}

		code := MapToGRPC(err, codes.Unknown)
		return nil, status.Error(code, err.Error())
	}
}

// StreamErrorEncodingInterceptor returns a stream interceptor that encodes
// handler errors into gRPC status details for wire transmission.
func StreamErrorEncodingInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, ss)
		if err == nil {
			return nil
		}

		if st, ok := status.FromError(err); ok {
			code := MapToGRPC(err, st.Code())
			if code != st.Code() {
				return status.Error(code, st.Message())
			}
			return err
		}

		code := MapToGRPC(err, codes.Unknown)
		return status.Error(code, err.Error())
	}
}
