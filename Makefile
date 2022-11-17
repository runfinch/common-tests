# Copyright Amazon.com, Inc. or its affiliates.
# SPDX-License-Identifier: Apache-2.0

# Note that Finch CLI may not be an option here, or there will be circular dependency.
# To be specific, when we update e2e tests in the future,
# it is possible that the functionality has not been implemented in Finch CLI yet,
# so the test logs won't be the expected ones.
#
# For example, if the test is "pulls an alpine image",
# then the test subject should pull an alpine image successfully and the corresponding logs should be printed.
# If the functionality is not in the test subject yet, the developer will see some failing logs instead.
#
# The assumption is that when adding/updating an e2e test,
# it is likely that the publicly available Finch CLI does not have the corresponding functionality
# (i.e., the "pulls an alpine image" in the example above) yet.
# It is because that the functionality cannot be added to Finch CLI until the corresponding test is in place.
# In other words, the test has to be merged into Finch Test and used by Finch CLI, then
# the CI of Finch CLI can run with the updated test and the added functionality.
#
# As a result, we may want to use another reference implementation (i.e., nerdctl here) when testing the tests themselves.
SUBJECT ?= nerdctl

VERBOSE ?= true
VERBOSE_FLAGS =
ifeq ($(VERBOSE),true)
    # https://github.com/onsi/ginkgo/issues/381
	VERBOSE_FLAGS = -test.v -ginkgo.v
endif

.PHONY: run
run:
	go test -timeout 30m ./run/... $(VERBOSE_FLAGS) -args --subject="$(SUBJECT)"

.PHONY: lint
# To run golangci-lint locally: https://golangci-lint.run/usage/install/#local-installation
lint:
	golangci-lint run
