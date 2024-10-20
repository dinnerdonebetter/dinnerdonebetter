package fakes

import "github.com/dinnerdonebetter/backend/pkg/types"

func BuildFakeUserDataCollectionResponse() *types.UserDataCollectionResponse {
	return &types.UserDataCollectionResponse{
		ReportID: BuildFakeID(),
	}
}
