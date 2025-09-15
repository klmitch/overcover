# Packages to test; can be overridden at the command line
PACKAGES           = ./...

# File for repository ignores
IGNORE             = .gitignore

# Additional arguments to pass to various target rules
BUILD_ARGS         ?=
TEST_REPORT_ARGS   ?= --format testname
TEST_ARGS          ?=
LINT_ARGS          ?=
FORMAT_ARGS        ?=
COVER_ARGS         ?= --summary
COVER_REPORT_ARGS  ?=
COVER_XML_ARGS     ?=

# Get the module root and name
PKG_ROOT           = $(shell grep '^module ' go.mod | awk '{print $$NF}')
PKG_NAME           = $(notdir $(PKG_ROOT))

# Tool-related definitions
TOOLDIR            = .tools
TOOLS              =
TOOLS_CONF         = .tools.conf
TOOLS_VERSION      := version() { \
	vers=latest; \
	if [ -r "$(TOOLS_CONF)" ]; then \
		tmp=`grep "$$1" "$(TOOLS_CONF)" | awk '{print $$2}'`; \
		if [ "$${tmp}" != "" ]; then \
			vers="$${tmp}"; \
		fi; \
	fi; \
	echo "$${vers}" ; \
}

# Names of the various commands
GO                 = go
GOLANGCI_LINT      = ./$(TOOLDIR)/golangci-lint
GOTESTSUM          = ./$(TOOLDIR)/gotestsum
TOOLS              += gotest.tools/gotestsum
OVERCOVER          = ./$(TOOLDIR)/overcover
TOOLS              += github.com/klmitch/overcover
COBERTURA          = ./$(TOOLDIR)/gocover-cobertura
TOOLS              += github.com/boumenot/gocover-cobertura

# Coverage configuration file
COV_CONF           = .overcover.yaml

# Linter configuration file and default list of linters to enable if
# generating it
LINT_CONF          = .golangci.yml
LINT_LINTERS       = asciicheck dupl err113 exhaustive exptostd gocognit goconst
LINT_LINTERS       += goprintffuncname gosec misspell nolintlint predeclared
LINT_LINTERS       += promlinter protogetter revive sloglint spancheck unconvert
LINT_LINTERS       += wastedassign whitespace zerologlint
LINT_FORMATTERS    = gofmt gofumpt goimports
LINT_NOTESTS       = dupl goconst
LINT_URL           = https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh

# CI-linked variables; these set up read-only behavior within a CI
# system
ifeq ($(CI),true)
MOD_ARG            = -mod=readonly
FIX_ARG            = --modules-download-mode=readonly
FMT_ARG            = --diff
COV_ARG            = --readonly
else
MOD_ARG            =
FIX_ARG            = --fix
FMT_ARG            =
COV_ARG            =
endif

# Directories within which to place artifacts
GO_BUILD_ARTIFACTS ?= .
GO_TEST_ARTIFACTS  ?= .

# Canonical coverage data and report file names
CCOVER_OUT         = cover.out
CCOVER_HTML        = cover.html
CCOVER_XML         = cover.xml
CJUNIT_OUT         = report.xml

# Macro to simplify a path to be relative to the top-level path.
SIMPPATH           = $(foreach name,$(1),$(subst $(abspath .)/,,$(abspath $(name))))

# Complete file paths for coverage data and file names
COVER_OUT          = $(call SIMPPATH,$(GO_TEST_ARTIFACTS)/$(CCOVER_OUT))
COVER_HTML         = $(call SIMPPATH,$(GO_TEST_ARTIFACTS)/$(CCOVER_HTML))
COVER_XML          = $(call SIMPPATH,$(GO_TEST_ARTIFACTS)/$(CCOVER_XML))
JUNIT_OUT          = $(call SIMPPATH,$(GO_TEST_ARTIFACTS)/$(CJUNIT_OUT))

# Collect the sources and test data files for dependencies; this also
# collects the list of sources that are not test files for detecting
# binaries and plugins to build
GO_MOD             = ./go.mod ./go.sum
SOURCES            = $(shell find . -name \*.go -print)
SRC_ONLY           = $(filter-out %_test.go,$(SOURCES))
TEST_DATA          = $(shell find . -path '*/testdata/*' -type f -print)

# Macros to convert a source file name to the corresponding expected
# binary or plugin name.
BINNAME            = $(call SIMPPATH,$(dir $(1))$(patsubst .,$(PKG_NAME),$(notdir $(patsubst %/,%,$(dir $(1))))))
FULLBINNAME        = $(call SIMPPATH,$(addprefix $(GO_BUILD_ARTIFACTS)/,$(call BINNAME,$(1))))

# The extension for the plugin.  As only UNIX-style OSes currently
# support golang plugins anyway, we can fix this to "so".
PLUG_EXT           = so

# Identify the binaries and plugins to build.  This starts by
# searching the non-test sources for source files that are "package
# main".  Binaries will have a "main" function, while plugins will
# have the special comment "//make:plugin".  Binary and plugin names
# will be drawn from the directory the files are in.  The outputs are
# BINSRC, BINS, PLUGSRC, and PLUGS; the CBINS and CPLUGS are
# "canonical" filenames used for ".gitignore" and the "clean" target.
_mainPkgRE         = ^\s*package\s\s*main\s*\(//.*\)*$$
_mainFuncRE        = ^\s*func\s\s*main(.*$$
_pluginRE          = ^\s*//\s*make:plugin\s*\(//.*\)*$$
MAINSRC            = $(shell echo "$(SRC_ONLY)" | xargs grep -H '$(_mainPkgRE)' | awk -F: '{print $$1}' | sort -u)
BINSRC             = $(shell echo "$(MAINSRC)" | xargs grep -H '$(_mainFuncRE)' | awk -F: '{print $$1}' | sort -u)
CBINS              = $(foreach bin,$(BINSRC),$(call BINNAME,$(bin)))
BINS               = $(foreach bin,$(BINSRC),$(call FULLBINNAME,$(bin)))
PLUGSRC            = $(shell echo "$(MAINSRC)" | xargs grep -H '$(_pluginRE)' | awk -F: '{print $$1}' | sort -u)
CPLUGS             = $(foreach plug,$(PLUGSRC),$(call BINNAME,$(plug)).$(PLUG_EXT))
PLUGS              = $(foreach plug,$(PLUGSRC),$(call FULLBINNAME,$(plug)).$(PLUG_EXT))

# Files to be cleaned up on "make clean"
CLEAN              = $(BINS) $(PLUGS) $(COVER_OUT) $(COVER_HTML) $(COVER_XML) $(JUNIT_OUT) $(IGNORE).tmp $(LINT_CONF).tmp $(TOOLDIR)

# Files to be ignored by git
IGNORE_FILES       = $(addprefix /,$(CBINS) $(CPLUGS) $(CCOVER_OUT) $(CCOVER_HTML) $(CCOVER_XML) $(CJUNIT_OUT) $(IGNORE).tmp $(LINT_CONF).tmp $(TOOLDIR))

# Compute the dependencies for the "all" and "build" targets
ALL_TARG           = $(IGNORE) test
BUILD_TARG         =
ifneq ($(BINS),)
BUILD_TARG         += build-bins
endif
ifneq ($(PLUGS),)
BUILD_TARG         += build-plugins
endif
ifneq ($(BUILD_TARG),)
ALL_TARG           += build
endif

# Set up dependencies for the "test" and "cover" targets
TEST_TARG          = lint test-only

include $(wildcard scripts/*.mk)

all: $(ALL_TARG) ## Run tests and build binaries and plugins (if any)

build: $(BUILD_TARG) ## Build binaries and plugins (if any)

build-bins: $(BINS) ## Build binaries (if any)

build-plugins: $(PLUGS) ## Build plugins (if any)

tidy: ## Ensure go.mod matches the source code
	$(GO) mod tidy

format: $(GOLANGCI_LINT) $(LINT_CONF) ## Reformat source files
	$(GOLANGCI_LINT) fmt -c $(LINT_CONF) $(FMT_ARG) $(FORMAT_ARGS) $(PACKAGES)

lint: $(GOLANGCI_LINT) $(LINT_CONF) ## Lint-check source files; may fix some lint issues
	$(GOLANGCI_LINT) run -c $(LINT_CONF) $(FIX_ARG) $(LINT_ARGS) $(PACKAGES)

test-only: $(GO_TEST_ARTIFACTS) $(GOTESTSUM) ## Run tests only
	$(GOTESTSUM) $(TEST_REPORT_ARGS) --junitfile $(JUNIT_OUT) -- $(MOD_ARG) $(TEST_ARGS) -coverprofile=$(COVER_OUT) -coverpkg=./... $(PACKAGES)

test: $(TEST_TARG) cover-test ## Run all tests

cover: $(TEST_TARG) cover-report cover-test ## Run tests and generate a coverage report

cover-report: $(COVER_HTML) $(COVER_XML) ## Generate a coverage report, running tests only if required

cover-test: $(COVER_OUT) $(OVERCOVER) ## Test that coverage meets minimum configured threshold
	$(OVERCOVER) --config $(COV_CONF) $(COV_ARG) --coverprofile $(COVER_OUT) $(COVER_ARGS) $(PACKAGES)

clean: ## Clean up intermediate files
	rm -rf $(CLEAN)

list-updates: ## List direct dependencies that are not up to date
	@$(GO) list -m -u -f '{{if and (not .Indirect) (not .Main) .Update}}{{.}}{{end}}' all

update: ## Check for available direct dependency updates and apply them
	for mod in $$($(GO) list -m -u -f '{{if and (not .Indirect) (not .Main) .Update}}{{.Path}}{{end}}' all); do \
		$(GO) get -u $${mod}; \
	done; \
	$(GO) mod tidy

$(LINT_CONF).tmp: $(MAKEFILE_LIST)
	@echo "Generating $(LINT_CONF).tmp"
	@echo 'version: "2"' > $(LINT_CONF).tmp; \
	echo "linters:" >> $(LINT_CONF).tmp; \
	echo "  enable:" >> $(LINT_CONF).tmp; \
	for linter in $(LINT_LINTERS); do \
		echo "  - $${linter}" >> $(LINT_CONF).tmp; \
	done; \
	echo "  exclusions:" >> $(LINT_CONF).tmp; \
	echo "    generated: lax" >> $(LINT_CONF).tmp; \
	echo "    presets:" >> $(LINT_CONF).tmp; \
	echo "    - comments" >> $(LINT_CONF).tmp; \
	echo "    - common-false-positives" >> $(LINT_CONF).tmp; \
	echo "    - legacy" >> $(LINT_CONF).tmp; \
	echo "    - std-error-handling" >> $(LINT_CONF).tmp; \
	echo "    rules:" >> $(LINT_CONF).tmp; \
	echo "    - linters:" >> $(LINT_CONF).tmp; \
	for linter in $(LINT_NOTESTS); do \
		echo "      - $${linter}" >> $(LINT_CONF).tmp; \
	done; \
	echo '      path: (.+)_test\.go' >> $(LINT_CONF).tmp; \
	echo "    paths:" >> $(LINT_CONF).tmp; \
	echo "    - third_party$$" >> $(LINT_CONF).tmp; \
	echo "    - builtin$$" >> $(LINT_CONF).tmp; \
	echo "    - examples$$" >> $(LINT_CONF).tmp; \
	echo "severity:" >> $(LINT_CONF).tmp; \
	echo "  default: blocker" >> $(LINT_CONF).tmp; \
	echo "formatters:" >> $(LINT_CONF).tmp; \
	echo "  enable:" >> $(LINT_CONF).tmp; \
	for formatter in $(LINT_FORMATTERS); do \
		echo "  - $${formatter}" >> $(LINT_CONF).tmp; \
	done; \
	echo "  settings:" >> $(LINT_CONF).tmp; \
	echo "    gofumpt:" >> $(LINT_CONF).tmp; \
	echo "      module-path: $(PKG_ROOT)" >> $(LINT_CONF).tmp; \
	echo "    goimports:" >> $(LINT_CONF).tmp; \
	echo "      local-prefixes:" >> $(LINT_CONF).tmp; \
	echo "      - $(PKG_ROOT)" >> $(LINT_CONF).tmp; \
	echo "  exclusions:" >> $(LINT_CONF).tmp; \
	echo "    generated: lax" >> $(LINT_CONF).tmp; \
	echo "    paths:" >> $(LINT_CONF).tmp; \
	echo "    - third_party$$" >> $(LINT_CONF).tmp; \
	echo "    - builtin$$" >> $(LINT_CONF).tmp; \
	echo "    - examples$$" >> $(LINT_CONF).tmp; \
	echo "output:" >> $(LINT_CONF).tmp; \
	echo "  formats:" >> $(LINT_CONF).tmp; \
	echo "    text:" >> $(LINT_CONF).tmp; \
	echo "      path: stdout" >> $(LINT_CONF).tmp; \
	echo "      print-issued-lines: false" >> $(LINT_CONF).tmp; \
	echo "      colors: true" >> $(LINT_CONF).tmp

$(LINT_CONF): $(LINT_CONF).tmp
ifeq ($(CI),true)
	@if cmp $(LINT_CONF) $(LINT_CONF).tmp >/dev/null 2>&1; then \
		:; \
	else \
		echo "The $(LINT_CONF) file requires regeneration."; \
		echo "Use \"make $(LINT_CONF)\" to regenerate it."; \
		exit 1; \
	fi
else
	cp $(LINT_CONF).tmp $(LINT_CONF)
endif

$(COVER_OUT): $(GO_MOD) $(SOURCES) $(TEST_DATA)
	$(MAKE) test-only

$(COVER_HTML): $(COVER_OUT)
	$(GO) tool cover -html=$(COVER_OUT) -o $(COVER_HTML) $(COVER_REPORT_ARGS)

$(COVER_XML): $(COBERTURA) $(COVER_OUT)
	$(COBERTURA) $(COVER_XML_ARGS) < $(COVER_OUT) > $(COVER_XML)

# Sets up build targets for each binary
ifneq ($(BINS),)
$(BINS): $(GO_MOD) $(SOURCES)

define BIN_template =
$$(call FULLBINNAME,$(1)):
	$(GO) build $(MOD_ARG) $(BUILD_ARGS) -o $$(call FULLBINNAME,$(1)) $(1)
endef

$(foreach bin,$(BINSRC),$(eval $(call BIN_template,$(bin))))
endif

# Sets up build targets for each plugin
ifneq ($(PLUGS),)
$(PLUGS): $(GO_MOD) $(SOURCES)

define PLUG_template =
$$(call FULLBINNAME,$(1)).so:
	$(GO) build -buildmode=plugin $(MOD_ARG) $(BUILD_ARGS) -o $$(call FULLBINNAME,$(1)).so $(1)
endef

$(foreach plug,$(PLUGSRC),$(eval $(call PLUG_template,$(plug))))
endif

# Sets up the test artifacts directory
$(GO_TEST_ARTIFACTS):
	mkdir -p $(GO_TEST_ARTIFACTS)

# Sets up the tools directory
$(TOOLDIR):
	mkdir $(TOOLDIR)

# Ensures that golangci-lint is available
$(GOLANGCI_LINT): $(TOOLDIR)
	$(TOOLS_VERSION); \
	LINT_VERSION=`version $$(basename $(GOLANGCI_LINT))`; \
	if command -v wget; then \
		echo "Installing golangci-lint version $${LINT_VERSION} with wget"; \
		wget -O- -nv $(LINT_URL) | sh -s -- -b $(TOOLDIR) $${LINT_VERSION}; \
	elif command -v curl; then \
		echo "Installing golangci-lint version $${LINT_VERSION} with curl"; \
		curl -sSfL $(LINT_URL) | sh -s -- -b $(TOOLDIR) $${LINT_VERSION}; \
	else \
		echo "Install curl or wget" >&2; \
		exit 1; \
	fi

# Sets up build targets for each required tool
define TOOL_template =
./$(TOOLDIR)/$$(notdir $(1)): $(TOOLDIR)
	$$(TOOLS_VERSION); \
	vers=`version $$$$(basename $(1))`; \
	echo "Installing $(1) version $$$${vers}"; \
	GOBIN=$(abspath $(TOOLDIR)) go install "$(1)@$$$${vers}"
endef

$(foreach tool,$(TOOLS),$(eval $(call TOOL_template,$(tool))))

$(IGNORE).tmp: $(MAKEFILE_LIST)
	echo $(IGNORE_FILES) | sed 's/ /\n/g' | sed 's@^\./@@g' > $(IGNORE).tmp

$(IGNORE): $(IGNORE).tmp
ifeq ($(CI),true)
	@if cmp $(IGNORE) $(IGNORE).tmp >/dev/null 2>&1; then \
		:; \
	else \
		echo "The $(IGNORE) file requires regeneration."; \
		echo "Use \"make $(IGNORE)\" to regenerate it."; \
		echo "Current contents:"; \
		echo; \
		cat $(IGNORE); \
		echo; \
		echo "Expected contents:"; \
		echo; \
		cat $(IGNORE).tmp; \
		exit 1; \
	fi
else
	cp $(IGNORE).tmp $(IGNORE)
endif

help: ## Emit help for the Makefile
	@echo "Available make targets:"
	@echo
	@grep -h '^[^ 	:].*:.*##' $(MAKEFILE_LIST) | sed 's/:.*## */:/g' | \
		LANG=C sort -u -t: -k1,1 | awk -F: '{ \
			if (length($$1) > width) { \
				width = length($$1); \
			} \
			targets[targetCnt++] = $$1; \
			help[$$1] = $$2; \
		} \
		END { \
			indent = sprintf("\n  %*s  ", width, ""); \
			for (i = 0; i < targetCnt; i++) { \
				target = targets[i]; \
				helpText = help[target]; \
				gsub("\\\\n", indent, helpText); \
				printf("  %-*s  %s\n", width, target, helpText); \
			} \
		}'

debug: # Emit debugging information; target hidden from help
	@echo "   COVER_OUT: $(COVER_OUT)"
	@echo "  COVER_HTML: $(COVER_HTML)"
	@echo "   COVER_XML: $(COVER_XML)"
	@echo "   JUNIT_OUT: $(JUNIT_OUT)"
	@echo "      GO_MOD: $(GO_MOD)"
	@echo "     SOURCES: $(SOURCES)"
	@echo "    SRC_ONLY: $(SRC_ONLY)"
	@echo "   TEST_DATA: $(TEST_DATA)"
	@echo "     MAINSRC: $(MAINSRC)"
	@echo "      BINSRC: $(BINSRC)"
	@echo "       CBINS: $(CBINS)"
	@echo "        BINS: $(BINS)"
	@echo "     PLUGSRC: $(PLUGSRC)"
	@echo "      CPLUGS: $(CPLUGS)"
	@echo "       PLUGS: $(PLUGS)"
	@echo "       CLEAN: $(CLEAN)"
	@echo "IGNORE_FILES: $(IGNORE_FILES)"
	@echo "    ALL_TARG: $(ALL_TARG)"
	@echo "  BUILD_TARG: $(BUILD_TARG)"
	@echo "   TEST_TARG: $(TEST_TARG)"

.PHONY: all build build-bins build-plugins clean cover cover-report cover-test format help lint list-updates test test-only tidy update debug
