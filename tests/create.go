// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Create tests creating a container.
func Create(o *option.Option) {
	ginkgo.Describe("create a container", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should create a container and able to start the container", func() {
			command.Run(o, "create", "--name", testContainerName, defaultImage, "sleep", "infinity")
			status := command.StdoutStr(o, "ps", "-a", "--filter", fmt.Sprintf("name=%s", testContainerName), "--format", "{{.Status}}")
			gomega.Expect(status).Should(gomega.Equal("Created"))

			command.Run(o, "start", testContainerName)
			containerShouldBeRunning(o, testContainerName)
		})

		ginkgo.It("should not create a container if the image doesn't exist", func() {
			command.RunWithoutSuccessfulExit(o, "create", nonexistentImageName)
		})
	})
}
