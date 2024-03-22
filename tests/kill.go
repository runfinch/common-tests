// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Kill tests killing a running container.
func Kill(o *option.Option) {
	ginkgo.Describe("kill a container", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.When("the container is running", func() {
			ginkgo.BeforeEach(func() {
				command.Run(o, "run", "-d", "--name", testContainerName, localImages[defaultImage], "sleep", "infinity")
			})

			ginkgo.It("should kill the running container", func() {
				containerShouldBeRunning(o, testContainerName)
				command.Run(o, "kill", testContainerName)
				command.RunWithoutSuccessfulExit(o, "exec", testContainerName, "echo", "foo")
				containerShouldNotBeRunning(o, testContainerName)
			})

			for _, signal := range []string{"-s", "--signal"} {
				signal := signal
				// With PID=1, `sleep infinity` will only exit when receiving SIGKILL. Default signal for kill is SIGKILL.
				// https://stackoverflow.com/questions/45148381/why-cant-i-ctrl-c-a-sleep-infinity-in-docker-when-it-runs-as-pid-1
				for _, term := range []string{"SIGTERM", "TERM"} {
					term := term
					ginkgo.It(fmt.Sprintf("should not kill the running container with %s %s", signal, term), func() {
						containerShouldBeRunning(o, testContainerName)
						command.Run(o, "kill", signal, term, testContainerName)
						containerShouldBeRunning(o, testContainerName)
					})
				}
			}
		})

		ginkgo.It("should fail to send the signal if the container doesn't exist", func() {
			command.RunWithoutSuccessfulExit(o, "kill", nonexistentContainerName)
		})
	})
}
