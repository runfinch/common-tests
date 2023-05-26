// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package fenv Package path implements utility routines for working with environment variables
package fenv

import (
	"os"
)

// GetEnv retrieves the value of an environment variable.
// It returns an empty string if the variable is not set.
func GetEnv(key string) string {
	return os.Getenv(key)
}
