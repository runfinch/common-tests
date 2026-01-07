// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/option"
	"github.com/runfinch/common-tests/testutil"
)

type testCase struct {
	name string
	// extraArgs are the extra cli flags that will be passed while running the container
	extraArgs []string
	// noHealthCheck if true, will ignore all other healthchecks flags
	noHealthCheck    bool
	healthCheckFlags *healthCheckFlags
	// delaySec is the amount of time to wait after a container is running to ensure certain number
	// of healthchecks have run.
	delaySec       int
	shouldErr      bool
	errMsg         string
	expectedStatus *healthCheckStatus
	expectedFlags  *healthCheckFlags
}

type healthCheckFlags struct {
	cmd         string
	interval    int
	timeout     int
	retries     int
	startPeriod int
}

func (hcf *healthCheckFlags) toArgs() []string {
	args := []string{}
	if hcf.cmd != "" {
		args = append(args, "--health-cmd", hcf.cmd)
	}
	if hcf.interval != 0 {
		args = append(args, "--health-interval", fmt.Sprintf("%ds", hcf.interval))
	}
	if hcf.timeout != 0 {
		args = append(args, "--health-timeout", fmt.Sprintf("%ds", hcf.timeout))
	}
	if hcf.retries != 0 {
		args = append(args, "--health-retries", fmt.Sprintf("%d", hcf.retries))
	}
	if hcf.startPeriod != 0 {
		args = append(args, "--health-start-period", fmt.Sprintf("%ds", hcf.startPeriod))
	}
	return args
}

type healthCheckStatus struct {
	status         string
	failingStreak  int
	logContainsStr string
	maxLogEntries  int
	minLogLen      int
}

// HealthCheck tests both CLI healthcheck flags, manual and automated healthcheck features.
func HealthCheck(o *option.Option) {
	testCliHealthCheckFlags(o)
	testImageHealthCheckFlags(o)
	testHealthCheckStatus(o)
	testHealthCheckStatusNegativeCases(o)
}

func testCliHealthCheckFlags(o *option.Option) {
	ginkgo.Describe("test container healthcheck flags", func() {
		ginkgo.BeforeEach(func() {
			testutil.RequireNerdctlVersion(o, ">= 2.2.1")
		})
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		testCases := []testCase{
			{
				name: "should run a container with valid health check flags",
				healthCheckFlags: &healthCheckFlags{
					cmd:         "curl -f http://localhost || exit 1",
					interval:    30,
					timeout:     50,
					retries:     3,
					startPeriod: 20,
				},
				expectedFlags: &healthCheckFlags{
					cmd:         "[CMD-SHELL curl -f http://localhost || exit 1]",
					interval:    30,
					timeout:     50,
					retries:     3,
					startPeriod: 20,
				},
			},
			{
				name:             "should run a container without any healthcheck flags",
				healthCheckFlags: nil,
			},
			{
				name:          "should throw error for conflicting healthcheck flags",
				noHealthCheck: true,
				extraArgs:     []string{"--health-cmd", "true"},
				shouldErr:     true,
				errMsg:        "--no-healthcheck conflicts with --health-* options",
			},
			{
				name: "should throw error for negative --health-retries flag",
				healthCheckFlags: &healthCheckFlags{
					cmd:     "true",
					retries: -2,
				},
				shouldErr: true,
				errMsg:    "--health-retries cannot be negative",
			},
			{
				name: "should throw error for negative --health-timeout flag",
				healthCheckFlags: &healthCheckFlags{
					cmd:     "true",
					timeout: -5,
				},
				shouldErr: true,
				errMsg:    "--health-timeout cannot be negative",
			},
		}

		for _, tc := range testCases {
			ginkgo.It(tc.name, func() {
				var healthCheckArgs []string
				if tc.noHealthCheck {
					healthCheckArgs = []string{"--no-healthcheck"}
				} else if tc.healthCheckFlags != nil {
					healthCheckArgs = tc.healthCheckFlags.toArgs()
				}
				args := append([]string{"run", "-d", "--name", testContainerName}, tc.extraArgs...)
				args = append(args, healthCheckArgs...)
				args = append(args, localImages[defaultImage], "sleep", "infinity")
				if tc.shouldErr {
					stdErr := command.RunWithoutSuccessfulExit(o, args...).Err.Contents()
					gomega.Expect(string(stdErr)).Should(gomega.ContainSubstring(tc.errMsg))
					return
				}
				command.Run(o, args...)
				validateInspectHealthCheckFlags(o, tc.expectedFlags)
			})
		}
	})
}

func testImageHealthCheckFlags(o *option.Option) {
	ginkgo.Describe("test image healthcheck flags", func() {
		var buildContext string

		ginkgo.BeforeEach(func() {
			testutil.RequireNerdctlVersion(o, ">= 2.2.1")
		})
		ginkgo.BeforeEach(func() {
			buildContext = ffs.CreateBuildContext(fmt.Sprintf(`FROM %s
			HEALTHCHECK --interval=30s --timeout=10s CMD wget -q http://localhost:8080 || exit 1
			`, localImages[defaultImage]))

			ginkgo.DeferCleanup(os.RemoveAll, buildContext)
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		testCases := []testCase{
			{
				name:             "should run a container with healthcheck from image",
				healthCheckFlags: nil,
				expectedFlags: &healthCheckFlags{
					cmd:      "[CMD-SHELL wget -q http://localhost:8080 || exit 1]", // from dockerfile
					interval: 30,                                                    // from dockerfile
					timeout:  10,                                                    // from dockerfile
				},
			},
			{
				name: "should merge healthcheck flags from image with CLI flags",
				healthCheckFlags: &healthCheckFlags{
					retries:     3,
					startPeriod: 20,
				},
				expectedFlags: &healthCheckFlags{
					cmd:         "[CMD-SHELL wget -q http://localhost:8080 || exit 1]", // from dockerfile
					interval:    30,                                                    // from dockerfile
					timeout:     10,                                                    // from dockerfile
					retries:     3,                                                     // from cli flags
					startPeriod: 20,                                                    // from cli flags
				},
			},
			{
				name:          "should disable image healthcheck via CLI flag",
				noHealthCheck: true,
				expectedFlags: &healthCheckFlags{
					cmd: "[NONE]",
				},
			},
		}

		for _, tc := range testCases {
			ginkgo.It(tc.name, func() {
				var healthCheckArgs []string
				if tc.noHealthCheck {
					healthCheckArgs = []string{"--no-healthcheck"}
				} else if tc.healthCheckFlags != nil {
					healthCheckArgs = tc.healthCheckFlags.toArgs()
				}
				command.Run(o, "build", "-t", testImageName, buildContext)
				imageShouldExist(o, testImageName)
				args := append([]string{"run", "-d", "--name", testContainerName}, tc.extraArgs...)
				args = append(args, healthCheckArgs...)
				args = append(args, testImageName)
				command.Run(o, args...)
				validateInspectHealthCheckFlags(o, tc.expectedFlags)
			})
		}
	})
}

func testHealthCheckStatus(o *option.Option) {
	ginkgo.Describe("test automated container healthcheck status", func() {
		ginkgo.BeforeEach(func() {
			testutil.RequireNerdctlVersion(o, ">= 2.2.1")
		})
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		testCases := []testCase{
			{
				name: "should report container as unhealthy when health-cmd fails retryCount times",
				healthCheckFlags: &healthCheckFlags{
					cmd:      "exit 1",
					retries:  3,
					interval: 1,
				},
				// every health-cmd run exits immediately and runs every 1 sec. So wait for 3 seconds to make sure
				// healthcheck has run 3 (retry count) times.
				delaySec: 3,
				expectedStatus: &healthCheckStatus{
					status:        "unhealthy",
					failingStreak: 3,
				},
			},
			{
				name: "should report container as unhealthy when health-cmd times out retryCount times",
				healthCheckFlags: &healthCheckFlags{
					cmd:      "sleep 10",
					retries:  3,
					interval: 1,
					timeout:  1,
				},
				// every health-cmd run times out after 1 second and runs every 1 sec. So wait for 6 seconds to make sure
				// healthcheck has run 3 (retryCount) times.
				delaySec: 6,
				expectedStatus: &healthCheckStatus{
					status:        "unhealthy",
					failingStreak: 3,
				},
			},
			{
				name: "should report container as unhealthy if an invalid health-cmd is specified",
				healthCheckFlags: &healthCheckFlags{
					cmd:      "invalid command",
					retries:  3,
					interval: 1,
				},
				delaySec: 3,
				expectedStatus: &healthCheckStatus{
					status:        "unhealthy",
					failingStreak: 3,
				},
			},
			{
				name: "should not count health-cmd failures if a start-period is specified",
				healthCheckFlags: &healthCheckFlags{
					cmd:         "exit 1",
					retries:     3,
					interval:    1,
					startPeriod: 60,
				},
				delaySec: 3,
				expectedStatus: &healthCheckStatus{
					status:        "starting",
					failingStreak: 0,
				},
			},
			{
				name: "should report container as healthy if health-cmd succeeds during start-period",
				healthCheckFlags: &healthCheckFlags{
					cmd:         "echo hello",
					retries:     3,
					interval:    1,
					startPeriod: 60,
				},
				expectedStatus: &healthCheckStatus{
					status:        "healthy",
					failingStreak: 0,
				},
			},
			{
				name: "should report container as healthy",
				healthCheckFlags: &healthCheckFlags{
					cmd:      "echo hello",
					interval: 30,
				},
				expectedStatus: &healthCheckStatus{
					status: "healthy",
				},
			},
			{
				name:           "should disable healthcheck --no-healthcheck flag is specified",
				noHealthCheck:  true,
				expectedStatus: nil,
			},
			{
				name:      "should work with container environment variables",
				extraArgs: []string{"-e", "ENV_VAR=hello"},
				healthCheckFlags: &healthCheckFlags{
					cmd:      "echo $ENV_VAR",
					interval: 30,
				},
				expectedStatus: &healthCheckStatus{
					status:         "healthy",
					logContainsStr: "hello",
				},
			},
			{
				name:      "should work with container workdir",
				extraArgs: []string{"--workdir", "/tmp"},
				healthCheckFlags: &healthCheckFlags{
					cmd:      "pwd",
					interval: 30,
				},
				expectedStatus: &healthCheckStatus{
					status:         "healthy",
					logContainsStr: "/tmp",
				},
			},
			{
				name: "should log large healthcheck output",
				healthCheckFlags: &healthCheckFlags{
					cmd:      "cat /dev/urandom | tr -dc A-Za-z0-9 | head -c 50000",
					interval: 1,
					timeout:  2,
				},
				expectedStatus: &healthCheckStatus{
					status:    "healthy",
					minLogLen: 1024,
				},
			},
			{
				name: "should truncate large log healthcheck output",
				healthCheckFlags: &healthCheckFlags{
					cmd:      "cat /dev/urandom | tr -dc A-Za-z0-9 | head -c 100000",
					interval: 1,
					timeout:  2,
				},
				expectedStatus: &healthCheckStatus{
					status:         "healthy",
					logContainsStr: "truncated",
					minLogLen:      1024,
				},
			},
			{
				name: "should contain at most 5 healthcheck log entries",
				healthCheckFlags: &healthCheckFlags{
					cmd:      "cat /dev/urandom | tr -dc A-Za-z0-9 | head -c 16",
					interval: 1,
					timeout:  1,
				},
				delaySec: 10,
				expectedStatus: &healthCheckStatus{
					status:        "healthy",
					maxLogEntries: 5,
				},
			},
		}

		for _, tc := range testCases {
			ginkgo.It(tc.name, func() {
				var healthCheckArgs []string
				if tc.noHealthCheck {
					healthCheckArgs = []string{"--no-healthcheck"}
				} else if tc.healthCheckFlags != nil {
					healthCheckArgs = tc.healthCheckFlags.toArgs()
				}
				args := append([]string{"run", "-d", "--name", testContainerName}, tc.extraArgs...)
				args = append(args, healthCheckArgs...)
				args = append(args, localImages[defaultImage], "sleep", "infinity")
				command.Run(o, args...)
				waitTillContainerStatus(o, "running")
				if tc.delaySec > 0 {
					time.Sleep(time.Duration(tc.delaySec) * time.Second)
				}
				validateInspectHealthCheckStatus(o, tc.expectedStatus)
			})
		}
	})
}

func testHealthCheckStatusNegativeCases(o *option.Option) {
	ginkgo.Describe("test manual container healthcheck status", func() {
		ginkgo.BeforeEach(func() {
			testutil.RequireNerdctlVersion(o, ">= 2.2.1")
		})
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should fail healthcheck on non-existent container", func() {
			stdErr := command.RunWithoutSuccessfulExit(o, "container", "healthcheck", "nonexistent").Err.Contents()
			gomega.Expect(string(stdErr)).Should(gomega.ContainSubstring("no such container nonexistent"))
		})

		ginkgo.It("should fail healthcheck on missing healthcheck config", func() {
			command.Run(o, "run", "-d", "--name", testContainerName, localImages[defaultImage], "sleep", "infinity")
			waitTillContainerStatus(o, "running")
			stdErr := command.RunWithoutSuccessfulExit(o, "container", "healthcheck", testContainerName).Err.Contents()
			gomega.Expect(string(stdErr)).Should(gomega.ContainSubstring("container has no health check configured"))
		})

		ginkgo.It("should fail healthcheck for paused container", func() {
			healthCheckArgs := []string{
				"--health-cmd", "echo healthy",
				"--health-interval", "30s",
			}
			args := append([]string{"run", "-d", "--name", testContainerName}, healthCheckArgs...)
			args = append(args, localImages[defaultImage], "sleep", "infinity")
			command.Run(o, args...)
			command.Run(o, "pause", testContainerName)
			waitTillContainerStatus(o, "paused")
			stdErr := command.RunWithoutSuccessfulExit(o, "container", "healthcheck", testContainerName).Err.Contents()
			gomega.Expect(string(stdErr)).Should(gomega.ContainSubstring("container is not running (status: paused)"))
		})

		ginkgo.It("should fail healthcheck for stopped container", func() {
			healthCheckArgs := []string{
				"--health-cmd", "echo healthy",
				"--health-interval", "30s",
			}
			args := append([]string{"run", "-d", "--name", testContainerName}, healthCheckArgs...)
			args = append(args, localImages[defaultImage], "sleep", "2s")
			command.Run(o, args...)
			waitTillContainerStatus(o, "exited")
			stdErr := command.RunWithoutSuccessfulExit(o, "container", "healthcheck", testContainerName).Err.Contents()
			gomega.Expect(string(stdErr)).Should(gomega.ContainSubstring("container is not running (status: stopped)"))
		})

		ginkgo.It("should fail healthcheck for nonexistent container task", func() {
			healthCheckArgs := []string{
				"--health-cmd", "echo healthy",
				"--health-interval", "30s",
			}
			args := append([]string{"create", "--name", testContainerName}, healthCheckArgs...)
			args = append(args, localImages[defaultImage], "sleep", "infinity")
			command.Run(o, args...)
			stdErr := command.RunWithoutSuccessfulExit(o, "container", "healthcheck", testContainerName).Err.Contents()
			gomega.Expect(string(stdErr)).Should(gomega.ContainSubstring("failed to get container task: no running task found"))
		})
	})
}

func validateInspectHealthCheckFlags(o *option.Option, expect *healthCheckFlags) {
	if expect == nil {
		inspectHealthCheck := command.StdoutStr(o, "inspect", "--format", "{{.Config.Healthcheck}}", testContainerName)
		gomega.Expect(inspectHealthCheck).Should(gomega.Equal("<nil>"))
		return
	}
	if expect.cmd != "" {
		inspectTest := command.StdoutStr(o, "inspect", "--format", "{{.Config.Healthcheck.Test}}", testContainerName)
		gomega.Expect(inspectTest).Should(gomega.Equal(expect.cmd))
	}
	if expect.interval != 0 {
		inspectHealthInterval := command.StdoutStr(o, "inspect", "--format", "{{.Config.Healthcheck.Interval}}", testContainerName)
		gomega.Expect(inspectHealthInterval).Should(gomega.Equal(fmt.Sprintf("%ds", expect.interval)))
	}
	if expect.timeout != 0 {
		inspectTimeout := command.StdoutStr(o, "inspect", "--format", "{{.Config.Healthcheck.Timeout}}", testContainerName)
		gomega.Expect(inspectTimeout).Should(gomega.Equal(fmt.Sprintf("%ds", expect.timeout)))
	}
	if expect.retries != 0 {
		inspectRetries := command.StdoutStr(o, "inspect", "--format", "{{.Config.Healthcheck.Retries}}", testContainerName)
		gomega.Expect(inspectRetries).Should(gomega.Equal(fmt.Sprintf("%d", expect.retries)))
	}
	if expect.startPeriod != 0 {
		inspectStartPeriod := command.StdoutStr(o, "inspect", "--format", "{{.Config.Healthcheck.StartPeriod}}", testContainerName)
		gomega.Expect(inspectStartPeriod).Should(gomega.Equal(fmt.Sprintf("%ds", expect.startPeriod)))
	}
}

func validateInspectHealthCheckStatus(o *option.Option, expect *healthCheckStatus) {
	if expect == nil {
		inspectHealthStatus := command.StdoutStr(o, "inspect", "--format", "{{.State.Health}}", testContainerName)
		gomega.Expect(inspectHealthStatus).Should(gomega.Equal("<nil>"))
		return
	}
	if expect.status != "" {
		inspectHealthStatus := command.StdoutStr(o, "inspect", "--format", "{{.State.Health.Status}}", testContainerName)
		gomega.Expect(inspectHealthStatus).Should(gomega.Equal(expect.status))
	}
	inspectFailingStreak, err := strconv.Atoi(command.StdoutStr(o, "inspect", "--format", "{{.State.Health.FailingStreak}}", testContainerName))
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	if expect.failingStreak > 0 {
		gomega.Expect(inspectFailingStreak).Should(gomega.BeNumerically(">=", expect.failingStreak))
	} else {
		gomega.Expect(inspectFailingStreak).Should(gomega.BeNumerically("==", expect.failingStreak))
	}
	if expect.maxLogEntries > 0 {
		inspectLogEntries, err := strconv.Atoi(command.StdoutStr(o, "inspect", "--format", "{{len .State.Health.Log}}", testContainerName))
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(inspectLogEntries).Should(gomega.BeNumerically("<=", expect.maxLogEntries))
	}
	if expect.logContainsStr != "" {
		inspectLog := command.StdoutStr(o, "inspect", "--format", "{{(index .State.Health.Log 0).Output}}", testContainerName)
		gomega.Expect(inspectLog).Should(gomega.ContainSubstring(expect.logContainsStr))
	}
	if expect.minLogLen > 0 {
		inspectLogLen := len(command.StdoutStr(o, "inspect", "--format", "{{(index .State.Health.Log 0).Output}}", testContainerName))
		gomega.Expect(inspectLogLen).Should(gomega.BeNumerically(">=", expect.minLogLen))
	}
}

func waitTillContainerStatus(o *option.Option, status string) {
	retryCount := 20
	retryInterval := time.Second
	for i := 1; i <= retryCount; i++ {
		time.Sleep(retryInterval)
		if command.StdoutStr(o, "inspect", "--format", "{{.State.Status}}", testContainerName) == status {
			return
		}
	}
	ginkgo.Fail(fmt.Sprintf("container is still not in status \"%s\" after %d attempts", status, retryCount))
}
