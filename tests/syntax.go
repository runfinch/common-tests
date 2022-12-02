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

// Syntax Tests for instructions described in https://docs.docker.com/engine/reference/builder/
// Syntax TODO: #syntax directive has a bug see: https://github.com/moby/buildkit/issues/3138, ONBUILD, HEALTHCHECK,STOPSIGNAL.
func Syntax(o *option.Option) {
	ginkgo.Describe("Syntax tests", func() {
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

		ginkgo.It("should correctly parse escape directive", func() {
			fileName := "Dockerfile.EscapeDirective"
			instructions := "# escape=`\nFROM public.ecr.aws/docker/library/alpine:3.1`3"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			out := command.StdErr(o, "build", "-f", dockerFilePath, "--no-cache", "--progress=plain", buildContext)
			gomega.Expect(out).Should(gomega.ContainSubstring("public.ecr.aws/docker/library/alpine:3.13"))
		})

		ginkgo.It("should correctly parse ENV instruction", func() {
			fileName := "Dockerfile.EnvInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\nENV FOO=\"ENV test\"\nRUN echo ${FOO}_arg"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			out := command.StdErr(o, "build", "-f", dockerFilePath, "--no-cache", "--progress=plain", buildContext)
			gomega.Expect(out).Should(gomega.ContainSubstring("ENV test_arg"))
		})

		ginkgo.It("should correctly parse FROM instruction", func() {
			fileName := "Dockerfile.FromInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			out := command.StdErr(o, "build", "-f", dockerFilePath, "--no-cache", "--progress=plain", buildContext)
			gomega.Expect(out).Should(gomega.ContainSubstring("public.ecr.aws/docker/library/alpine:3.13"))
		})

		ginkgo.It("should correctly parse LABEL instruction", func() {
			fileName := "Dockerfile.LabelInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\nLABEL \"Maintainer\" = \"maintainer@maintainer.com\""
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			command.Run(o, "build", "-f", dockerFilePath, "-t", "label", buildContext)
			out := command.StdOut(o, "inspect", "label")
			gomega.Expect(out).Should(gomega.ContainSubstring("maintainer@maintainer.com"))
		})

		ginkgo.It("should correctly parse CMD instruction in shell form", func() {
			fileName := "Dockerfile.CmdInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\nCMD echo \"shell form\""
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			command.Run(o, "build", "-f", dockerFilePath, "-t", "cmd-shell", buildContext)
			out := command.StdOut(o, "run", "cmd-shell")
			gomega.Expect(out).Should(gomega.ContainSubstring("shell form"))
		})

		ginkgo.It("should correctly parse CMD instruction in exec form", func() {
			fileName := "Dockerfile.CmdInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n CMD [\"sh\", \"-c\", \"echo exec form\" ]"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			command.Run(o, "build", "-f", dockerFilePath, "-t", "cmd-exec", buildContext)
			out := command.StdOut(o, "run", "cmd-exec")
			gomega.Expect(out).Should(gomega.ContainSubstring("exec form"))
		})

		ginkgo.It("should correctly parse RUN instruction", func() {
			fileName := "Dockerfile.CmdInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n RUN [\"sh\", \"-c\", \"echo hello from run\" ]"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			out := command.StdErr(o, "build", "-f", dockerFilePath, "--no-cache", "--progress=plain", buildContext)
			gomega.Expect(out).Should(gomega.ContainSubstring("hello from run"))
		})

		ginkgo.It("should correctly parse EXPOSE instruction", func() {
			fileName := "Dockerfile.ExposeInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n EXPOSE 80/tcp"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			command.Run(o, "build", "-f", dockerFilePath, "-t", "expose", buildContext)
			out := command.StdOut(o, "inspect", "expose")
			gomega.Expect(out).Should(gomega.ContainSubstring("80/tcp"))
		})

		ginkgo.It("should correctly parse ADD instruction", func() {
			fileName := "Dockerfile.AddInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n ADD Dockerfile.AddInstruction .\n CMD ls"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			command.Run(o, "build", "-f", dockerFilePath, "-t", "add", buildContext)
			out := command.StdOut(o, "run", "add")
			gomega.Expect(out).Should(gomega.ContainSubstring("Dockerfile.AddInstruction"))
		})

		ginkgo.It("should correctly parse COPY instruction", func() {
			fileName := "Dockerfile.CopyInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n COPY Dockerfile.CopyInstruction .\n CMD ls"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			command.Run(o, "build", "-f", dockerFilePath, "-t", "copy", buildContext)
			out := command.StdOut(o, "run", "copy")
			gomega.Expect(out).Should(gomega.ContainSubstring("Dockerfile.CopyInstruction"))
		})

		ginkgo.It("should correctly parse ENTRYPOINT instruction", func() {
			fileName := "Dockerfile.EntrypointInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n ENTRYPOINT [\"echo\"," +
				" \"entrypoint\"]\n CMD echo cmd \" "
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			command.Run(o, "build", "-f", dockerFilePath, "-t", "entrypoint", buildContext)
			out := command.StdOut(o, "run", "entrypoint", "echo", "override cmd")
			gomega.Expect(out).Should(gomega.ContainSubstring("entrypoint echo override cmd"))
		})

		ginkgo.It("should correctly parse VOLUME instruction", func() {
			fileName := "Dockerfile.VolumeInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n RUN mkdir /myvol\n " +
				"RUN echo \"Greetings Volume\" > /myvol/greeting\n VOLUME /myvol\n CMD [\"cat\", \"/myvol/greeting\"]"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			command.Run(o, "build", "-f", dockerFilePath, "-t", "volume", buildContext)
			out := command.StdOut(o, "run", "volume")
			gomega.Expect(out).Should(gomega.ContainSubstring("Greetings Volume"))
		})

		ginkgo.It("should correctly parse USER instruction", func() {
			fileName := "Dockerfile.UserInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n RUN addgroup -S appgroup && adduser " +
				"-S someuser -G appgroup\n USER someuser\n CMD whoami"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			command.Run(o, "build", "-f", dockerFilePath, "-t", "user", buildContext)
			out := command.StdOut(o, "run", "user")
			gomega.Expect(out).Should(gomega.ContainSubstring("someuser"))
		})

		ginkgo.It("should correctly parse WORKDIR instruction", func() {
			fileName := "Dockerfile.WorkdirInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n WORKDIR /a/b/c \n CMD pwd"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			command.Run(o, "build", "-f", dockerFilePath, "-t", "workdir", buildContext)
			out := command.StdOut(o, "run", "workdir")
			gomega.Expect(out).Should(gomega.ContainSubstring("/a/b/c"))
		})

		ginkgo.It("should correctly parse ARG instruction", func() {
			fileName := "Dockerfile.ArgInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n ARG user=someuser \n RUN echo ${user}_arg_test"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			out := command.StdErr(o, "build", "-f", dockerFilePath, "--no-cache", "--progress=plain", buildContext)
			gomega.Expect(out).Should(gomega.ContainSubstring("someuser_arg_test"))
		})

		ginkgo.It("should correctly parse SHELL instruction", func() {
			fileName := "Dockerfile.ShellInstruction"
			instructions := "FROM public.ecr.aws/docker/library/alpine:3.13\n SHELL [\"/bin/sh\", \"-c\"] \n RUN echo $0"
			dockerFilePath := filepath.Join(buildContext, fileName)
			ffs.WriteFile(dockerFilePath, instructions)
			out := command.StdErr(o, "build", "-f", dockerFilePath, "--no-cache", "--progress=plain", buildContext)
			gomega.Expect(out).Should(gomega.ContainSubstring("/bin/sh"))
		})

		ginkgo.Context("Docker file syntax negative tests", func() {
			negativeTests := []struct {
				test         string
				fileName     string
				instructions string
				errorMessage string
			}{
				{
					test:         "Empty Dockerfile",
					fileName:     "Dockerfile.Empty",
					instructions: "",
					errorMessage: "Dockerfile cannot be empty",
				},
				{
					test:     "Env no value",
					fileName: "Dockerfile.NoEnv",
					instructions: fmt.Sprintf(`FROM %s
				ENV PATH
				`, defaultImage),
					errorMessage: "ENV must have two arguments",
				},
				{
					test:         "Only comments",
					fileName:     "Dockerfile.OnlyComments",
					instructions: "# Hello\n# These are just comments",
					errorMessage: "file with no instructions",
				},
			}

			for _, test := range negativeTests {
				test := test
				ginkgo.It("should not successfully build a container", func() {
					dockerFilePath := filepath.Join(buildContext, test.fileName)
					ffs.WriteFile(dockerFilePath, test.instructions)
					stdErr := command.RunWithoutSuccessfulExit(o, "build", "-f", dockerFilePath, buildContext).Err.Contents()
					gomega.Expect(stdErr).Should(gomega.ContainSubstring(test.errorMessage))
				})
			}
		})
	})
}
