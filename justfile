alias tv := test-verbose

test:
  go test ./...

test-verbose:
  go test ./... -test.v

proto:
	protoc \
		--proto_path=. \
		--go_out=:. \
		proto/tiktok.proto
