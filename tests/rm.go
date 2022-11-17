// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Rm tests removing a container.
func Rm(o *option.Option) {
	ginkgo.Describe("remove a container", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should remove the container when it is not running", func() {
			command.Run(o, "run", "--name", testContainerName, defaultImage)
			containerShouldExist(o, testContainerName)

			command.Run(o, "rm", testContainerName)
			containerShouldNotExist(o, testContainerName)
		})

		ginkgo.Context("when the container is running", func() {
			ginkgo.BeforeEach(func() {
				command.Run(o, "run", "-d", "--name", testContainerName, defaultImage, "sleep", "infinity")
			})

			ginkgo.It("should not be able to remove the container without -f/--force flag", func() {
				command.RunWithoutSuccessfulExit(o, "rm", testContainerName)
				containerShouldExist(o, testContainerName)
			})

			for _, force := range []string{"-f", "--force"} {
				force := force
				ginkgo.It(fmt.Sprintf("should be able to remove the container with %s flag", force), func() {
					command.Run(o, "rm", force, testContainerName)
					containerShouldNotExist(o, testContainerName)
				})
			}
		})

		ginkgo.Context("when a volume is used by the container", func() {
			for _, volumes := range []string{"-v", "--volumes"} {
				volumes := volumes
				ginkgo.It(fmt.Sprintf("with %s flag, should remove the container and the anonymous volume used by the container", volumes),
					func() {
						command.Run(o, "run", "-v", "/usr/share", "--name", testContainerName, defaultImage)
						anonymousVolume := command.StdoutStr(o, "inspect", testContainerName,
							"--format", "{{range .Mounts}}{{.Name}}{{end}}")
						containerShouldExist(o, testContainerName)
						volumeShouldExist(o, anonymousVolume)
						command.Run(o, "rm", volumes, testContainerName)
						containerShouldNotExist(o, testContainerName)
						volumeShouldNotExist(o, anonymousVolume)
					},
				)

				ginkgo.It(fmt.Sprintf("with %s flag, should remove the container but can't remove the named volume used by container", volumes),
					func() {
						command.Run(o, "run", "-v", "foo:/usr/share", "--name", testContainerName, defaultImage)
						volumeShouldExist(o, "foo")

						command.Run(o, "rm", volumes, testContainerName)
						containerShouldNotExist(o, testContainerName)
						volumeShouldExist(o, "foo")
					},
				)
			}
		})
	})
}
