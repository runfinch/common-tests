// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package command invokes external commands.
//
// It is designed in a way that the users of the `tests` package can also utilize this package when writing their own tests.
package command

import (
	"bufio"
	"io"
	"strings"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/runfinch/common-tests/option"
)

// Command represents a to-be-executed shell command.
type Command struct {
	opt                 *option.Option
	args                []string
	stdout              io.Writer
	stdin               io.Reader
	timeout             time.Duration
	shouldWait          bool
	shouldCheckExitCode bool
	shouldSucceed       bool
}

// New creates a command with the default configuration.
// Invoke the WithXXX methods to configure it.
// After it's configured, invoke Run to run the command.
//
// Also check the RunXXX and StdXXX wrapper functions in this package for simple use cases.
// As a rule of thumb, please don't create a wrapper function for a use case that involves more than one thing.
// For example,
//
//	command.RunWithoutSuccessfulExit(o, arg).Err.Contents()
//
// may be more readable than
//
//	command.RunWithoutSuccessfulExitAndReturnStderr(o, arg)
func New(opt *option.Option, args ...string) *Command {
	return &Command{
		opt:                 opt,
		args:                args,
		stdout:              ginkgo.GinkgoWriter,
		stdin:               nil,
		timeout:             10 * time.Second,
		shouldWait:          true,
		shouldCheckExitCode: true,
		shouldSucceed:       true,
	}
}

// WithTimeout updates the timeout for the session.
func (c *Command) WithTimeout(timeout time.Duration) *Command {
	c.timeout = timeout
	return c
}

// WithTimeoutInSeconds updates the timeout (in seconds) for the session.
func (c *Command) WithTimeoutInSeconds(timeout time.Duration) *Command {
	return c.WithTimeout(timeout * time.Second)
}

// WithoutSuccessfulExit ensures that the exit code of the command is not 0.
func (c *Command) WithoutSuccessfulExit() *Command {
	c.shouldSucceed = false
	c.shouldCheckExitCode = true
	c.shouldWait = true
	return c
}

// WithoutWait disables waiting for a session to finish.
func (c *Command) WithoutWait() *Command {
	c.shouldWait = false
	return c
}

// WithoutCheckingExitCode disables exit code checking after the session ends.
func (c *Command) WithoutCheckingExitCode() *Command {
	c.shouldCheckExitCode = false
	c.shouldWait = true
	return c
}

// WithStdout specifies the output writer for gexec.Start.
func (c *Command) WithStdout(stdout io.Writer) *Command {
	c.stdout = stdout
	return c
}

// WithStdin specifies the input reader for gexec.Start.
func (c *Command) WithStdin(stdin io.Reader) *Command {
	c.stdin = stdin
	return c
}

// Run starts a session and waits for it to finish.
// It's behavior can be modified by using other Command methods.
// It returns the ended session for further assertions.
func (c *Command) Run() *gexec.Session {
	cmd := c.opt.NewCmd(c.args...)
	cmd.Stdin = c.stdin
	session, err := gexec.Start(cmd, c.stdout, ginkgo.GinkgoWriter)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	if !c.shouldWait {
		return session
	}
	session.Wait(c.timeout)

	if !c.shouldCheckExitCode {
		return session
	}

	if c.shouldSucceed {
		gomega.Expect(session).Should(gexec.Exit(0))
	} else {
		gomega.Expect(session).ShouldNot(gexec.Exit(0))
	}

	return session
}

// Run starts a session, waits for it to finish, ensures the exit code to be 0,
// and returns the ended session to be used for assertions.
func Run(o *option.Option, args ...string) *gexec.Session {
	return New(o, args...).Run()
}

// RunWithoutSuccessfulExit starts a session, waits for it to finish, ensures the exit code not to be 0,
// and returns the ended session to be used for assertions.
func RunWithoutSuccessfulExit(o *option.Option, args ...string) *gexec.Session {
	return New(o, args...).WithoutSuccessfulExit().Run()
}

// Stdout invokes Run and returns the stdout.
func Stdout(o *option.Option, args ...string) []byte {
	return Run(o, args...).Out.Contents()
}

// StdoutStr invokes Run and returns the output in string format.
func StdoutStr(o *option.Option, args ...string) string {
	return strings.TrimSpace(string(Stdout(o, args...)))
}

// StdoutAsLines invokes Run and returns the stdout as lines.
func StdoutAsLines(o *option.Option, args ...string) []string {
	return toLines(Run(o, args...).Out)
}

// Stderr invokes Run and returns the stderr.
func Stderr(o *option.Option, args ...string) []byte {
	return Run(o, args...).Err.Contents()
}

// StderrAsLines invokes Run and returns the stderr as lines.
func StderrAsLines(o *option.Option, args ...string) []string {
	return toLines(Run(o, args...).Err)
}

// StderrStr invokes Run and returns the output in string format.
func StderrStr(o *option.Option, args ...string) string {
	return string(Run(o, args...).Err.Contents())
}

// RunWithoutWait starts a session without waiting for it to finish.
func RunWithoutWait(o *option.Option, args ...string) *gexec.Session {
	return New(o, args...).WithoutWait().Run()
}

func toLines(r io.Reader) []string {
	var lines []string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	gomega.Expect(scanner.Err()).ShouldNot(gomega.HaveOccurred())

	return lines
}
