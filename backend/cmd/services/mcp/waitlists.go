package main

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/waitlists"
	grpcconverters "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/converters"
	waitlistsgrpc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/waitlists"
	waitlistsconverters "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/waitlists/grpc/converters"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/verygoodsoftwarenotvirus/platform/database/filtering"
)

var waitlistSchema = map[string]any{
	"ID":            stringField("The ID of the waitlist"),
	"Name":          stringField("The waitlist name"),
	"Description":   stringField("The waitlist description"),
	"ValidUntil":    timestampField("When the waitlist expires"),
	"CreatedAt":     timestampField("When the waitlist was created"),
	"LastUpdatedAt": timestampField("When the waitlist was last updated"),
	"ArchivedAt":    timestampField("When the waitlist was archived"),
}

var waitlistSignupSchema = map[string]any{
	"ID":                stringField("The ID of the waitlist signup"),
	"Notes":             stringField("Notes about the signup"),
	"BelongsToWaitlist": stringField("The ID of the waitlist"),
	"BelongsToUser":     stringField("The ID of the user who signed up"),
	"BelongsToAccount":  stringField("The ID of the account"),
	"CreatedAt":         timestampField("When the signup was created"),
	"LastUpdatedAt":     timestampField("When the signup was last updated"),
	"ArchivedAt":        timestampField("When the signup was archived"),
}

var getWaitlistTool = &mcp.Tool{
	Name:        "GetWaitlist",
	Description: "Get a waitlist by its ID",
	InputSchema: schemaObject(map[string]any{
		"WaitlistID": stringField("The ID of the waitlist to get"),
	}),
	OutputSchema: schemaObject(waitlistSchema),
}

type GetWaitlistInvocation struct {
	WaitlistID string `jsonschema:"description=The waitlist ID"`
}

func (h *mcpToolManager) GetWaitlist() mcp.ToolHandlerFor[*GetWaitlistInvocation, *waitlists.Waitlist] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetWaitlistInvocation) (*mcp.CallToolResult, *waitlists.Waitlist, error) {
		result, err := h.client.GetWaitlist(ctx, &waitlistsgrpc.GetWaitlistRequest{
			WaitlistId: x.WaitlistID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, waitlistsconverters.ConvertGRPCWaitlistToWaitlist(result.Result), nil
	}
}

var getWaitlistsTool = &mcp.Tool{
	Name:        "GetWaitlists",
	Description: "Get waitlists with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(waitlistSchema)),
	}),
}

type (
	GetWaitlistsInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetWaitlistsResult struct {
		Results []*waitlists.Waitlist
	}
)

func (h *mcpToolManager) GetWaitlists() mcp.ToolHandlerFor[*GetWaitlistsInvocation, *GetWaitlistsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetWaitlistsInvocation) (*mcp.CallToolResult, *GetWaitlistsResult, error) {
		results, err := h.client.GetWaitlists(ctx, &waitlistsgrpc.GetWaitlistsRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetWaitlistsResult{}
		for _, w := range results.Results {
			out.Results = append(out.Results, waitlistsconverters.ConvertGRPCWaitlistToWaitlist(w))
		}
		return nil, out, nil
	}
}

var getActiveWaitlistsTool = &mcp.Tool{
	Name:        "GetActiveWaitlists",
	Description: "Get waitlists that are currently active (not expired)",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(waitlistSchema)),
	}),
}

type (
	GetActiveWaitlistsInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetActiveWaitlistsResult struct {
		Results []*waitlists.Waitlist
	}
)

func (h *mcpToolManager) GetActiveWaitlists() mcp.ToolHandlerFor[*GetActiveWaitlistsInvocation, *GetActiveWaitlistsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetActiveWaitlistsInvocation) (*mcp.CallToolResult, *GetActiveWaitlistsResult, error) {
		results, err := h.client.GetActiveWaitlists(ctx, &waitlistsgrpc.GetActiveWaitlistsRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetActiveWaitlistsResult{}
		for _, w := range results.Results {
			out.Results = append(out.Results, waitlistsconverters.ConvertGRPCWaitlistToWaitlist(w))
		}
		return nil, out, nil
	}
}

var createWaitlistTool = &mcp.Tool{
	Name:        "CreateWaitlist",
	Description: "Create a new waitlist",
	InputSchema: schemaObject(map[string]any{
		"Name":        stringField("The waitlist name"),
		"Description": stringField("The waitlist description"),
		"ValidUntil":  timestampField("When the waitlist expires (ISO 8601)"),
	}),
	OutputSchema: schemaObject(waitlistSchema),
}

type CreateWaitlistInvocation struct {
	ValidUntil  time.Time `jsonschema:"description=When the waitlist expires"`
	Name        string    `jsonschema:"description=The waitlist name"`
	Description string    `jsonschema:"description=The waitlist description"`
}

func (h *mcpToolManager) CreateWaitlist() mcp.ToolHandlerFor[*CreateWaitlistInvocation, *waitlists.Waitlist] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateWaitlistInvocation) (*mcp.CallToolResult, *waitlists.Waitlist, error) {
		result, err := h.client.CreateWaitlist(ctx, &waitlistsgrpc.CreateWaitlistRequest{
			Input: &waitlistsgrpc.WaitlistCreationRequestInput{
				Name:        x.Name,
				Description: x.Description,
				ValidUntil:  grpcconverters.ConvertTimeToPBTimestamp(x.ValidUntil),
			},
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, waitlistsconverters.ConvertGRPCWaitlistToWaitlist(result.Created), nil
	}
}

var archiveWaitlistTool = &mcp.Tool{
	Name:        "ArchiveWaitlist",
	Description: "Archive (soft-delete) a waitlist",
	InputSchema: schemaObject(map[string]any{
		"WaitlistID": stringField("The ID of the waitlist to archive"),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Success": boolField("Whether the archive was successful"),
	}),
}

type ArchiveWaitlistInvocation struct {
	WaitlistID string `jsonschema:"required,description=The ID of the waitlist to archive"`
}

func (h *mcpToolManager) ArchiveWaitlist() mcp.ToolHandlerFor[*ArchiveWaitlistInvocation, *boolResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *ArchiveWaitlistInvocation) (*mcp.CallToolResult, *boolResult, error) {
		_, err := h.client.ArchiveWaitlist(ctx, &waitlistsgrpc.ArchiveWaitlistRequest{
			WaitlistId: x.WaitlistID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, &boolResult{Success: true}, nil
	}
}

var getWaitlistSignupTool = &mcp.Tool{
	Name:        "GetWaitlistSignup",
	Description: "Get a specific waitlist signup",
	InputSchema: schemaObject(map[string]any{
		"WaitlistSignupID": stringField("The ID of the waitlist signup"),
		"WaitlistID":       stringField("The ID of the waitlist"),
	}),
	OutputSchema: schemaObject(waitlistSignupSchema),
}

type GetWaitlistSignupInvocation struct {
	WaitlistSignupID string `jsonschema:"description=The waitlist signup ID"`
	WaitlistID       string `jsonschema:"description=The waitlist ID"`
}

func (h *mcpToolManager) GetWaitlistSignup() mcp.ToolHandlerFor[*GetWaitlistSignupInvocation, *waitlists.WaitlistSignup] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetWaitlistSignupInvocation) (*mcp.CallToolResult, *waitlists.WaitlistSignup, error) {
		result, err := h.client.GetWaitlistSignup(ctx, &waitlistsgrpc.GetWaitlistSignupRequest{
			WaitlistSignupId: x.WaitlistSignupID,
			WaitlistId:       x.WaitlistID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, waitlistsconverters.ConvertGRPCWaitlistSignupToWaitlistSignup(result.Result), nil
	}
}

var getWaitlistSignupsForWaitlistTool = &mcp.Tool{
	Name:        "GetWaitlistSignupsForWaitlist",
	Description: "Get all signups for a specific waitlist",
	InputSchema: schemaObject(map[string]any{
		"WaitlistID": stringField("The ID of the waitlist"),
		"Filter":     queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(waitlistSignupSchema)),
	}),
}

type (
	GetWaitlistSignupsForWaitlistInvocation struct {
		Filter     *filtering.QueryFilter
		WaitlistID string `jsonschema:"description=The waitlist ID"`
	}

	GetWaitlistSignupsForWaitlistResult struct {
		Results []*waitlists.WaitlistSignup
	}
)

func (h *mcpToolManager) GetWaitlistSignupsForWaitlist() mcp.ToolHandlerFor[*GetWaitlistSignupsForWaitlistInvocation, *GetWaitlistSignupsForWaitlistResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetWaitlistSignupsForWaitlistInvocation) (*mcp.CallToolResult, *GetWaitlistSignupsForWaitlistResult, error) {
		results, err := h.client.GetWaitlistSignupsForWaitlist(ctx, &waitlistsgrpc.GetWaitlistSignupsForWaitlistRequest{
			WaitlistId: x.WaitlistID,
			Filter:     grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetWaitlistSignupsForWaitlistResult{}
		for _, s := range results.Results {
			out.Results = append(out.Results, waitlistsconverters.ConvertGRPCWaitlistSignupToWaitlistSignup(s))
		}
		return nil, out, nil
	}
}

var createWaitlistSignupTool = &mcp.Tool{
	Name:        "CreateWaitlistSignup",
	Description: "Create a new waitlist signup",
	InputSchema: schemaObject(map[string]any{
		"Notes":             stringField("Notes about the signup"),
		"BelongsToWaitlist": stringField("The ID of the waitlist to sign up for"),
		"BelongsToUser":     stringField("The ID of the user signing up"),
		"BelongsToAccount":  stringField("The ID of the account"),
	}),
	OutputSchema: schemaObject(waitlistSignupSchema),
}

type CreateWaitlistSignupInvocation struct {
	Notes             string `jsonschema:"description=Notes about the signup"`
	BelongsToWaitlist string `jsonschema:"description=The ID of the waitlist"`
	BelongsToUser     string `jsonschema:"description=The ID of the user"`
	BelongsToAccount  string `jsonschema:"description=The ID of the account"`
}

func (h *mcpToolManager) CreateWaitlistSignup() mcp.ToolHandlerFor[*CreateWaitlistSignupInvocation, *waitlists.WaitlistSignup] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateWaitlistSignupInvocation) (*mcp.CallToolResult, *waitlists.WaitlistSignup, error) {
		result, err := h.client.CreateWaitlistSignup(ctx, &waitlistsgrpc.CreateWaitlistSignupRequest{
			Input: &waitlistsgrpc.WaitlistSignupCreationRequestInput{
				Notes:             x.Notes,
				BelongsToWaitlist: x.BelongsToWaitlist,
				BelongsToUser:     x.BelongsToUser,
				BelongsToAccount:  x.BelongsToAccount,
			},
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, waitlistsconverters.ConvertGRPCWaitlistSignupToWaitlistSignup(result.Created), nil
	}
}

var archiveWaitlistSignupTool = &mcp.Tool{
	Name:        "ArchiveWaitlistSignup",
	Description: "Archive (soft-delete) a waitlist signup",
	InputSchema: schemaObject(map[string]any{
		"WaitlistSignupID": stringField("The ID of the waitlist signup to archive"),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Success": boolField("Whether the archive was successful"),
	}),
}

type ArchiveWaitlistSignupInvocation struct {
	WaitlistSignupID string `jsonschema:"required,description=The ID of the waitlist signup to archive"`
}

func (h *mcpToolManager) ArchiveWaitlistSignup() mcp.ToolHandlerFor[*ArchiveWaitlistSignupInvocation, *boolResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *ArchiveWaitlistSignupInvocation) (*mcp.CallToolResult, *boolResult, error) {
		_, err := h.client.ArchiveWaitlistSignup(ctx, &waitlistsgrpc.ArchiveWaitlistSignupRequest{
			WaitlistSignupId: x.WaitlistSignupID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, &boolResult{Success: true}, nil
	}
}
