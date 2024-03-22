// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/gbytes"
	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// VolumePrune tests "volume prune" command that removes all unused volumes.
func VolumePrune(o *option.Option) {
	ginkgo.Describe("remove all unused volumes", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should not remove a volume if it is used by a container", func() {
			command.Run(o, "run", "-v", fmt.Sprintf("%s:/tmp", testVolumeName), "--name", testContainerName, localImages["defaultImage"])
			command.Run(o, "volume", "prune", "--force", "--all")
			volumeShouldExist(o, testVolumeName)
		})

		ginkgo.It("should remove all unused volumes with inputting y in prompt confirmation", func() {
			command.Run(o, "volume", "create", testVolumeName)
			command.New(o, "volume", "prune", "--all").WithStdin(gbytes.BufferWithBytes([]byte("y"))).Run()
			volumeShouldNotExist(o, testVolumeName)
		})

		ginkgo.It("should not remove all unused volumes with inputting n in prompt confirmation", func() {
			command.Run(o, "volume", "create", testVolumeName)
			command.New(o, "volume", "prune", "--all").WithStdin(gbytes.BufferWithBytes([]byte("n"))).Run()
			volumeShouldExist(o, testVolumeName)
		})

		for _, force := range []string{"--force", "-f"} {
			force := force
			ginkgo.It(fmt.Sprintf("should remove all unused volumes without prompting for confirmation with %s flag", force), func() {
				command.Run(o, "volume", "create", testVolumeName)
				command.Run(o, "volume", "prune", force, "--all")
				volumeShouldNotExist(o, testVolumeName)
			})
		}
	})
}
