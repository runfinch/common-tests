// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package ffs contains functions that manipulate file system
package ffs

import (
	"os"
	"path/filepath"

	"github.com/onsi/gomega"
)

// CreateBuildContext creates a directory which contains a Dockerfile with the specified content and returns the path to the directory.
// It is the caller's responsibility to remove the directory when it is no longer needed.
func CreateBuildContext(dockerfile string) string {
	return filepath.Dir(CreateTempFile("Dockerfile", dockerfile))
}

// WriteFile writs data to the file specified by name. The file will be created if not existing.
func WriteFile(name string, data string) {
	err := os.WriteFile(name, []byte(data), 0o644)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}

// CreateTarFilePath creates a directory and a tar file path appended to the directory and returns the tar file path.
// It is the caller's responsibility to remove the directory when it is no longer needed.
// TODO: replace with CreateTempDir.
func CreateTarFilePath() string {
	tempDir := CreateTempDir("finch-test-save")
	tarFilePath := filepath.Join(tempDir, "test.tar")
	return tarFilePath
}

// CreateComposeYmlContext creates a temp directory along with a docker-compose.yml file.
// It is the caller's responsibility to remove the directory when it is no longer needed.
func CreateComposeYmlContext(composeYmlContent string) (string, string) {
	composeFileName := "docker-compose.yml"
	tempDir := CreateTempDir("finch-compose")
	composeFilePath := filepath.Join(tempDir, composeFileName)
	WriteFile(composeFilePath, composeYmlContent)
	return tempDir, composeFilePath
}

// CreateTempFile creates a temp directory which contains a temp file and returns the path to the temp file.
// It is the caller's responsibility to remove the directory when it is no longer needed.
func CreateTempFile(filename string, content string) string {
	tempDir := CreateTempDir("finch-test")
	filepath := filepath.Join(tempDir, filename)
	WriteFile(filepath, content)
	return filepath
}

// CreateTempDir creates a temp directory and returns the path of the created directory.
// It is the caller's responsibility to remove the directory when it is no longer needed.
func CreateTempDir(directoryPrefix string) string {
	homeDir, err := os.UserHomeDir()
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	tempDir, err := os.MkdirTemp(homeDir, directoryPrefix)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return tempDir
}
