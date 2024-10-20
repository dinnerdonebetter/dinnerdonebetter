package types

import (
	"net/http"
)

type (
	// DataPrivacyService describes a structure capable of serving CCPA/GRPC-related requests.
	DataPrivacyService interface {
		DataDeletionHandler(http.ResponseWriter, *http.Request)
	}
)
