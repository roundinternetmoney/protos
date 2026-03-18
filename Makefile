.PHONY: examples buf all version tag release

LATEST_TAG := $(shell git describe --tags --abbrev=0 --match 'v[0-9]*' 2>/dev/null || echo v0.0.0)
VERSION := $(patsubst v%,%,$(LATEST_TAG))
PATCH_VERSION := $(shell echo $(VERSION) | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3+1}')
NEW_VERSION ?= $(PATCH_VERSION)

buf:
	buf dep update
	buf lint
	buf build

all: buf
	buf push

version:
	@echo "Current version: $(VERSION)"
	@echo "Release version: $(NEW_VERSION)"

tag:
	git tag -a v$(NEW_VERSION) -m "Release v$(NEW_VERSION)"
	git push origin v$(NEW_VERSION)

release: buf tag
	gh release create v$(NEW_VERSION) \
		--title "v$(NEW_VERSION)" \
		--notes "Release v$(NEW_VERSION)"