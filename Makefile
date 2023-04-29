INSTALL ?= install
INSTALL_DIR = $(INSTALL) -d
INSTALL_DATA = $(INSTALL) -m0644
INSTALL_PROGRAM = $(INSTALL) -m0755

GO ?= go
GOOS ?= $(shell $(GO) env GOOS)
GOLD_FLAGS ?= -s -w

ifeq ($(GOOS),windows)
EXT := .exe
else
EXT :=
endif

ENTRYPOINTS := encval
PROGRAM_files = $(addsuffix $(EXT),$(ENTRYPOINTS))

.PHONY: all
all: binaries plugins

.PHONY: binaries
binaries: $(PROGRAM_files)

$(PROGRAM_files):
	$(GO) build $(GOBUILD_FLAGS) \
		-ldflags="$(GOLD_FLAGS)" \
		-o $@ \
		./cmd/$(basename $@)

.PHONY: test
test:
	$(GO) test ./...

.PHONY: format
format:
	$(GO) fmt ./...

.PHONY: lint
lint:
	$(GO) vet ./...

.PHONY: plugins
plugins:
	$(MAKE) -C plugins

.PHONY: update-deps
update-deps:
	$(GO) mod tidy
	$(MAKE) -C plugins $@

.PHONY: clean
clean:
	$(RM) $(PROGRAM_files)
	$(MAKE) -C plugins $@
