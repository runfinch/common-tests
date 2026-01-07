// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package testutil provides utility functions for the common tests.
package testutil

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/onsi/ginkgo/v2"

	"github.com/runfinch/common-tests/option"
)

// RequireNerdctlVersion skips a test if the nerdctl version does not satisfy the constraint.
func RequireNerdctlVersion(o *option.Option, constraint string) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("failed to construct constraint from %s: %s", constraint, err.Error()))
		return
	}
	nerdctlVersion, err := o.GetNerdctlVersion()
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("failed to get nerdctl version: %s", err.Error()))
		return
	}
	v, err := semver.NewVersion(nerdctlVersion)
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("failed to construct semver from %s: %s", nerdctlVersion, err.Error()))
		return
	}
	if !c.Check(v) {
		ginkgo.Skip(fmt.Sprintf("nerdctl version %s does not satisfy constraint %s", v, constraint))
	}
}
