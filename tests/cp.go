// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/option"
)

// Cp tests `finch cp` command to copy files between container and host filesystems.
func Cp(o *option.Option) {
	filename := "test-file"
	content := "test-content"
	containerFilepath := filepath.Join("/tmp", filename)
	containerResource := fmt.Sprintf("%s:%s", testContainerName, containerFilepath)

	ginkgo.Describe("copy from container to host and vice versa", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.When("the container is running", func() {
			ginkgo.BeforeEach(func() {
				command.Run(o, "run", "-d", "--name", testContainerName, defaultImage, "sleep", "infinity")
			})

			ginkgo.It("should be able to copy file from host to container", func() {
				path := ffs.CreateTempFile(filename, content)
				ginkgo.DeferCleanup(os.RemoveAll, filepath.Dir(path))

				command.Run(o, "cp", path, containerResource)
				fileShouldExistInContainer(o, testContainerName, containerFilepath, content)
			})

			ginkgo.It("should be able to copy file from container to host", func() {
				cmd := fmt.Sprintf("echo -n %s > %s", content, containerFilepath)
				command.Run(o, "exec", testContainerName, "sh", "-c", cmd)
				fileDir := ffs.CreateTempDir("finch-test")
				path := filepath.Join(fileDir, filename)
				ginkgo.DeferCleanup(os.RemoveAll, fileDir)

				command.Run(o, "cp", containerResource, path)
				fileShouldExist(path, content)
			})

			for _, link := range []string{"-L", "--follow-link"} {
				link := link
				ginkgo.It(fmt.Sprintf("with %s flag, should be able to copy file from host to container and follow symbolic link",
					link), func() {
					path := ffs.CreateTempFile(filename, content)
					fileDir := filepath.Dir(path)
					ginkgo.DeferCleanup(os.RemoveAll, fileDir)
					symlink := filepath.Join(fileDir, "symlink")
					err := os.Symlink(path, symlink)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())

					command.Run(o, "cp", link, symlink, containerResource)
					fileShouldExistInContainer(o, testContainerName, containerFilepath, content)
				})

				ginkgo.It(fmt.Sprintf("with %s flag, should be able to copy file from container to host and follow symbolic link",
					link), func() {
					cmd := fmt.Sprintf("echo -n %s > %s", content, containerFilepath)
					command.Run(o, "exec", testContainerName, "sh", "-c", cmd)
					containerSymlink := filepath.Join("/tmp", "symlink")
					command.Run(o, "exec", testContainerName, "ln", "-s", containerFilepath, containerSymlink)
					fileDir := ffs.CreateTempDir("finch-test")
					path := filepath.Join(fileDir, filename)
					ginkgo.DeferCleanup(os.RemoveAll, fileDir)

					command.Run(o, "cp", link, fmt.Sprintf("%s:%s", testContainerName, containerSymlink), path)
					fileShouldExist(path, content)
				})
			}

			ginkgo.It("should not be able to copy nonexistent file from host to container", func() {
				fileDir := ffs.CreateTempDir("finch-test")
				ginkgo.DeferCleanup(os.RemoveAll, fileDir)

				command.RunWithoutSuccessfulExit(o, "cp", filepath.Join(fileDir, filename), containerResource)
				fileShouldNotExistInContainer(o, testContainerName, containerFilepath)
			})

			ginkgo.It("should not be able to copy nonexistent file from container to host", func() {
				fileDir := ffs.CreateTempDir("finch-test")
				path := filepath.Join(fileDir, filename)
				ginkgo.DeferCleanup(os.RemoveAll, fileDir)

				command.RunWithoutSuccessfulExit(o, "cp", containerResource, path)
				fileShouldNotExist(path)
			})
		})

		ginkgo.When("the container is not running", func() {
			ginkgo.It("should be able to copy file from host to container", func() {
				command.Run(o, "run", "--name", testContainerName, defaultImage, "sleep", "1")
				command.Run(o, "stop", testContainerName)
				path := ffs.CreateTempFile(filename, content)
				ginkgo.DeferCleanup(os.RemoveAll, filepath.Dir(path))
				command.Run(o, "cp", path, containerResource)

				// Need to run container to cat file, can't exec in stopped container.
				// Start here will sleep 1s again so we can check file in container.
				command.Run(o, "container", "start", testContainerName)
				fileShouldExistInContainer(o, testContainerName, containerFilepath, content)
			})

			ginkgo.It("should be able to copy file from container to host", func() {
				cmd := fmt.Sprintf("echo -n %s > %s", content, containerFilepath)
				command.Run(o, "run", "--name", testContainerName, defaultImage, "sh", "-c", cmd)
				fileDir := ffs.CreateTempDir("finch-test")
				path := filepath.Join(fileDir, filename)
				ginkgo.DeferCleanup(os.RemoveAll, fileDir)
				command.Run(o, "cp", containerResource, path)
				fileShouldExist(path, content)
			})
		})
	})
}
