# ===== Settings =====
BUILD_DIR := tmp
BINARY    ?= main
GOOS      ?= linux
GOARCH    ?= arm64
CGO_ENABLED ?= 0

CMD_DIRS  := $(wildcard backend/cmd/*)
NAMES     := $(notdir $(CMD_DIRS))
ZIPS      := $(addprefix $(BUILD_DIR)/,$(addsuffix .zip,$(NAMES)))

.PHONY: package clean
package:
	@echo "==> Packaging all commands"
	@for dir in $(wildcard backend/cmd/*); do \
		name=$$(basename $$dir); \
		echo "Building $$name ..."; \
		mkdir -p $(BUILD_DIR)/$$name; \
		GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
			go build -o $(BUILD_DIR)/$$name/$(BINARY) $$dir; \
		(cd $(BUILD_DIR)/$$name && zip -j ../$$name.zip $(BINARY) >/dev/null); \
		rm -rf $(BUILD_DIR)/$$name; \
		echo "==> Created $(BUILD_DIR)/$$name.zip"; \
	done

clean:
	rm -rf $(BUILD_DIR)
