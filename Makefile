## Copyright (c) 2020 Kevin L. Mitchell
##
## Licensed under the Apache License, Version 2.0 (the "License"); you
## may not use this file except in compliance with the License.  You
## may obtain a copy of the License at
##
##      http://www.apache.org/licenses/LICENSE-2.0
##
## Unless required by applicable law or agreed to in writing, software
## distributed under the License is distributed on an "AS IS" BASIS,
## WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
## implied.  See the License for the specific language governing
## permissions and limitations under the License.

# Packages to test; can be overridden at the command line
PACKAGES    = ./...

# File for repository ignores
IGNORE      = .gitignore

# Get the module root and name
PKG_ROOT    = $(shell grep '^module ' go.mod | awk '{print $$NF}')
PKG_NAME    = $(notdir $(PKG_ROOT))

# Tool-related definitions
TOOLDIR     = tools
TOOLS       =

# Names of the various commands
GO          = go
GOFMT       = gofmt
GOIMPORTS   = ./$(TOOLDIR)/goimports
TOOLS       += golang.org/x/tools/cmd/goimports
GOLINT      = ./$(TOOLDIR)/golint
TOOLS       += golang.org/x/lint/golint
OVERCOVER   = ./$(TOOLDIR)/overcover
TOOLS       += github.com/klmitch/overcover
GOVERALLS   = ./$(TOOLDIR)/goveralls
TOOLS       += github.com/mattn/goveralls

# Coverage configuration file
COV_CONF    = .overcover.yaml

# Additional arguments to pass to overcover
COVER_ARGS  = --summary

# Coverage data and report files
COVER_OUT   = coverage.out
COVER_HTML  = coverage.html

# Collect the sources and test data files for dependencies
SOURCES     = $(shell find . -name \*.go -print)
TEST_DATA   = $(shell find . -path '*/testdata/*' -type f -print)

# Macro to convert a source file name to the corresponding expected
# binary name.
BINNAME     = $(patsubst .,$(PKG_NAME),$(notdir $(patsubst %/,%,$(dir $(1)))))

# Identify the binaries to build; searches for source files that are
# "package main" that contain a "main" function.  Binary names will be
# drawn from the directory the files are in.
_mainPkgRE  = ^\s*package\s\s*main\s*\(\#.*\)*$$
_mainFuncRE = ^\s*func\s\s*main(.*$$
BINSRC      = $(shell echo "$(SOURCES)" | xargs grep -H '$(_mainPkgRE)' | awk -F: '{print $$1}' | sort -u | xargs grep -H '$(_mainFuncRE)' | awk -F: '{print $$1}' | sort -u)
BINS        = $(call BINNAME,$(BINSRC))

# Files to be cleaned up on "make clean"
CLEAN       = $(BINS) $(COVER_OUT) $(COVER_HTML) $(IGNORE).tmp

# CI-linked variables; these set up read-only behavior within a CI
# system
ifeq ($(CI),true)
FORMAT_TARG = format-test
MOD_ARG     = -mod=readonly
COV_ARG     = --readonly
else
FORMAT_TARG = format
MOD_ARG     =
COV_ARG     =
endif

# Compute the dependencies for the "all" target
ALL_TARG    = $(IGNORE) test
ifneq ($(BINS),)
ALL_TARG    += build
endif

# Set up dependencies for the "test" and "cover" targets
TEST_TARG   = $(FORMAT_TARG) lint vet test-only

include $(wildcard scripts/*.mk)

# Add TOOLDIR to CLEAN unless told not to
ifneq ($(KEEP_TOOLDIR),true)
CLEAN       += $(TOOLDIR)
endif

all: $(ALL_TARG) ## Run tests and build binaries (if any)

build: $(BINS) ## Build binaries (if any)

tidy: ## Ensure go.mod matches the source code
	$(GO) mod tidy

format-test: $(GOIMPORTS) ## Check that source files are properly formatted
	@all=`( \
		$(GOIMPORTS) -l -local $(PKG_ROOT) $(SOURCES); \
		$(GOFMT) -l -s $(SOURCES) \
	) | sort -u`; \
	if [ "$$all" != "" ]; then \
		echo "The following files require formatting updates:"; \
		echo; \
		echo "$$all"; \
		echo; \
		echo "Use \"make format\" to reformat these files."; \
		exit 1; \
	fi

format: $(GOIMPORTS) ## Format source files properly
	$(GOIMPORTS) -l -local $(PKG_ROOT) -w $(SOURCES)
	$(GOFMT) -l -s -w $(SOURCES)

lint: $(GOLINT) ## Lint-check source files
	$(GOLINT) -set_exit_status $(PACKAGES)

vet: ## Vet source files
	$(GO) vet $(MOD_ARG) $(PACKAGES)

test-only: ## Run tests only
	$(GO) test $(MOD_ARG) -race -coverprofile=$(COVER_OUT) $(PACKAGES)

test: $(TEST_TARG) cover-test ## Run all tests

cover: $(TEST_TARG) $(COVER_HTML) cover-test ## Run tests and generate a coverage report

cover-test: $(COVER_OUT) $(OVERCOVER) ## Test that coverage meets minimum configured threshold
	$(OVERCOVER) --config $(COV_CONF) $(COV_ARG) --coverprofile $(COVER_OUT) $(COVER_ARGS) $(PACKAGES)

# Travis-specific target for submitting coverage to coveralls.io;
# explicitly undocumented
goveralls: $(COVER_OUT) $(GOVERALLS)
	$(GOVERALLS) -coverprofile=$(COVER_OUT) -service=travis-ci

$(COVER_OUT): $(SOURCES) $(TEST_DATA)
	$(MAKE) test-only

$(COVER_HTML): $(COVER_OUT)
	$(GO) tool cover -html=$(COVER_OUT) -o $(COVER_HTML)

clean: ## Clean up intermediate files
	rm -rf $(CLEAN)

# Sets up build targets for each binary
ifneq ($(BINS),)
$(BINS): $(SOURCES)

define BIN_template =
$$(call BINNAME,$(1)):
	$(GO) build $(MOD_ARG) -o $$(call BINNAME,$(1)) $(1)
endef

$(foreach bin,$(BINSRC),$(eval $(call BIN_template,$(bin))))
endif

# Sets up the tools directory
$(TOOLDIR):
	mkdir $(TOOLDIR)
	cd $(TOOLDIR) && go mod init $(PKG_ROOT)/$(notdir $(TOOLDIR))

# Sets up build targets for each required tool
define TOOL_template =
./$(TOOLDIR)/$$(notdir $(1)): $(TOOLDIR)
	cd $(TOOLDIR) && GOBIN=$(abspath $(TOOLDIR)) go install $(1)
endef

$(foreach tool,$(TOOLS),$(eval $(call TOOL_template,$(tool))))

$(IGNORE).tmp: $(MAKEFILE_LIST)
	echo $(CLEAN) | sed 's/ /\n/g' > $(IGNORE).tmp

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
		sort -u | awk -F: '{ \
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

.PHONY: all build tidy format-test format lint vet test-only test cover cover-test goveralls clean help
