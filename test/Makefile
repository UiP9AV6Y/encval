INSTALL ?= install
INSTALL_DIR = $(INSTALL) -d
INSTALL_DATA = $(INSTALL) -m0644
INSTALL_PROGRAM = $(INSTALL) -m0755

ENCVAL_PKCS_7_PRIVATE_KEY ?= .keys/private_key.pkcs7.pem
ENCVAL_PKCS_7_PUBLIC_KEY ?= .keys/public_key.pkcs7.pem
ENCVAL_ENCRYPT_METHOD ?= pkcs7

ENCVAL ?= encval
EYAML ?= eyaml
EYAML_ARGS := --pkcs7-private-key=$(ENCVAL_PKCS_7_PRIVATE_KEY) \
							--pkcs7-public-key=$(ENCVAL_PKCS_7_PUBLIC_KEY) \
							--encrypt-method=$(ENCVAL_ENCRYPT_METHOD)

TEST_PAYLOAD := $(shell date)

export ENCVAL_PKCS_7_PRIVATE_KEY ENCVAL_PKCS_7_PUBLIC_KEY ENCVAL_ENCRYPT_METHOD

.PHONY: all
all: eyaml-compat encval-compat smoke

$(ENCVAL_PKCS_7_PRIVATE_KEY) $(ENCVAL_PKCS_7_PUBLIC_KEY):
	$(ENCVAL) create-keys

.PHONY: eyaml-compat
eyaml-compat: $(ENCVAL_PKCS_7_PRIVATE_KEY) $(ENCVAL_PKCS_7_PUBLIC_KEY)
	echo "$(TEST_PAYLOAD)" \
	| $(ENCVAL) encrypt - \
	| $(EYAML) decrypt $(EYAML_ARGS) --stdin \
	| ./bin/compare-output.sh "$(TEST_PAYLOAD)"

	echo "$(TEST_PAYLOAD)" \
	| $(EYAML) encrypt $(EYAML_ARGS) --stdin --output=string \
	| $(ENCVAL) decrypt - \
	| ./bin/compare-output.sh "$(TEST_PAYLOAD)"

.PHONY: encval-compat
encval-compat: $(ENCVAL_PKCS_7_PRIVATE_KEY) $(ENCVAL_PKCS_7_PUBLIC_KEY)
	echo "$(TEST_PAYLOAD)" \
	| $(ENCVAL) encrypt - \
	| $(ENCVAL) decrypt - \
	| ./bin/compare-output.sh "$(TEST_PAYLOAD)"

.PHONY: smoke
smoke: $(ENCVAL_PKCS_7_PRIVATE_KEY) $(ENCVAL_PKCS_7_PUBLIC_KEY)
	$(ENCVAL) --help
	$(ENCVAL) version
	echo $(TEST_PAYLOAD) | $(ENCVAL) encrypt -

.PHONY: clean
clean:
	$(RM) -r $(ENCVAL_PKCS_7_PRIVATE_KEY) $(ENCVAL_PKCS_7_PUBLIC_KEY)

