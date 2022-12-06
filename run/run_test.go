// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"flag"
	"strings"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/option"
	"github.com/runfinch/common-tests/tests"
)

// From https://pkg.go.dev/testing#hdr-Main:
// "Command line flags are always parsed by the time test or benchmark functions run."
// As a result, we need to define the custom flags below as global variables so that
// when `flag.Parse` is invoked by the `testing` package, they can also be parsed.
var subject = flag.String("subject", "", "the subject to be tested, potentially containing spaces")

//nolint:paralleltest // TestRun is like TestMain for the e2e tests.
func TestRun(t *testing.T) {
	o, err := option.New(strings.Split(*subject, " "))
	if err != nil {
		t.Fatalf("failed to initialize a testing option: %v", err)
	}

	ginkgo.SynchronizedBeforeSuite(func() []byte {
		tests.SetupLocalRegistry(o)
		return nil
	}, func(bytes []byte) {})

	ginkgo.SynchronizedAfterSuite(func() {
		tests.CleanupLocalRegistry(o)
	}, func() {})

	const description = "Finch Shared E2E Tests"
	ginkgo.Describe(description, func() {
		// Every test should be listed here.
		// TODO: add tests for "system prune" and "network prune" after upgrading nerdctl to v0.23
		tests.Pull(o)
		tests.Rm(o)
		tests.Rmi(o)
		tests.Run(o)
		tests.Start(o)
		tests.Stop(o)
		tests.Cp(o)
		tests.Tag(o)
		tests.Save(o)
		tests.Load(o)
		tests.Build(o)
		tests.Push(o)
		tests.Images(o)
		tests.ComposeBuild(o)
		tests.ComposeDown(o)
		tests.ComposeKill(o)
		tests.ComposePs(o)
		tests.ComposePull(o)
		tests.ComposeLogs(o)
		tests.Create(o)
		tests.Port(o)
		tests.Kill(o)
		tests.Restart(o)
		tests.Stats(o)
		tests.BuilderPrune(o)
		tests.Exec(o)
		tests.Logs(o)
		tests.Login(o)
		tests.Logout(o)
		tests.VolumeCreate(o)
		tests.VolumeInspect(o)
		tests.VolumeLs(o)
		tests.VolumeRm(o)
		tests.VolumePrune(o)
		tests.ImageHistory(o)
		tests.ImageInspect(o)
		tests.ImagePrune(o)
		tests.Info(o)
		tests.Events(o)
		tests.Inspect(o)
		tests.NetworkCreate(o)
		tests.NetworkInspect(o)
		tests.NetworkLs(o)
		tests.NetworkRm(o)
	})

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, description)
}
