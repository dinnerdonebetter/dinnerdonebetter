package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	grpcfiltering "github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/protobuf/types/known/timestamppb"
	g "maragu.dev/gomponents"
)

func fetchErrorString(err error, key string) string {
	var validErr validation.Errors
	if errors.As(err, &validErr) {
		if validationErr := validErr[key]; validationErr != nil {
			var validationLibError validation.ErrorObject
			if errors.As(validationErr, &validationLibError) {
				return validationLibError.Error()
			}
		}
	}

	return ""
}

func renderTimestamp(value any) g.Node {
	if value == nil {
		return g.Text("-")
	}

	switch v := value.(type) {
	case *timestamppb.Timestamp:
		if v == nil {
			return g.Text("-")
		}
		return g.Text(v.AsTime().Format("2006-01-02 15:04:05"))
	case timestamppb.Timestamp:
		return g.Text(v.AsTime().Format("2006-01-02 15:04:05"))
	default:
		return g.Text(fmt.Sprintf("%v", v))
	}
}

// buildCookie provides a consistent way of constructing an HTTP cookie.
func (s *AdminFrontendServer) buildCookie(ctx context.Context, value string) *http.Cookie {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	expiry := time.Now().Add(s.config.Cookies.Lifetime)

	// https://www.calhoun.io/securing-cookies-in-go/
	cookie := &http.Cookie{
		Name:     s.config.Cookies.CookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.config.Cookies.SecureOnly,
		// Domain:   s.config.Cookies.Domain,
		Expires:  expiry,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(expiry).Seconds()),
	}

	return cookie
}

// buildQueryFilterFromRequest extracts a QueryFilter from the HTTP request and converts it to a gRPC filter.
// This is a reusable helper for handlers that need to process pagination and filtering.
func buildQueryFilterFromRequest(req *http.Request) (*filtering.QueryFilter, *grpcfiltering.QueryFilter) {
	queryFilter := filtering.ExtractQueryFilterFromRequest(req)
	grpcFilter := grpcconverters.ConvertQueryFilterToGRPCQueryFilter(queryFilter, filtering.Pagination{})
	return queryFilter, grpcFilter
}

// buildPaginationFromGRPCResponse constructs a Pagination object from a gRPC Pagination response.
// It extracts cursor, counts, max response size, and previous cursor from the service response.
func buildPaginationFromGRPCResponse(grpcPagination *grpcfiltering.Pagination) *filtering.Pagination {
	if grpcPagination == nil {
		return nil
	}

	// Convert AppliedQueryFilter from gRPC to internal format
	var appliedQueryFilter *filtering.QueryFilter
	if grpcPagination.AppliedQueryFilter != nil {
		appliedQueryFilter = grpcconverters.ConvertGRPCQueryFilterToQueryFilter(grpcPagination.AppliedQueryFilter)
	}

	pagination := &filtering.Pagination{
		Cursor:             grpcPagination.Cursor,
		MaxResponseSize:    uint8(grpcPagination.MaxResponseSize),
		FilteredCount:      grpcPagination.FilteredCount,
		TotalCount:         grpcPagination.TotalCount,
		AppliedQueryFilter: appliedQueryFilter,
	}

	// Diagnostic logging
	var appliedCursorStr string
	if appliedQueryFilter != nil && appliedQueryFilter.Cursor != nil {
		appliedCursorStr = *appliedQueryFilter.Cursor
	}
	fmt.Printf("[PAGINATION BUILD DEBUG] grpcPreviousCursor='%s', appliedCursor='%s', grpcCursor='%s'\n",
		grpcPagination.PreviousCursor, appliedCursorStr, grpcPagination.Cursor)

	// Determine PreviousCursor for navigation:
	// Cursor-based pagination limitation: We can only go back one page without cursor history.
	//
	// The backend's PreviousCursor is set to the input filter's cursor (the cursor used to get to this page),
	// which is the same as AppliedCursor. This means:
	// - On page 2: AppliedCursor = cursor from page 1, PreviousCursor = cursor from page 1
	// - On page 3: AppliedCursor = cursor from page 2, PreviousCursor = cursor from page 2
	//
	// To go back:
	// - From page 2 to page 1: need empty cursor (we can detect this: if AppliedCursor exists, we're on page 2+)
	// - From page 3 to page 2: need cursor from page 1 (but we don't have it - this is the limitation)
	//
	// Current workaround: Only support going back from page 2 to page 1.
	// For page 3+, backward navigation is disabled (PreviousCursor is empty).
	//
	// TODO: Backend should provide cursor history or a way to get the previous page's cursor.
	if appliedQueryFilter != nil && appliedQueryFilter.Cursor != nil {
		appliedCursor := *appliedQueryFilter.Cursor
		if appliedCursor != "" {
			// We're on page 2+ (navigated here with a cursor)
			// Only support going back from page 2 to page 1.
			// We detect page 2 by checking if PreviousCursor equals AppliedCursor (both are cursor from page 1).
			// For page 3+, PreviousCursor equals AppliedCursor (both are cursor from page 2),
			// but we can't go back because we don't have the cursor from page 1.
			if grpcPagination.PreviousCursor != "" && grpcPagination.PreviousCursor == appliedCursor {
				// Likely page 2: can go back to page 1 with empty cursor
				pagination.PreviousCursor = ""
			} else {
				// Page 3+: Can't determine previous page cursor, disable backward navigation
				// TODO: Backend should provide cursor history
				pagination.PreviousCursor = ""
			}
		} else {
			// No cursor was used, so we're on the first page
			pagination.PreviousCursor = ""
		}
	} else {
		// No AppliedQueryFilter means we're on the first page
		pagination.PreviousCursor = ""
	}

	return pagination
}

// buildPaginationURLGenerator creates a reusable PaginationURLGenerator function that preserves
// query parameters (like search) and adds the cursor for pagination.
// This uses the main page URL for deep linking.
func buildPaginationURLGenerator(req *http.Request, baseURL string, queryFilter *filtering.QueryFilter) components.PaginationURLGenerator {
	return func(cursor string) string {
		queryParams := url.Values{}

		// Preserve existing query parameters from the request
		for key, values := range req.URL.Query() {
			// Skip cursor as we'll set it fresh (or omit it if empty)
			if key == "cursor" {
				continue
			}
			// Preserve all other params (search, limit, etc.)
			for _, value := range values {
				queryParams.Add(key, value)
			}
		}

		// Only add cursor if it's non-empty (empty cursor means go to first page)
		if cursor != "" {
			queryParams.Set("cursor", cursor)
		}

		// Ensure limit is included if it was in the original filter
		if queryFilter != nil && queryFilter.Limit != nil {
			queryParams.Set("limit", fmt.Sprintf("%d", *queryFilter.Limit))
		}

		return baseURL + "?" + queryParams.Encode()
	}
}

// buildPaginationURLGeneratorForSearch creates a PaginationURLGenerator that uses the search endpoint.
// This is used for HTMX pagination buttons to return just the table content.
func buildPaginationURLGeneratorForSearch(req *http.Request, searchEndpoint string, queryFilter *filtering.QueryFilter) components.PaginationURLGenerator {
	return func(cursor string) string {
		queryParams := url.Values{}

		// Preserve existing query parameters from the request
		for key, values := range req.URL.Query() {
			// Skip cursor as we'll set it fresh (or omit it if empty)
			if key == "cursor" {
				continue
			}
			// Preserve all other params (search, limit, etc.)
			for _, value := range values {
				queryParams.Add(key, value)
			}
		}

		// Only add cursor if it's non-empty (empty cursor means go to first page)
		if cursor != "" {
			queryParams.Set("cursor", cursor)
		}

		// Ensure limit is included if it was in the original filter
		if queryFilter != nil && queryFilter.Limit != nil {
			queryParams.Set("limit", fmt.Sprintf("%d", *queryFilter.Limit))
		}

		return searchEndpoint + "?" + queryParams.Encode()
	}
}
