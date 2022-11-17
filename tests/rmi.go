// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"

	"github.com/onsi/ginkgo/v2"
)

// Rmi tests removing a container image.
func Rmi(o *option.Option) {
	ginkgo.Describe("remove a container image", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should remove an image when container is not running", func() {
			pullImage(o, defaultImage)

			command.Run(o, "rmi", defaultImage)
			imageShouldNotExist(o, defaultImage)
		})

		ginkgo.Context("when there is a container based on the image to be removed", func() {
			ginkgo.BeforeEach(func() {
				pullImage(o, defaultImage)
				command.Run(o, "run", defaultImage)
			})

			ginkgo.It("should not be able to remove the image without -f/--force flag", func() {
				command.RunWithoutSuccessfulExit(o, "rmi", defaultImage)
				imageShouldExist(o, defaultImage)
			})

			for _, force := range []string{"-f", "--force"} {
				force := force
				ginkgo.It(fmt.Sprintf("should be able to remove the image with %s flag", force), func() {
					command.Run(o, "rmi", force, defaultImage)
					imageShouldNotExist(o, defaultImage)
				})
			}
		})
	})
}
