// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package option

import (
	"path/filepath"
	"slices"
	"strings"
)

// The path translation logic here is ported from the Finch CLI:
// https://github.com/runfinch/finch/blob/ff1346b1d76f083ba86433e4501cbb5e5ce29634/cmd/finch/nerdctl_windows.go

var aliasMap = map[string]string{
	"build": "image build",
	"run":   "container run",
	"cp":    "container cp",
}

type commandHandler func(args []string) error

type argHandler func(args []string, index int) error

var commandHandlerMap = map[string]commandHandler{
	"container cp": cpHandler,
	"image build":  imageBuildHandler,
}

var argHandlerMap = map[string]map[string]argHandler{
	"image build": {
		"-f":        handleFilePath,
		"--file":    handleFilePath,
		"-o":        handleOutputOption,
		"--output":  handleOutputOption,
		"--secret":  handleSecretOption,
		"--iidfile": handleFilePath,
	},
	"image save": {
		"-o":       handleFilePath,
		"--output": handleFilePath,
	},
	"image load": {
		"-i":      handleFilePath,
		"--input": handleFilePath,
	},
	"save": {
		"-o":       handleFilePath,
		"--output": handleFilePath,
	},
	"load": {
		"-i":      handleFilePath,
		"--input": handleFilePath,
	},
	"container run": {
		"-v":           handleVolume,
		"--volume":     handleVolume,
		"--mount":      handleMount,
		"-w":           handleWorkdir,
		"--workdir":    handleWorkdir,
		"--env-file":   handleFilePath,
		"--cidfile":    handleFilePath,
		"--label-file": handleFilePath,
	},
	"run": {
		"-v":        handleVolume,
		"--volume":  handleVolume,
		"--mount":   handleMount,
		"-w":        handleWorkdir,
		"--workdir": handleWorkdir,
	},
	"create": {
		"-v":        handleVolume,
		"--volume":  handleVolume,
		"--mount":   handleMount,
		"-w":        handleWorkdir,
		"--workdir": handleWorkdir,
	},
	"ps": {
		"-f":       handleFilter,
		"--filter": handleFilter,
	},
	"container ps": {
		"-f":       handleFilter,
		"--filter": handleFilter,
	},
	"exec": {
		"--env-file": handleFilePath,
	},
	"compose": {
		"-f":     handleFilePath,
		"--file": handleFilePath,
	},
}

func translateWindowsHostPaths(args []string) []string {
	if len(args) == 0 {
		return args
	}

	out := slices.Clone(args)

	cmdName := out[0]
	rest := out[1:]

	cmdKey := cmdName
	if alias, ok := aliasMap[cmdName]; ok {
		cmdKey = alias
	} else if len(rest) > 0 {
		if _, ok := commandHandlerMap[cmdName+" "+rest[0]]; ok {
			cmdKey = cmdName + " " + rest[0]
		} else if _, ok := argHandlerMap[cmdName+" "+rest[0]]; ok {
			cmdKey = cmdName + " " + rest[0]
		}
	}

	if h, ok := commandHandlerMap[cmdKey]; ok {
		if err := h(out); err != nil {
			return out
		}
	}

	if aMap, ok := argHandlerMap[cmdKey]; ok {
		for i := range out {
			flag, _, _ := strings.Cut(out[i], "=")
			if h, ok := aMap[flag]; ok {
				if err := h(out, i); err != nil {
					return out
				}
			}
		}
	}

	return out
}

// convertToWSLPath converts an absolute Windows path (C:\Users\foo) to its WSL2 mount equivalent (/mnt/c/Users/foo).
func convertToWSLPath(winPath string) string {
	if len(winPath) < 2 || winPath[1] != ':' {
		return winPath
	}
	drive := strings.ToLower(string(winPath[0]))
	remaining := ""
	if len(winPath) > 2 {
		remaining = winPath[2:]
	}
	remaining = strings.ReplaceAll(remaining, `\`, "/")
	return filepath.ToSlash("/mnt/" + drive + remaining)
}

func isWindowsPath(s string) bool {
	return len(s) >= 3 && s[1] == ':' && (s[2] == '\\' || s[2] == '/')
}

func handleFilePath(args []string, index int) error {
	arg := args[index]
	if strings.Contains(arg, "=") {
		before, after, _ := strings.Cut(arg, "=")
		args[index] = before + "=" + convertToWSLPath(after)
		return nil
	}
	if index+1 < len(args) {
		args[index+1] = convertToWSLPath(args[index+1])
	}
	return nil
}

func handleVolume(args []string, index int) error {
	arg := args[index]

	var (
		value        string
		before       string
		hasEqualForm bool
	)
	switch {
	case strings.Contains(arg, "="):
		before, value, _ = strings.Cut(arg, "=")
		hasEqualForm = true
	case index+1 < len(args):
		value = args[index+1]
	default:
		return nil
	}

	cleanArg := value
	readWrite := ""
	switch {
	case strings.HasSuffix(value, ":ro"), strings.HasSuffix(value, ":rw"):
		readWrite = value[len(value)-3:]
		cleanArg = value[:len(value)-3]
	case strings.HasSuffix(value, ":rro"):
		readWrite = value[len(value)-4:]
		cleanArg = value[:len(value)-4]
	}

	// The host path must be a Windows drive path (eg C:\...).
	if !isWindowsPath(cleanArg) {
		return nil
	}

	// Split the host and container paths at the separator colon that follows the
	// host path, skipping the drive-letter colon at index 1 (eg the ':' between
	// "C:\host" and "/container").
	sep := strings.IndexByte(cleanArg[2:], ':')
	if sep < 0 {
		return nil
	}
	sep += 2
	hostPath := cleanArg[:sep]
	containerPath := cleanArg[sep+1:]

	wslHostPath := convertToWSLPath(hostPath)
	// The container path is normally a Linux path, but some tests reuse the host
	// path (eg `-v <pwd>:<pwd>`); convert it too when it is a Windows path.
	if isWindowsPath(containerPath) {
		containerPath = convertToWSLPath(containerPath)
	}

	if hasEqualForm {
		args[index] = before + "=" + wslHostPath + ":" + containerPath + readWrite
	} else {
		args[index+1] = wslHostPath + ":" + containerPath + readWrite
	}
	return nil
}

func handleMount(args []string, index int) error {
	arg := args[index]

	var (
		value        string
		before       string
		hasEqualForm bool
	)
	switch {
	case strings.Contains(arg, "="):
		before, value, _ = strings.Cut(arg, "=")
		hasEqualForm = true
	case index+1 < len(args):
		value = args[index+1]
	default:
		return nil
	}

	entries := strings.Split(value, ",")
	m := make(map[string]string)
	for _, e := range entries {
		k, v, found := strings.Cut(e, "=")
		if found {
			m[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}

	// Only translate bind mounts; volume/tmpfs sources are not host paths.
	if m["type"] != "bind" {
		return nil
	}

	key := "src"
	hostPath, ok := m[key]
	if !ok {
		key = "source"
		hostPath, ok = m[key]
	}
	if !ok || !strings.Contains(hostPath, `\`) {
		return nil
	}

	for i, e := range entries {
		k, _, found := strings.Cut(e, "=")
		if found && strings.TrimSpace(k) == key {
			entries[i] = k + "=" + convertToWSLPath(hostPath)
		}
	}
	newValue := strings.Join(entries, ",")

	if hasEqualForm {
		args[index] = before + "=" + newValue
	} else {
		args[index+1] = newValue
	}
	return nil
}

func handleWorkdir(args []string, index int) error {
	arg := args[index]
	if strings.Contains(arg, "=") {
		before, after, _ := strings.Cut(arg, "=")
		if isWindowsPath(after) {
			args[index] = before + "=" + convertToWSLPath(after)
		}
		return nil
	}
	if index+1 < len(args) && isWindowsPath(args[index+1]) {
		args[index+1] = convertToWSLPath(args[index+1])
	}
	return nil
}

func handleFilter(args []string, index int) error {
	arg := args[index]

	var (
		value        string
		before       string
		hasEqualForm bool
	)
	switch {
	case strings.HasPrefix(arg, "--filter="):
		before, value, _ = strings.Cut(arg, "=")
		hasEqualForm = true
	case strings.HasPrefix(arg, "-f="):
		before, value, _ = strings.Cut(arg, "=")
		hasEqualForm = true
	case index+1 < len(args):
		value = args[index+1]
	default:
		return nil
	}

	key, path, found := strings.Cut(value, "=")
	if !found || key != "volume" || !isWindowsPath(path) {
		return nil
	}
	newValue := key + "=" + convertToWSLPath(path)

	if hasEqualForm {
		args[index] = before + "=" + newValue
	} else {
		args[index+1] = newValue
	}
	return nil
}

func handleOutputOption(args []string, index int) error {
	return convertKeyValueOptionPath(args, index, "dest")
}

func handleSecretOption(args []string, index int) error {
	return convertKeyValueOptionPath(args, index, "src")
}

func convertKeyValueOptionPath(args []string, index int, pathKey string) error {
	arg := args[index]

	var (
		value        string
		before       string
		hasEqualForm bool
	)
	switch {
	case strings.Contains(arg, "="):
		before, value, _ = strings.Cut(arg, "=")
		hasEqualForm = true
	case index+1 < len(args):
		value = args[index+1]
	default:
		return nil
	}

	entries := strings.Split(value, ",")
	converted := false
	for i, e := range entries {
		k, v, found := strings.Cut(e, "=")
		if found && strings.TrimSpace(k) == pathKey {
			entries[i] = k + "=" + convertToWSLPath(v)
			converted = true
		}
	}
	if !converted {
		return nil
	}
	newValue := strings.Join(entries, ",")

	if hasEqualForm {
		args[index] = before + "=" + newValue
	} else {
		args[index+1] = newValue
	}
	return nil
}

func cpHandler(args []string) error {
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") || arg == "cp" || arg == "container" {
			continue
		}
		if colon := strings.Index(arg, ":"); colon > 1 {
			continue
		}
		args[i] = convertToWSLPath(arg)
	}
	return nil
}

func imageBuildHandler(args []string) error {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return nil
		}
	}
	last := len(args) - 1
	if last < 0 {
		return nil
	}
	if args[last] != "--debug" {
		if isWindowsPath(args[last]) {
			args[last] = convertToWSLPath(args[last])
		}
	} else if last-1 >= 0 && isWindowsPath(args[last-1]) {
		args[last-1] = convertToWSLPath(args[last-1])
	}
	return nil
}
