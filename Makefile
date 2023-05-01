INSTALL ?= install
INSTALL_DIR = $(INSTALL) -d
INSTALL_DATA = $(INSTALL) -m0644
INSTALL_PROGRAM = $(INSTALL) -m0755

PREFIX ?= /usr/local

GO ?= go
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)
GOLD_FLAGS ?= -s -w
GOTEST_FLAGS ?= -race
GOFMT_FLAGS ?=
GOVET_FLAGS ?=
CGO_ENABLED ?= 0

ifeq ($(GOOS),windows)
EXT := .exe
else
EXT :=
endif

ENTRYPOINTS := $(notdir $(patsubst %/main.go,%,$(wildcard ./cmd/*/main.go)))
PROGRAM_files := $(addsuffix $(EXT),$(ENTRYPOINTS))

.PHONY: all
all: binaries plugins

.PHONY: binaries
binaries: $(PROGRAM_files)

.PHONY: plugins
plugins:
	$(MAKE) -C plugins

$(PROGRAM_files):
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOBUILD_FLAGS) \
		-ldflags="$(GOLD_FLAGS)" \
		-o $@ \
		./cmd/$(basename $@)

.PHONY: test
test:
	$(GO) test $(GOTEST_FLAGS) ./...

.PHONY: format
format:
	$(GO) fmt $(GOFMT_FLAGS) ./...

.PHONY: lint
lint:
	$(GO) vet $(GOVET_FLAGS) ./...

.PHONY: update-deps
update-deps:
	$(GO) mod tidy

.PHONY: install
install: $(PROGRAM_files)
	$(INSTALL_PROGRAM) -D -t $(DESTDIR)$(PREFIX)/bin $^

.PHONY: clean
clean:
	$(RM) $(PROGRAM_files)

.PHONY: %-all
%-all: %
	$(MAKE) -C plugins $*

