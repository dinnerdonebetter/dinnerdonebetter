package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	auditsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"
	grpctypes "github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/audit/grpc/converters"

	"google.golang.org/grpc/codes"
)

const (
	o11yName = "audit_service"
)

var _ auditsvc.AuditServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		auditsvc.UnimplementedAuditServiceServer
		tracer          tracing.Tracer
		logger          logging.Logger
		auditRepository audit.Repository
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditRepository audit.Repository,
) auditsvc.AuditServiceServer {
	return &ServiceImpl{
		logger:          logging.EnsureLogger(logger).WithName(o11yName),
		tracer:          tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		auditRepository: auditRepository,
	}
}

func (s *ServiceImpl) GetAuditLogEntriesForAccount(ctx context.Context, request *auditsvc.GetAuditLogEntriesForAccountRequest) (*auditsvc.GetAuditLogEntriesForAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("", "")
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	auditLogEntries, err := s.auditRepository.GetAuditLogEntriesForAccount(ctx, "TODO", filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "")
	}

	x := &auditsvc.GetAuditLogEntriesForAccountResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			CurrentAccountID: "",
			TraceID:          span.SpanContext().TraceID().String(),
		},
		Filter:  request.Filter,
		Results: nil,
	}

	for _, y := range auditLogEntries.Data {
		x.Results = append(x.Results, converters.ConvertAuditLogEntryToGRPCAuditLogEntry(y))
	}

	return x, nil
}

func (s *ServiceImpl) GetAuditLogEntriesForUser(ctx context.Context, request *auditsvc.GetAuditLogEntriesForUserRequest) (*auditsvc.GetAuditLogEntriesForUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("", "")
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	auditLogEntries, err := s.auditRepository.GetAuditLogEntriesForUser(ctx, "TODO", filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "")
	}

	x := &auditsvc.GetAuditLogEntriesForUserResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			CurrentAccountID: "",
			TraceID:          span.SpanContext().TraceID().String(),
		},
		Filter:  request.Filter,
		Results: nil,
	}

	for _, y := range auditLogEntries.Data {
		x.Results = append(x.Results, converters.ConvertAuditLogEntryToGRPCAuditLogEntry(y))
	}

	return x, nil
}

func (s *ServiceImpl) GetAuditLogEntryByID(ctx context.Context, request *auditsvc.GetAuditLogEntryByIDRequest) (*auditsvc.GetAuditLogEntryByIDResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.AuditLogEntryIDKey, request.AuditLogEntryID)
	auditLogEntry, err := s.auditRepository.GetAuditLogEntry(ctx, request.AuditLogEntryID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "")
	}

	returnValue := converters.ConvertAuditLogEntryToGRPCAuditLogEntry(auditLogEntry)

	x := &auditsvc.GetAuditLogEntryByIDResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			CurrentAccountID: "",
			TraceID:          span.SpanContext().TraceID().String(),
		},
		Result: returnValue,
	}

	return x, nil
}
