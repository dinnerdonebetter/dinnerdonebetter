# TODO

just a list of gripes and things I'd like to accomplish eventually

- [ ] concrete and distinct `authn` and `authz` packages. Everything's smushed together right now in the worst way.
- [ ] `platform/errors` is for better error building, `internal/errors` is for common service-level errors and is used ubiquitously in application code
- [ ] full use of `go tool` for tool invocations
- [ ] ensure all repositories create audit log entries
- [ ] ensure all gRPC methods are wrappers of a "manager" method
- [ ] integration tests verify ownership things (i.e. created recipes cannot be deleted by people who didn't create them)
- [ ] integration tests verify audit log entries are made
- [ ] integration tests verify that query filtering works
- [ ] all GRPC payload routes (creates/updates) reference an "input" field
- [ ] better error handling (i.e. naming missing parameters)