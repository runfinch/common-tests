// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ffs

import (
	"os"
)

// CheckIfFileExists checks if a file exists in the filesystem.
func CheckIfFileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}
	return false
}
