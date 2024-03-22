// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"

	"github.com/onsi/ginkgo/v2"
)

// Pull tests pulling a container image.
func Pull(o *option.Option) {
	ginkgo.Describe("pull a container image", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveImages(o)
		})

		ginkgo.AfterEach(func() {
			command.RemoveImages(o)
		})

		ginkgo.It("should pull the default image successfully", func() {
			command.Run(o, "pull", localImages[defaultImage])
			imageShouldExist(o, localImages[defaultImage])
		})
	})
}
