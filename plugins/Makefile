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
EXT := .dll
else
EXT := .so
endif

PROJECT := encval

PLUGINS := $(notdir $(patsubst %/main.go,%,$(wildcard ./*/main.go)))
PLUGIN_objects := $(addprefix $(PROJECT)_,$(addsuffix $(EXT),$(PLUGINS)))

.PHONY: all
all: plugins

.PHONY: plugins
plugins: $(PLUGIN_objects)

$(PROJECT)_%$(EXT):
	cd $* ; GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -buildmode=plugin \
		$(GOBUILD_FLAGS) \
		-ldflags="$(GOLD_FLAGS)" \
		-o ../$(PROJECT)_$*$(EXT) \
		.

.PHONY: update-deps
update-deps:
	for P in $(PLUGINS); do \
		$(MAKE) mod-tidy PLUGIN=$$P; \
	done

.PHONY: mod-tidy
mod-tidy: PLUGIN=$(firstword $(PLUGINS))
mod-tidy:
	cd $(PLUGIN); $(GO) mod tidy

.PHONY: install
install:
	$(INSTALL_DATA) -D -t $(DESTDIR)$(PREFIX)/lib/$(PROJECT) $(PLUGIN_objects)

.PHONY: clean
clean:
	$(RM) $(PLUGIN_objects)
