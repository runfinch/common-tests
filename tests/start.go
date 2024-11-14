// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"strings"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Start tests starting a container.
func Start(o *option.Option) {
	ginkgo.Describe("start a container", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should start the container if it is in Exited status", func() {
			command.Run(o, "run", "-d", "--name", testContainerName, localImages[defaultImage], "nc", "-l")
			containerShouldBeRunning(o, testContainerName)

			command.Run(o, "stop", testContainerName)
			command.RunWithoutSuccessfulExit(o, "exec", testContainerName, "echo", "foo")

			command.Run(o, "start", testContainerName)
			containerShouldBeRunning(o, testContainerName)
		})

		for _, attach := range []string{"--attach", "-a", "-a=true", "--attach=true"} {
			attach := attach
			ginkgo.It(fmt.Sprintf("with %s flag, should start the container with stdout", attach), func() {
				command.Run(o, "create", "--name", testContainerName, localImages[defaultImage], "echo", "foo")
				output := command.StdoutStr(o, "start", attach, testContainerName)
				gomega.Expect(output).To(gomega.Equal("foo"))
			})
		}

		ginkgo.It("should run a container without an init process when --init=false flag is used", func() {
			command.Run(o, "run", "--name", testContainerName, "--init=false", localImages[defaultImage], "ps", "-ao", "pid,comm")
			psOutput := command.StdoutStr(o, "logs", testContainerName)

			// Split the output into lines
			lines := strings.Split(strings.TrimSpace(psOutput), "\n")

			processLine := lines[1] // Second line (after header)
			fields := strings.Fields(processLine)

			pid := fields[0]
			command := fields[1]
			gomega.Expect(pid).To(gomega.Equal("1"), "The only process should have PID 1")
			gomega.Expect(command).To(gomega.Equal("ps"), "The only process should be ps")

			// Verify there's no init process
			gomega.Expect(psOutput).NotTo(gomega.ContainSubstring("tini"), "There should be no tini process")
		})
	})
}
