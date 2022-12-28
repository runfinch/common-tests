// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"os"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Ps tests functionality of `ps` command.
func Ps(o *option.Option) {
	sha256RegexTruncated := `^[a-f0-9]{12}$`
	sha256RegexFull := `^[a-f0-9]{64}$`

	containerNames := []string{"ctr_1", "ctr_2"}
	ginkgo.Describe("Ps command", func() {
		pwd, _ := os.Getwd()
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			command.Run(o, "network", "create", testNetwork)
			command.Run(o, "run", "-d",
				"--name", containerNames[0],
				"--label", "color=red",
				"-v", fmt.Sprintf("%s:%s", pwd, pwd),
				"-w", pwd,
				defaultImage)
			command.Run(o, "run", "-d",
				"--label", "color=green",
				"--network", testNetwork,
				"-p", "8081:80",
				"--name", containerNames[1],
				defaultImage, "sleep", "infinity")
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.It("should list only running containers", func() {
			psOutput := command.StdOutAsLines(o, "ps", "--format", "{{.Names}}")
			gomega.Expect(psOutput).ShouldNot(gomega.ContainElement(containerNames[0]))
			gomega.Expect(psOutput).Should(gomega.ContainElement(containerNames[1]))
		})

		for _, flag := range []string{"-a", "--all"} {
			flag := flag
			ginkgo.It(fmt.Sprintf("should list all containers [%s flag]", flag), func() {
				psOutput := command.StdOutAsLines(o, "ps", "--format", "{{.Names}}", flag)
				gomega.Expect(psOutput).Should(gomega.ContainElements(containerNames))
			})
		}

		ginkgo.It("should list ID of the containers", func() {
			psOutput := command.StdOutAsLines(o, "ps", "--format", "{{.ID}}")
			gomega.Expect(psOutput).Should(gomega.HaveEach(gomega.MatchRegexp(sha256RegexTruncated)))
		})
		ginkgo.It("should list image of the containers", func() {
			psOutput := command.StdOutAsLines(o, "ps", "--format", "{{.Image}}")
			gomega.Expect(psOutput).Should(gomega.ContainElement(defaultImage))
		})
		ginkgo.It("should list command of the containers", func() {
			psOutput := command.StdoutStr(o, "ps", "--format", "{{.Command}}")
			gomega.Expect(psOutput).Should(gomega.ContainSubstring("sleep infinity"))
		})
		ginkgo.It("should list creation date of the containers", func() {
			psOutput := command.StdOutAsLines(o, "ps", "--format", "{{.CreatedAt}}")
			gomega.Expect(psOutput).ShouldNot(gomega.ContainElement(gomega.BeEmpty()))
		})
		ginkgo.It("should list port forwarding info of the containers", func() {
			psOutput := command.StdoutStr(o, "ps", "--format", "{{.Ports}}")
			gomega.Expect(psOutput).Should(gomega.ContainSubstring("8081"))
		})
		ginkgo.It("should list only running containers", func() {
			psOutput := command.StdOutAsLines(o, "ps", "--format", "{{.Status}}")
			gomega.Expect(psOutput).Should(gomega.ContainElement("Up"))
		})

		for _, flag := range []string{"-q", "--quiet"} {
			flag := flag
			ginkgo.It(fmt.Sprintf("should list truncated container IDs [%s flag]", flag), func() {
				psOutput := command.StdOutAsLines(o, "ps", flag)
				gomega.Expect(psOutput).ShouldNot(gomega.BeEmpty())
				gomega.Expect(psOutput).Should(gomega.HaveEach(gomega.MatchRegexp(sha256RegexTruncated)))
			})
		}

		ginkgo.It("should list full container IDs", func() {
			psOutput := command.StdOutAsLines(o, "ps", "--format", "{{.ID}}", "--no-trunc")
			gomega.Expect(psOutput).ShouldNot(gomega.BeEmpty())
			gomega.Expect(psOutput).Should(gomega.HaveEach(gomega.MatchRegexp(sha256RegexFull)))
		})

		for _, flag := range []string{"-s", "--size"} {
			flag := flag
			ginkgo.It(fmt.Sprintf("should list container size [%s flag]", flag), func() {
				psOutput := command.StdoutStr(o, "ps", "--format", "{{.Size}}", flag)
				gomega.Expect(psOutput).ShouldNot(gomega.BeEmpty())
			})
		}

		for _, flag := range []string{"-n", "--last"} {
			flag := flag
			ginkgo.It(fmt.Sprintf("should list last 1 containers [%s flag]", flag), func() {
				psOutput := command.StdOutAsLines(o, "ps", "--format", "{{.Names}}", flag, "1")
				gomega.Expect(psOutput).ShouldNot(gomega.ContainElement(containerNames[0]))
				gomega.Expect(psOutput).Should(gomega.ContainElement(containerNames[1]))
			})
		}

		for _, flag := range []string{"-l", "--latest"} {
			flag := flag
			ginkgo.It(fmt.Sprintf("should list last 1 containers [%s flag]", flag), func() {
				psOutput := command.StdOutAsLines(o, "ps", "--format", "{{.Names}}", flag)
				gomega.Expect(psOutput).ShouldNot(gomega.ContainElement(containerNames[0]))
				gomega.Expect(psOutput).Should(gomega.ContainElement(containerNames[1]))
			})
		}

		ginkgo.Context("should list container ", func() {
			filterTests := []struct {
				filter         string
				expectedOutput []string
			}{
				{
					filter:         fmt.Sprintf("name=%s", containerNames[0]),
					expectedOutput: []string{containerNames[0]},
				},
				{
					filter:         "label=color=green",
					expectedOutput: []string{containerNames[1]},
				},
				{
					filter:         "label=color",
					expectedOutput: containerNames,
				},
				{
					filter:         "exited=0",
					expectedOutput: []string{containerNames[0]},
				},
				{
					filter:         "status=exited",
					expectedOutput: []string{containerNames[0]},
				},
				{
					filter:         "status=running",
					expectedOutput: []string{containerNames[1]},
				},
				{
					filter:         fmt.Sprintf("before=%s", containerNames[1]),
					expectedOutput: []string{containerNames[0]},
				},
				{
					filter:         fmt.Sprintf("since=%s", containerNames[0]),
					expectedOutput: []string{containerNames[1]},
				},
				{
					filter:         fmt.Sprintf("volume=%s", pwd),
					expectedOutput: []string{containerNames[0]},
				},
				{
					filter:         fmt.Sprintf("network=%s", testNetwork),
					expectedOutput: []string{containerNames[1]},
				},
			}

			for _, test := range filterTests {
				test := test
				ginkgo.It(fmt.Sprintf("with filter %s", test.filter), func() {
					output := command.StdOutAsLines(o, "ps", "-a", "--format", "{{.Names}}", "--filter", test.filter)
					gomega.Expect(output).Should(gomega.ContainElements(test.expectedOutput))
				})
			}
		})
	})
}
