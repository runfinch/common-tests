// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/option"
)

// Build command building an image.
//
// TODO:  --no-cache, --output, --ssh
// --no-cache flag is added to tests asserting the output from `RUN` command.
// [Discussion]: https://github.com/runfinch/common-tests/pull/4#discussion_r971338825
func Build(o *option.Option) {
	ginkgo.Describe("Build container image", func() {
		ginkgo.Context("Build container image using default image", func() {
			var buildContext string
			ginkgo.BeforeEach(func() {
				buildContext = ffs.CreateBuildContext(fmt.Sprintf(`FROM %s
			CMD ["echo", "finch-test-dummy-output"]
			`, defaultImage))
				ginkgo.DeferCleanup(os.RemoveAll, buildContext)
				command.RemoveAll(o)
			})

			ginkgo.AfterEach(func() {
				command.RemoveAll(o)
			})

			for _, tag := range []string{"-t", "--tag"} {
				tag := tag
				ginkgo.It(fmt.Sprintf("build basic alpine image with %s option", tag), func() {
					command.Run(o, "build", tag, testImageName, buildContext)
					imageShouldExist(o, testImageName)
				})
			}

			ginkgo.Context("Build an image with file option", func() {
				var dockerFilePath string
				ginkgo.BeforeEach(func() {
					dockerFilePath = filepath.Join(buildContext, "AnotherDockerfile")
					ffs.WriteFile(dockerFilePath, fmt.Sprintf(`FROM %s
			RUN ["echo", "built from AnotherDockerfile"]
			`, defaultImage))
				})

				for _, file := range []string{"-f", "--file"} {
					file := file
					ginkgo.It(fmt.Sprintf("build an image with %s option", file), func() {
						stdErr := command.StdErr(o, "build", "--no-cache", file, dockerFilePath, buildContext)
						gomega.Expect(stdErr).Should(gomega.ContainSubstring("built from AnotherDockerfile"))
					})
				}
			})

			ginkgo.It("build image with --secret option", func() {
				containerWithSecret := fmt.Sprintf(`FROM %s
			RUN --mount=type=secret,id=mysecret cat /run/secrets/mysecret
			`, defaultImage)
				dockerFilePath := filepath.Join(buildContext, "Dockerfile.with-secret")
				ffs.WriteFile(dockerFilePath, containerWithSecret)
				secretFile := filepath.Join(buildContext, "secret.txt")
				ffs.WriteFile(secretFile, "somesecret")
				secret := fmt.Sprintf("id=mysecret,src=%s", secretFile)
				stdErr := command.StdErr(o, "build", "--progress=plain", "--no-cache", "-f", dockerFilePath, "--secret", secret, buildContext)
				gomega.Expect(stdErr).Should(gomega.ContainSubstring("somesecret"))
			})

			ginkgo.It("build image with --target option", func() {
				containerWithTarget := fmt.Sprintf(`FROM %s AS build_env
			RUN echo output from build_env
			FROM %s AS prod_env
			RUN  echo "output from prod_env
			`, defaultImage, defaultImage)
				dockerFilePath := filepath.Join(buildContext, "Dockerfile.with-target")
				ffs.WriteFile(dockerFilePath, containerWithTarget)
				stdEr := command.StdErr(o, "build", "--progress=plain", "--no-cache",
					"-f", dockerFilePath, "--target", "build_env", buildContext)
				gomega.Expect(stdEr).Should(gomega.ContainSubstring("output from build_env"))
				gomega.Expect(stdEr).ShouldNot(gomega.ContainSubstring("output from prod_env"))
			})

			// "--output=type=docker" is intentional for the imageId to show up
			ginkgo.It("build image with --quiet option", func() {
				commandOut := command.StdoutStr(o, "build", "--output=type=docker", "--quiet", buildContext)
				gomega.Expect(len(strings.Split(commandOut, "\n"))).To(gomega.Equal(1))
			})

			ginkgo.It("build image with --build-arg option", func() {
				containerWithBuildArg := "ARG VERSION=latest \n FROM public.ecr.aws/docker/library/alpine:${VERSION}"
				dockerFilePath := filepath.Join(buildContext, "Dockerfile.with-build-arg")
				ffs.WriteFile(dockerFilePath, containerWithBuildArg)
				stdErr := command.StdErr(o, "build", "-f", dockerFilePath, "--no-cache", "--progress=plain",
					"--build-arg", "VERSION=3.13", buildContext)
				gomega.Expect(stdErr).Should(gomega.ContainSubstring("public.ecr.aws/docker/library/alpine:3.13"))
			})

			ginkgo.It("build image with --progress=plain", func() {
				dockerFile := fmt.Sprintf(`FROM %s
				RUN echo "progress flag set:$((1 + 1))"
			`, defaultImage)
				dockerFilePath := filepath.Join(buildContext, "Dockerfile.progress")
				ffs.WriteFile(dockerFilePath, dockerFile)
				stdErr := command.StdErr(o, "build", "-f", dockerFilePath, "--no-cache", "--progress=plain", buildContext)
				gomega.Expect(stdErr).Should(gomega.ContainSubstring("progress flag set:2"))
			})
		})

		ginkgo.Context("Build container image using alpine image", func() {
			var buildContext string
			ginkgo.BeforeEach(func() {
				buildContext = ffs.CreateBuildContext(fmt.Sprintf(`FROM %s
			CMD ["echo", "finch-test-dummy-output"]
			`, alpineImage))
				ginkgo.DeferCleanup(os.RemoveAll, buildContext)
				command.RemoveAll(o)
			})

			ginkgo.AfterEach(func() {
				command.RemoveAll(o)
			})
			// If SetupLocalRegistry is invoked before this test case,
			// then defaultImage will point to the image in the local registry,
			// and there will be only one platform (i.e., the platform of the running machine) available for that image in the local registry.
			// As a result, to make this test case not flaky even when SetupLocalRegistry is used,
			// we need to pull alpineImage instead of defaultImage
			// because we can be sure that the registry associated with the former provides the image with the platform specified below.
			ginkgo.It("build basic alpine image with --platform option", func() {
				command.Run(o, "build", "-t", testImageName, "--platform=amd64", buildContext)
				platform := command.StdoutStr(o, "images", testImageName, "--format", "{{.Platform}}")
				gomega.Expect(platform).Should(gomega.Equal("linux/amd64"))
			})
		})
	})
}
