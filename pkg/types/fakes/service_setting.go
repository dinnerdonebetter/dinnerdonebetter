package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeServiceSetting builds a faked valid preparation.
func BuildFakeServiceSetting() *types.ServiceSetting {
	return &types.ServiceSetting{
		ID:           BuildFakeID(),
		Name:         buildUniqueString(),
		Type:         "user",
		Description:  buildUniqueString(),
		Enumeration:  []string{},
		DefaultValue: pointers.Pointer(buildUniqueString()),
		AdminsOnly:   true,
		CreatedAt:    BuildFakeTime(),
	}
}

// BuildFakeServiceSettingList builds a faked ServiceSettingList.
func BuildFakeServiceSettingList() *types.QueryFilteredResult[types.ServiceSetting] {
	var examples []*types.ServiceSetting
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeServiceSetting())
	}

	return &types.QueryFilteredResult[types.ServiceSetting]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeServiceSettingUpdateRequestInput builds a faked ServiceSettingUpdateRequestInput from a valid preparation.
func BuildFakeServiceSettingUpdateRequestInput() *types.ServiceSettingUpdateRequestInput {
	validPreparation := BuildFakeServiceSetting()
	return converters.ConvertServiceSettingToServiceSettingUpdateRequestInput(validPreparation)
}

// BuildFakeServiceSettingCreationRequestInput builds a faked ServiceSettingCreationRequestInput.
func BuildFakeServiceSettingCreationRequestInput() *types.ServiceSettingCreationRequestInput {
	validPreparation := BuildFakeServiceSetting()
	return converters.ConvertServiceSettingToServiceSettingCreationRequestInput(validPreparation)
}
