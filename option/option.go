// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package option customizes how tests are run.
package option

import (
	"errors"
	"os"
	"os/exec"
	"strings"
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

// UpdateEnv updates the environment variable for the key name of the input.
func (o *Option) UpdateEnv(env string) {
	if i, exists := containsEnv(o.env, env); exists {
		o.env[i] = env
	} else {
		o.env = append(o.env, env)
	}
}

// DeleteEnv deletes the environment variable for the key name of the input.
func (o *Option) DeleteEnv(env string) {
	if i, exists := containsEnv(o.env, env); exists {
		o.env = append(o.env[:i], o.env[i+1:]...)
	}
}

// containsEnv determines whether an environment variable exists.
func containsEnv(envs []string, targetEnv string) (int, bool) {
	for i, env := range envs {
		envKey := strings.Split(env, "=")[0]
		targetEnvKey := strings.Split(targetEnv, "=")[0]
		if envKey == targetEnvKey {
			return i, true
		}
	}

	return -1, false
}
