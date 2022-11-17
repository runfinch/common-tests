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

// VolumeCreate tests "volume create" command that creates a volume.
func VolumeCreate(o *option.Option) {
	ginkgo.Describe("create a volume", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should create a volume with name", func() {
			command.Run(o, "volume", "create", testVolumeName)
			volumeShouldExist(o, testVolumeName)
		})

		ginkgo.It("data in volume should be shared between containers", func() {
			command.Run(o, "volume", "create", testVolumeName)
			command.Run(o, "run", "-v", fmt.Sprintf("%s:/tmp", testVolumeName), defaultImage, "sh", "-c", "echo foo > /tmp/test.txt")
			output := command.StdoutStr(o, "run", "-v", fmt.Sprintf("%s:/tmp", testVolumeName), defaultImage, "cat", "/tmp/test.txt")
			gomega.Expect(output).Should(gomega.Equal("foo"))
		})

		ginkgo.It("should create a volume with label with --label flag", func() {
			command.Run(o, "volume", "create", "--label", "label=tag", testVolumeName)
			output := command.StdoutStr(o, "volume", "inspect", testVolumeName, "--format", "{{.Labels.label}}")
			gomega.Expect(output).Should(gomega.Equal("tag"))
		})

		ginkgo.It("should create multiple labels with --label flag", func() {
			command.Run(o, "volume", "create", "--label", "label=tag", "--label", "label1=tag1", testVolumeName)
			tag := command.StdoutStr(o, "volume", "inspect", testVolumeName, "--format", "{{.Labels.label}}")
			tag1 := command.StdoutStr(o, "volume", "inspect", testVolumeName, "--format", "{{.Labels.label1}}")
			gomega.Expect(tag).Should(gomega.Equal("tag"))
			gomega.Expect(tag1).Should(gomega.Equal("tag1"))
		})

		ginkgo.It("should not create a volume if the volume with the same name exists", func() {
			command.Run(o, "volume", "create", testVolumeName)
			command.RunWithoutSuccessfulExit(o, "volume", "create", testVolumeName)
		})
	})
}
