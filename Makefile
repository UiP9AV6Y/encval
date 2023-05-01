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
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOBUILD_FLAGS) \
		-ldflags="$(GOLD_FLAGS)" \
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
	cd $< ; tar -czf ../$@ *

$(TARGET_TRIPLE).zip: $(TARGET_TRIPLE) $(TARGET_TRIPLE)/bin
	cd $< ; zip -r ../$@ *

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
	$(RM) *.zip *.tgz
	$(RM) $(PROGRAM_files)
	$(RM) -r $(PROJECT)_*

.PHONY: %-all
%-all: %
	$(MAKE) -C plugins $*

