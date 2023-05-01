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
else
EXE_EXT :=
LIB_EXT := .so
endif

PROJECT := encval

PLUGINS := $(notdir $(patsubst %/main.go,%,$(wildcard ./*/main.go)))
PLUGIN_objects := $(addprefix $(PROJECT)_,$(addsuffix $(LIB_EXT),$(PLUGINS)))

.PHONY: all
all: plugins

.PHONY: plugins
plugins: $(PLUGIN_objects)

$(PROJECT)_%$(LIB_EXT):
	cd $* ; GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) $(GO) build -buildmode=plugin \
		$(GOBUILD_FLAGS) \
		-ldflags="$(GOLD_FLAGS)" \
		-o ../$(PROJECT)_$*$(LIB_EXT) \
		.

.PHONY: test
test: $(addprefix test-,$(PLUGINS))

.PHONY: test-%
test-%:
	cd $*; $(GO) test $(GOTEST_FLAGS) ./...

.PHONY: format
format: $(addprefix format-,$(PLUGINS))

.PHONY: format-%
format-%:
	cd $*; $(GO) fmt $(GOFMT_FLAGS) ./...

.PHONY: lint
lint: $(addprefix lint-,$(PLUGINS))

.PHONY: lint-%
lint-%:
	cd $*; $(GO) vet $(GOVET_FLAGS) ./...

.PHONY: update-deps
update-deps: $(addprefix update-deps-,$(PLUGINS))

.PHONY: update-deps-%
update-deps-%:
	cd $*; $(GO) mod tidy

.PHONY: install
install: $(addprefix install-,$(PLUGINS))

.PHONY: install-%
install-%: $(PROJECT)_%$(LIB_EXT)
	$(INSTALL_DATA) -D -t $(DESTDIR)$(PREFIX)/lib/$(PROJECT) $^

.PHONY: clean
clean:
	$(RM) $(PLUGIN_objects)
