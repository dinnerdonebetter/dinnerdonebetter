// Package domainreg defines the extension-point types that domain registration modules use
// to inject their gRPC services and method permissions into the platform gRPC API builder
// without creating an import cycle between the builder and the domain module.
//
// Usage pattern:
//   - The grpcapi builder (extras.go) reads these types from the DI container.
//   - Each domain registration module provides these types into the same DI container.
//   - The grpcapi builder (build.go) remains the single point that calls each domain's
//     RegisterForGRPCAPI function — one line per domain, nothing else needs to change.
package domainreg

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/auth/grpc/interceptors"

	platformgrpc "github.com/primandproper/platform/server/grpc"

	grpc "google.golang.org/grpc"
)

// ExtraRegistrationFuncs is the named DI type that domain registration modules provide
// to register their gRPC service implementations on the underlying *grpc.Server.
// Each entry should call the generated RegisterXxxServiceServer helper.
type ExtraRegistrationFuncs []func(s *grpc.Server)

// ToRegistrationFuncs converts ExtraRegistrationFuncs to the platform RegistrationFunc slice.
func (e ExtraRegistrationFuncs) ToRegistrationFuncs() []platformgrpc.RegistrationFunc {
	out := make([]platformgrpc.RegistrationFunc, len(e))
	for idx, fn := range e {
		out[idx] = fn
	}
	return out
}

// ExtraMethodPermissions is the named DI type that domain registration modules provide
// to contribute their method→permission mappings to the auth interceptor.
type ExtraMethodPermissions interceptors.MethodPermissionsMap
