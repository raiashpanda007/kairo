PROTO_DIR   := proto
OUT_DIR     := internal/pb

PROTOC      := protoc
PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

.PHONY: proto clean build

proto: $(PROTO_FILES)
	@mkdir -p $(OUT_DIR)
	$(PROTOC) \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		--proto_path=$(PROTO_DIR) \
		$(PROTO_FILES)
	@echo "proto generation done -> $(OUT_DIR)"

build:
	go build -o bin/kairo ./cmd/...

clean:
	rm -rf $(OUT_DIR) bin/
