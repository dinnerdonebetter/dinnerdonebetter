setup:
	(cd backend && make setup)
	(cd frontend && make setup)

format:
	(cd backend && make format)
	(cd frontend && make format)

BACKEND_PROTO_DESTINATION = backend/internal/proto

.PHONY: proto
proto:
	rm -rf $(BACKEND_PROTO_DESTINATION)
	mkdir -p $(BACKEND_PROTO_DESTINATION)/types
	protoc -I=proto \
		--go_opt=paths=source_relative \
		--go_out=$(BACKEND_PROTO_DESTINATION)/types \
		proto/types.proto
	mkdir -p $(BACKEND_PROTO_DESTINATION)/service
	protoc -I=proto \
		--go_opt=paths=source_relative \
		--go_out=$(BACKEND_PROTO_DESTINATION)/service \
		proto/service.proto