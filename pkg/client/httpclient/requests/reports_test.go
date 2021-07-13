package requests

import (
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_BuildReportExistsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/reports/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleReport := fakes.BuildFakeReport()

		actual, err := helper.builder.BuildReportExistsRequest(helper.ctx, exampleReport.ID)
		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, exampleReport.ID)

		assert.NoError(t, err)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid report ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildReportExistsRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleReport := fakes.BuildFakeReport()

		actual, err := helper.builder.BuildReportExistsRequest(helper.ctx, exampleReport.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetReportRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/reports/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleReport := fakes.BuildFakeReport()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleReport.ID)

		actual, err := helper.builder.BuildGetReportRequest(helper.ctx, exampleReport.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid report ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetReportRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleReport := fakes.BuildFakeReport()

		actual, err := helper.builder.BuildGetReportRequest(helper.ctx, exampleReport.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetReportsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/reports"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetReportsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetReportsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateReportRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/reports"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeReportCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateReportRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateReportRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateReportRequest(helper.ctx, &types.ReportCreationInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeReportCreationInput()

		actual, err := helper.builder.BuildCreateReportRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateReportRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/reports/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleReport := fakes.BuildFakeReport()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleReport.ID)

		actual, err := helper.builder.BuildUpdateReportRequest(helper.ctx, exampleReport)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateReportRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleReport := fakes.BuildFakeReport()

		actual, err := helper.builder.BuildUpdateReportRequest(helper.ctx, exampleReport)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveReportRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/reports/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleReport := fakes.BuildFakeReport()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleReport.ID)

		actual, err := helper.builder.BuildArchiveReportRequest(helper.ctx, exampleReport.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid report ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveReportRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleReport := fakes.BuildFakeReport()

		actual, err := helper.builder.BuildArchiveReportRequest(helper.ctx, exampleReport.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetAuditLogForReportRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/reports/%d/audit"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleReport := fakes.BuildFakeReport()

		actual, err := helper.builder.BuildGetAuditLogForReportRequest(helper.ctx, exampleReport.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, exampleReport.ID)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid report ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetAuditLogForReportRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleReport := fakes.BuildFakeReport()

		actual, err := helper.builder.BuildGetAuditLogForReportRequest(helper.ctx, exampleReport.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
