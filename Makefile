INSTALL ?= install
INSTALL_DIR = $(INSTALL) -d
INSTALL_DATA = $(INSTALL) -m0644
INSTALL_PROGRAM = $(INSTALL) -m0755

PREFIX ?= /usr/local

GO ?= go
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)
GOLD_FLAGS ?= -s -w

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
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(GOBUILD_FLAGS) \
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

