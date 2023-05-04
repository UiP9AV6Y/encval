INSTALL ?= install
INSTALL_DIR = $(INSTALL) -d
INSTALL_DATA = $(INSTALL) -m0644
INSTALL_PROGRAM = $(INSTALL) -m0755

TAR ?= tar
ZIP ?= zip
UNZIP ?= unzip

SUDO ?= sudo

PREFIX ?= /usr/local
GITHUB_ENV ?= /dev/null
GITHUB_PATH ?= /dev/null

GO ?= go
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)
GOBUILD_TAGS ?=
GOBUILD_FLAGS ?=
GOLD_FLAGS ?= -s -w
GOTEST_FLAGS ?=
GOFMT_FLAGS ?=
GOVET_FLAGS ?=
CGO_ENABLED ?= 0
GO_DO := GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) $(GO)

ifeq ($(GOOS),windows)
EXE_EXT := .exe
LIB_EXT := .dll
DIST_EXT := .zip
else
EXE_EXT :=
LIB_EXT := .so
DIST_EXT := .tgz
endif

PROJECT := encval

ENTRYPOINTS := $(notdir $(patsubst %/main.go,%,$(wildcard ./cmd/*/main.go)))
PROGRAM_files := $(addsuffix $(EXE_EXT),$(ENTRYPOINTS))

TARGET_TRIPLE := $(PROJECT)_$(GOOS)_$(GOARCH)$(GOARM)

.PHONY: all
all: binaries plugins

.PHONY: binaries
binaries: $(PROGRAM_files)

.PHONY: plugins
plugins:
	$(MAKE) -C plugins

$(PROGRAM_files):
	 $(GO_DO) build \
		$(GOBUILD_FLAGS) \
		-tags "$(GOBUILD_TAGS)" \
		-ldflags "$(GOLD_FLAGS)" \
		-o $@ \
		./cmd/$(basename $@)

.PHONY: dist
dist: $(TARGET_TRIPLE)$(DIST_EXT)

.PHONY: $(TARGET_TRIPLE)
$(TARGET_TRIPLE): LICENSE.txt
	$(INSTALL_DATA) -D -t $@ $^

.PHONY: $(TARGET_TRIPLE)/bin
$(TARGET_TRIPLE)/bin: $(PROGRAM_files)
	$(INSTALL_PROGRAM) -D -t $@ $^

.PHONY: $(TARGET_TRIPLE)/lib
$(TARGET_TRIPLE)/lib: plugins
	$(INSTALL_DATA) -D -t $@ $(wildcard plugins/*$(LIB_EXT))

$(TARGET_TRIPLE).tgz: $(TARGET_TRIPLE) $(TARGET_TRIPLE)/bin $(TARGET_TRIPLE)/lib
	cd $< ; $(TAR) -czf ../$@ *

$(TARGET_TRIPLE).zip: $(TARGET_TRIPLE) $(TARGET_TRIPLE)/bin
	cd $< ; $(ZIP) -r ../$@ *

.PHONY: ci-deps
ci-deps:
	$(SUDO) apt-get update -qq
	$(SUDO) apt-get install -y --no-install-recommends \
		build-essential

.PHONY: ci-env
ci-env:
	echo "LD_LIBRARY_PATH=$(DESTDIR)$(PREFIX)/lib" >> $(GITHUB_ENV)
	echo "$(DESTDIR)$(PREFIX)/bin" >> $(GITHUB_PATH)

.PHONY: test
test:
	$(GO_DO) test $(GOTEST_FLAGS) ./...

.PHONY: format
format:
	$(GO_DO) fmt $(GOFMT_FLAGS) ./...

.PHONY: lint
lint:
	$(GO_DO) vet $(GOVET_FLAGS) ./...

.PHONY: update-deps
update-deps:
	$(GO_DO) mod tidy

.PHONY: install
install: $(PROGRAM_files)
	$(INSTALL_PROGRAM) -D -t $(DESTDIR)$(PREFIX)/bin $^

.PHONY: install-dist
install-dist: install$(DIST_EXT)

$(DESTDIR)$(PREFIX):
	$(INSTALL_DIR) $(DESTDIR)$(PREFIX)

.PHONY: install.tgz
install.tgz: $(TARGET_TRIPLE).tgz $(DESTDIR)$(PREFIX)
	$(TAR) -xzf $< -C $(DESTDIR)$(PREFIX)

.PHONY: install.zip
install.zip: $(TARGET_TRIPLE).zip $(DESTDIR)$(PREFIX)
	$(UNZIP) $< -d $(DESTDIR)$(PREFIX)

.PHONY: clean
clean:
	$(RM) *.zip *.tgz
	$(RM) $(PROGRAM_files)
	$(RM) -r $(PROJECT)_*

.PHONY: %-all
%-all: %
	$(MAKE) -C plugins $*

