// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package option customizes how tests are run.
package option

import (
	"errors"
	"os"
	"os/exec"
)

// Option customizes how tests are run.
//
// If a testing function needs special customizations other than the ones specified in Option,
// we can use composition to extend it.
// For example, to test login functionality,
// we may create a struct named LoginOption that embeds Option and contains additional fields like Username and Password.
type Option struct {
	subject []string
	env     []string
}

// New does some sanity checks on the arguments before initializing an Option.
//
// subject specifies the subject to be tested.
// It is intentionally not designed as an (optional) Modifier because it must contain at least one element.
// Essentially it is used as a prefix when invoking all the binaries during testing.
//
// For example, if subject is ["foo", "bar"], then to test pulling a image, the command name would be "foo",
// and the command args would be something like ["bar", "pull", "alpine"].
func New(subject []string, modifiers ...Modifier) (*Option, error) {
	if len(subject) == 0 {
		return nil, errors.New("missing subject")
	}

	o := &Option{subject: subject}
	for _, modifier := range modifiers {
		modifier.modify(o)
	}

	return o, nil
}

// NewCmd creates a command using the stored option and the provided args.
func (o *Option) NewCmd(args ...string) *exec.Cmd {
	cmdName := o.subject[0]
	cmdArgs := append(o.subject[1:], args...) //nolint:gocritic // appendAssign does not apply to our case.
	cmd := exec.Command(cmdName, cmdArgs...)  //nolint:gosec // G204 is not an issue because cmdName is fully controlled by the user.
	cmd.Env = append(os.Environ(), o.env...)
	return cmd
}
