// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/option"
)

// ComposeBuild tests functionality of `compose build` command.
func ComposeBuild(o *option.Option) {
	services := []string{"svc1_build_cmd", "svc2_build_cmd"}
	imageSuffix := []string{"alpine:latest", "-svc2_build_cmd:latest"}
	ginkgo.Describe("Compose build command", func() {
		var composeContext string
		var composeFilePath string
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			composeContext, composeFilePath = createComposeYmlForBuildCmd(services)
			ginkgo.DeferCleanup(os.RemoveAll, composeContext)
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should build services defined in the compose file", func() {
			command.Run(o, "compose", "build", "--file", composeFilePath)

			imageList := command.GetAllImageNames(o)
			gomega.Expect(imageList).Should(gomega.ContainElement(gomega.HaveSuffix(imageSuffix[0])))
			gomega.Expect(imageList).Should(gomega.ContainElement(gomega.HaveSuffix(imageSuffix[1])))
			// The built image should print 'Compose build test' when run.
			output := command.StdoutStr(o, "run", localImages[defaultImage])
			gomega.Expect(output).Should(gomega.Equal("Compose build test"))
		})

		ginkgo.It("should build services defined in the compose file specified by the COMPOSE_FILE environment variable", func() {
			envKey := "COMPOSE_FILE"
			o.UpdateEnv(envKey, composeFilePath)

			command.Run(o, "compose", "build")

			imageList := command.GetAllImageNames(o)
			gomega.Expect(imageList).Should(gomega.ContainElement(gomega.HaveSuffix(imageSuffix[0])))
			gomega.Expect(imageList).Should(gomega.ContainElement(gomega.HaveSuffix(imageSuffix[1])))
			// The built image should print 'Compose build test' when run.
			output := command.StdoutStr(o, "run", localImages[defaultImage])
			gomega.Expect(output).Should(gomega.Equal("Compose build test"))

			o.DeleteEnv(envKey)
		})

		ginkgo.It("should output progress in plain text format", func() {
			composeBuildOutput := command.StderrStr(o, "compose", "build", "--progress",
				"plain", "--no-cache", "--file", composeFilePath)
			// The docker file contains following command.
			// RUN printf 'should only see the final answer when "--progress" is set to be "plain": %d\n' $(expr 1 + 1)
			// where the expression "$(expr 1 + 1)" will be evaluated to 2 only for "--progress plain" output.
			gomega.Expect(composeBuildOutput).Should(gomega.ContainSubstring(
				`should only see the final answer when '--progress' is set to be 'plain': 2`))

			imageList := command.GetAllImageNames(o)
			gomega.Expect(imageList).Should(gomega.ContainElement(gomega.HaveSuffix(imageSuffix[0])))
			gomega.Expect(imageList).Should(gomega.ContainElement(gomega.HaveSuffix(imageSuffix[1])))
		})

		ginkgo.It("should build services defined in the compose file with --build-args", func() {
			command.Run(o, "compose", "build", "--build-arg",
				`CMD_MSG=Compose build with --build-arg`, "--file", composeFilePath)
			output := command.StdoutStr(o, "compose", "up", "--file", composeFilePath)
			gomega.Expect(output).Should(gomega.ContainSubstring("Compose build with --build-arg"))
			command.Run(o, "compose", "down", "--file", composeFilePath)
		})
		ginkgo.It("should build services defined in the compose file without --build-args", func() {
			command.Run(o, "compose", "build", "--file", composeFilePath)

			// The built image should print default value of the build-arg which is 'Compose build test'.
			output := command.StdoutStr(o, "compose", "up", "--file", composeFilePath)
			gomega.Expect(output).Should(gomega.ContainSubstring("Compose build test"))
			command.Run(o, "compose", "down", "--file", composeFilePath)
		})
		// TODO: --no-cache does not have any effect on the build output.
		ginkgo.It("should build services defined in the compose file without using cache", func() {
			command.Run(o, "compose", "build", "--no-cache", "--file", composeFilePath)
			imageList := command.GetAllImageNames(o)
			gomega.Expect(imageList).Should(gomega.ContainElement(gomega.HaveSuffix(imageSuffix[0])))
			gomega.Expect(imageList).Should(gomega.ContainElement(gomega.HaveSuffix(imageSuffix[1])))
		})
		// TODO: add functional test for --ipfs
	})
}

func createComposeYmlForBuildCmd(serviceNames []string) (string, string) {
	gomega.Expect(serviceNames).Should(gomega.HaveLen(2))

	dockerFileContent := fmt.Sprintf(`
FROM %s
ARG CMD_MSG="Compose build test"
RUN printf "should only see the final answer when '--progress' is set to be 'plain': %%d\n" $(expr 1 + 1) 
ENV ENV_CMD_MSG=${CMD_MSG}
CMD echo ${ENV_CMD_MSG}
`, localImages[defaultImage])

	composeYmlContent := fmt.Sprintf(
		`
services:
  %[1]s:
    build: 
      context: .
      dockerfile: Dockerfile
    image: %[3]s
  %[2]s:
    build: 
      context: .
      dockerfile: Dockerfile
`, serviceNames[0], serviceNames[1], localImages[defaultImage])

	composeDir, composeFilePath := ffs.CreateComposeYmlContext(composeYmlContent)
	ffs.WriteFile(filepath.Join(composeDir, "Dockerfile"), dockerFileContent)
	return composeDir, composeFilePath
}
