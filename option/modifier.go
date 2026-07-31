// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package option

// Modifier modifies an Option.
//
// It is not intended to be implemented by code outside this package.
// It is created to provide the flexibility to pass more things to option.New in the future
// without the need to update its signature.
type Modifier interface {
	modify(*Option)
}

type funcModifier struct {
	f func(*Option)
}

func newFuncModifier(f func(*Option)) *funcModifier {
	return &funcModifier{f: f}
}

func (fm *funcModifier) modify(o *Option) {
	fm.f(o)
}

// Env specifies the environment variables to be used during testing. It has the same format as Cmd.Env in os/exec.
func Env(env []string) Modifier {
	return newFuncModifier(func(o *Option) {
		o.env = env
	})
}

// WithNoEnvironmentVariablePassthrough denotes the option does not support environment variable passthrough.
//
// This is useful for disabling tests that require this feature.
func WithNoEnvironmentVariablePassthrough() Modifier {
	return newFuncModifier(func(o *Option) {
		delete(o.features, environmentVariablePassthrough)
	})
}

// WithNerdctlVersion denotes the underlying nerdctl version.
//
// This is useful for tests whose expectations change based on
// the underlying nerdctl version.
func WithNerdctlVersion(version string) Modifier {
	return newFuncModifier(func(o *Option) {
		o.features[nerdctlVersion] = version
	})
}

// WithWindowsHostPathTranslation makes the option rewrite Windows drive-letter
// paths (e.g. `C:\Users\foo`) in command arguments to their WSL2 equivalents
// (e.g. `/mnt/c/Users/foo`) before executing.
//
// This is done in the Finch CLI too.
// See (https://github.com/runfinch/finch/blob/ff1346b1d76f083ba86433e4501cbb5e5ce29634/cmd/finch/nerdctl_windows.go#L72),
// But since common-tests is meant to be executed outside of Finch CLI (e.g Finch Core), we need this as an option.
func WithWindowsHostPathTranslation() Modifier {
	return newFuncModifier(func(o *Option) {
		o.features[windowsHostPathTranslation] = true
	})
}
