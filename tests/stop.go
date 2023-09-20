// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Stop tests stopping a container.
func Stop(o *option.Option) {
	ginkgo.Describe("stop a container", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should stop the container if the container is running", func() {
			command.Run(o, "run", "-d", "--name", testContainerName, defaultImage, "nc", "-l")
			containerShouldBeRunning(o, testContainerName)

			command.Run(o, "stop", testContainerName)
			command.RunWithoutSuccessfulExit(o, "exec", testContainerName, "echo", "foo")
		})

		for _, timeFlag := range []string{"-t", "--time"} {
			timeFlag := timeFlag
			ginkgo.It(fmt.Sprintf("should stop running container within specified time by %s flag", timeFlag), func() {
				// With PID=1, `sleep infinity` does not exit due to receiving a SIGTERM, which is sent by the stop command.
				// Ref. https://superuser.com/a/1299463/730265
				command.Run(o, "run", "-d", "--name", testContainerName, defaultImage, "sleep", "infinity")
				gomega.Expect(command.StdoutStr(o, "exec", testContainerName, "echo", "foo")).To(gomega.Equal("foo"))
				startTime := time.Now()
				command.Run(o, "stop", "-t", "1", testContainerName)
				gomega.Expect(time.Since(startTime)).To(gomega.BeNumerically("~", 1*time.Second, 750*time.Millisecond))
				command.RunWithoutSuccessfulExit(o, "exec", testContainerName, "echo", "foo")
			})
		}
	})
}
