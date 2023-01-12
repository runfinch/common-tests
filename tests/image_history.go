// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"os"
	"strings"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/ffs"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// ImageHistory tests "image history" command that shows the history of an image.
func ImageHistory(o *option.Option) {
	ginkgo.Describe("show the history of an image", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			pullImage(o, defaultImage)
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should display image history", func() {
			gomega.Expect(command.StdoutStr(o, "image", "history", defaultImage)).ShouldNot(gomega.BeEmpty())
		})

		for _, quiet := range []string{"-q", "--quiet"} {
			quiet := quiet
			ginkgo.It(fmt.Sprintf("should only display snapshot ID with %s flag", quiet), func() {
				ids := removeMissingID(command.StdoutAsLines(o, "image", "history", quiet, defaultImage))
				gomega.Expect(ids).Should(gomega.HaveEach(gomega.MatchRegexp(sha256RegexFull)))
			})
		}

		ginkgo.It("should only display snapshot ID with --format flag", func() {
			ids := removeMissingID(command.StdoutAsLines(o, "image", "history", defaultImage, "--format", "{{.Snapshot}}"))
			gomega.Expect(ids).Should(gomega.HaveEach(gomega.MatchRegexp(sha256RegexFull)))
		})

		ginkgo.It("should display image history with --no-trunc flag", func() {
			const text = "a very very very very long test phrase that only serves for testing purpose"
			buildContext := ffs.CreateBuildContext(fmt.Sprintf(`FROM %s
			CMD ["echo", %s]
			`, defaultImage, text))
			ginkgo.DeferCleanup(os.RemoveAll, buildContext)

			command.Run(o, "build", "-t", testImageName, buildContext)
			gomega.Expect(command.StdoutStr(o, "image", "history", testImageName)).ShouldNot(gomega.ContainSubstring(text))
			gomega.Expect(command.StdoutStr(o, "image", "history", "--no-trunc", testImageName)).Should(gomega.ContainSubstring(text))
		})
	})
}

func removeMissingID(ids []string) []string {
	var res []string
	for _, id := range ids {
		if !strings.Contains(id, "missing") {
			res = append(res, id)
		}
	}
	return res
}
