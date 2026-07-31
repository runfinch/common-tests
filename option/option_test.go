// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package option

import (
	"testing"
)

func TestSupportsEnvVarPassthrough(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mods   []Modifier
		assert func(*testing.T, *Option)
	}{
		{
			name: "IsEnvVarPassthroughByDefault",
			mods: []Modifier{},
			assert: func(t *testing.T, uut *Option) {
				if !uut.SupportsEnvVarPassthrough() {
					t.Fatal("expected default SupportsEnvVarPassthrough to be true")
				}
			},
		},
		{
			name: "IsNotEnvVarPassthrough",
			mods: []Modifier{
				WithNoEnvironmentVariablePassthrough(),
			},
			assert: func(t *testing.T, uut *Option) {
				if uut.SupportsEnvVarPassthrough() {
					t.Fatal("expected SupportsEnvVarPassthrough to be false")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			uut, err := New([]string{"nerdctl"}, test.mods...)
			if err != nil {
				t.Fatal(err)
			}

			test.assert(t, uut)
		})
	}
}

func TestSupportsWindowsHostPathTranslation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mods   []Modifier
		assert func(*testing.T, *Option)
	}{
		{
			name: "IsNotWindowsHostPathTranslation",
			mods: []Modifier{},
			assert: func(t *testing.T, uut *Option) {
				if uut.SupportsWindowsHostPathTranslation() {
					t.Fatal("expected default SupportsWindowsHostPathTranslation to be false")
				}
			},
		},
		{
			name: "IsWindowsHostPathTranslation",
			mods: []Modifier{
				WithWindowsHostPathTranslation(),
			},
			assert: func(t *testing.T, uut *Option) {
				if !uut.SupportsWindowsHostPathTranslation() {
					t.Fatal("expected SupportsWindowsHostPathTranslation to be true")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			uut, err := New([]string{"nerdctl"}, test.mods...)
			if err != nil {
				t.Fatal(err)
			}

			test.assert(t, uut)
		})
	}
}

func TestTranslateWindowsHostPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "BuildContextPositional",
			in:   []string{"build", "-t", "foo", `C:\Users\foo\finch-test123`},
			want: []string{"build", "-t", "foo", "/mnt/c/Users/foo/finch-test123"},
		},
		{
			name: "BuildDockerfileAndContext",
			in:   []string{"build", "-f", `D:\a\Dockerfile`, "--no-cache", `D:\a`},
			want: []string{"build", "-f", "/mnt/d/a/Dockerfile", "--no-cache", "/mnt/d/a"},
		},
		{
			name: "BuildOutputDest",
			in:   []string{"build", "-t", "output:tag", "--output", `type=tar,dest=C:\Users\foo\out.tar`, `C:\ctx`},
			want: []string{"build", "-t", "output:tag", "--output", "type=tar,dest=/mnt/c/Users/foo/out.tar", "/mnt/c/ctx"},
		},
		{
			name: "BuildOutputEqualsForm",
			in:   []string{"build", `--output=type=docker`, `C:\ctx`},
			want: []string{"build", "--output=type=docker", "/mnt/c/ctx"},
		},
		{
			name: "BuildSecretSrc",
			in:   []string{"build", "--secret", `id=mysecret,src=C:\Users\foo\secret.txt`, "-f", `C:\Users\foo\Dockerfile`, `C:\Users\foo`},
			want: []string{
				"build", "--secret", "id=mysecret,src=/mnt/c/Users/foo/secret.txt",
				"-f", "/mnt/c/Users/foo/Dockerfile", "/mnt/c/Users/foo",
			},
		},
		{
			name: "SaveOutputFlag",
			in:   []string{"save", "-o", `C:\Users\foo\test.tar`, "alpine:latest"},
			want: []string{"save", "-o", "/mnt/c/Users/foo/test.tar", "alpine:latest"},
		},
		{
			name: "SaveOutputFlagBeforeImages",
			in:   []string{"save", "--output", `C:\Users\foo\test.tar`, "alpine:latest", "alpine:3.13"},
			want: []string{"save", "--output", "/mnt/c/Users/foo/test.tar", "alpine:latest", "alpine:3.13"},
		},
		{
			name: "LoadInputFlag",
			in:   []string{"load", "-i", `C:\Users\foo\test.tar`},
			want: []string{"load", "-i", "/mnt/c/Users/foo/test.tar"},
		},
		{
			name: "ComposeFileFlag",
			in:   []string{"compose", "up", "--file", `C:\Users\foo\docker-compose.yml`},
			want: []string{"compose", "up", "--file", "/mnt/c/Users/foo/docker-compose.yml"},
		},
		{
			name: "CpHostToContainer",
			in:   []string{"cp", `C:\Users\foo\test-file`, "finch-test-ctr:/tmp/test-file"},
			want: []string{"cp", "/mnt/c/Users/foo/test-file", "finch-test-ctr:/tmp/test-file"},
		},
		{
			name: "CpContainerToHost",
			in:   []string{"cp", "finch-test-ctr:/tmp/test-file", `C:\Users\foo\test-file`},
			want: []string{"cp", "finch-test-ctr:/tmp/test-file", "/mnt/c/Users/foo/test-file"},
		},
		{
			name: "CpWithFollowLinkFlagAndContainerSpec",
			in:   []string{"cp", "-L", `C:\Users\foo\symlink`, "finch-test-ctr:/tmp/test-file"},
			want: []string{"cp", "-L", "/mnt/c/Users/foo/symlink", "finch-test-ctr:/tmp/test-file"},
		},
		{
			name: "ExecEnvFileFlag",
			in:   []string{"exec", "--env-file", `C:\Users\foo\env`, "finch-test-ctr", "env"},
			want: []string{"exec", "--env-file", "/mnt/c/Users/foo/env", "finch-test-ctr", "env"},
		},
		{
			name: "RunEnvFileFlag",
			in:   []string{"run", "--rm", "--env-file", `C:\Users\foo\env`, "alpine:latest", "env"},
			want: []string{"run", "--rm", "--env-file", "/mnt/c/Users/foo/env", "alpine:latest", "env"},
		},
		{
			name: "RunCidFileFlag",
			in:   []string{"run", "-d", "--cidfile", `C:\Users\foo\test.cid`, "alpine:latest"},
			want: []string{"run", "-d", "--cidfile", "/mnt/c/Users/foo/test.cid", "alpine:latest"},
		},
		{
			name: "RunLabelFileFlag",
			in:   []string{"run", "--name", "ctr", "--label-file", `C:\Users\foo\label-file`, "alpine:latest"},
			want: []string{"run", "--name", "ctr", "--label-file", "/mnt/c/Users/foo/label-file", "alpine:latest"},
		},
		{
			name: "RunBindMountHostPathOnly",
			in:   []string{"run", "-v", `C:\host:/container`, "alpine:3.13"},
			want: []string{"run", "-v", "/mnt/c/host:/container", "alpine:3.13"},
		},
		{
			name: "RunNamedVolumeUntouched",
			in:   []string{"run", "-v", "foo:/usr/share", "--name", "ctr", "alpine:3.13"},
			want: []string{"run", "-v", "foo:/usr/share", "--name", "ctr", "alpine:3.13"},
		},
		{
			name: "RunAnonymousVolumeUntouched",
			in:   []string{"run", "-v", "/usr/share", "--name", "ctr", "alpine:3.13"},
			want: []string{"run", "-v", "/usr/share", "--name", "ctr", "alpine:3.13"},
		},
		{
			name: "RunBindMountWindowsPathBothSides",
			in:   []string{"run", "-v", `C:\host:C:\host`, "alpine:3.13"},
			want: []string{"run", "-v", "/mnt/c/host:/mnt/c/host", "alpine:3.13"},
		},
		{
			name: "RunWorkdirWindowsPath",
			in:   []string{"run", "-w", `D:\a\finch`, "alpine:3.13"},
			want: []string{"run", "-w", "/mnt/d/a/finch", "alpine:3.13"},
		},
		{
			name: "RunWorkdirLongFlagWindowsPath",
			in:   []string{"run", "--workdir", `C:\Users\foo`, "alpine:3.13"},
			want: []string{"run", "--workdir", "/mnt/c/Users/foo", "alpine:3.13"},
		},
		{
			name: "RunWorkdirLinuxPathUntouched",
			in:   []string{"run", "-w", "/app", "alpine:3.13"},
			want: []string{"run", "-w", "/app", "alpine:3.13"},
		},
		{
			name: "PsVolumeFilterWindowsPath",
			in:   []string{"ps", "-a", "--format", "{{.Names}}", "--filter", `volume=D:\a\finch`},
			want: []string{"ps", "-a", "--format", "{{.Names}}", "--filter", "volume=/mnt/d/a/finch"},
		},
		{
			name: "PsVolumeFilterEqualsForm",
			in:   []string{"ps", `--filter=volume=C:\host`},
			want: []string{"ps", "--filter=volume=/mnt/c/host"},
		},
		{
			name: "PsNameFilterUntouched",
			in:   []string{"ps", "--filter", "name=ctr_1"},
			want: []string{"ps", "--filter", "name=ctr_1"},
		},
		{
			name: "RunMountBindSource",
			in:   []string{"run", "-d", "--mount", `type=bind,source=C:\Users\foo,target=/app`, "alpine:3.13"},
			want: []string{"run", "-d", "--mount", "type=bind,source=/mnt/c/Users/foo,target=/app", "alpine:3.13"},
		},
		{
			name: "RunMountBindSrcReadonly",
			in:   []string{"run", "--mount", `type=bind,src=C:\host,target=/app,ro`, "alpine:3.13"},
			want: []string{"run", "--mount", "type=bind,src=/mnt/c/host,target=/app,ro", "alpine:3.13"},
		},
		{
			name: "RunMountBindEqualsForm",
			in:   []string{"run", `--mount=type=bind,source=C:\host,target=/app`, "alpine:3.13"},
			want: []string{"run", "--mount=type=bind,source=/mnt/c/host,target=/app", "alpine:3.13"},
		},
		{
			name: "RunMountVolumeUntouched",
			in:   []string{"run", "--mount", "type=volume,source=myvol,target=/app", "alpine:3.13"},
			want: []string{"run", "--mount", "type=volume,source=myvol,target=/app", "alpine:3.13"},
		},
		{
			name: "RunMountTmpfsUntouched",
			in:   []string{"run", "--mount", "type=tmpfs,target=/app", "alpine:3.13"},
			want: []string{"run", "--mount", "type=tmpfs,target=/app", "alpine:3.13"},
		},
		{
			name: "PullImageTagAndFlagsUntouched",
			in:   []string{"pull", "alpine:3.13", "--platform", "linux/amd64"},
			want: []string{"pull", "alpine:3.13", "--platform", "linux/amd64"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := translateWindowsHostPaths(test.in)
			if len(got) != len(test.want) {
				t.Fatalf("length mismatch: got %v, want %v", got, test.want)
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Errorf("arg %d: got %q, want %q", i, got[i], test.want[i])
				}
			}
		})
	}
}

func TestNerdctlVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mods   []Modifier
		assert func(*testing.T, *Option)
	}{
		{
			name: "IsNerdctlV2ByDefault",
			mods: []Modifier{},
			assert: func(t *testing.T, uut *Option) {
				if !uut.IsNerdctlV2() {
					t.Fatal("expected IsNerdctlV2 to be true")
				}
			},
		},
		{
			name: "IsNerdctlV1",
			mods: []Modifier{
				WithNerdctlVersion("1.7.7"),
			},
			assert: func(t *testing.T, uut *Option) {
				if !uut.IsNerdctlV1() {
					t.Fatal("expected IsNerdctlV1 to be true")
				}
			},
		},
		{
			name: "IsNerdctlV2",
			mods: []Modifier{
				WithNerdctlVersion("2.0.2"),
			},
			assert: func(t *testing.T, uut *Option) {
				if !uut.IsNerdctlV2() {
					t.Fatal("expected IsNerdctlV2 to be true")
				}
			},
		},
		{
			name: "IsPatchedNerdctlV2",
			mods: []Modifier{
				WithNerdctlVersion("2.0.2.m"),
			},
			assert: func(t *testing.T, uut *Option) {
				if !uut.IsNerdctlV2() {
					t.Fatal("expected IsNerdctlV2 to be true")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			uut, err := New([]string{"nerdctl"}, test.mods...)
			if err != nil {
				t.Fatal(err)
			}

			test.assert(t, uut)
		})
	}
}
