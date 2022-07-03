.PHONY: all
all: slurpd slurp

OUT=bin

.PHONY: slurpd
slurpd: $(OUT)/slurpd

.PHONY: slurp
slurp: $(OUT)/slurp

.PHONY: clean
clean:
	@rm $(OUT)/*
	@rm internal/api/slurp/*.pb.go

internal/api/slurp/*.pb.go: ./internal/api/slurp/slurp.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $<

$(OUT)/slurpd: internal/api/slurp/*.pb.go cmd/slurpd/*.go
	mkdir -p $(OUT)
	go build -o $@ ./cmd/slurpd

$(OUT)/slurp: internal/api/slurp/*.pb.go cmd/slurp/*.go
	mkdir -p $(OUT)
	go build -o $@ ./cmd/slurp
