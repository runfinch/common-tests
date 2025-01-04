// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package option

import "regexp"

const (
	nerdctl1xx            = "1.x.x"
	nerdctl2xx            = "2.x.x"
	defaultNerdctlVersion = nerdctl2xx
)

var (
	nerdctl1xxRegex = regexp.MustCompile(`^1\.[x0-9]+\.[x0-9]+`)
	nerdctl2xxRegex = regexp.MustCompile(`^2\.[x0-9]+\.[x0-9]+`)
)

func isNerdctl1xx(version string) bool {
	return nerdctl1xxRegex.MatchString(version)
}

func isNerdctl2xx(version string) bool {
	return nerdctl2xxRegex.MatchString(version)
}
