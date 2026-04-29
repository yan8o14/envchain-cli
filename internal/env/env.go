package env

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Injector handles injecting environment variables into a process.
type Injector struct {
	vars map[string]string
}

// New creates a new Injector with the given environment variables.
func New(vars map[string]string) *Injector {
	return &Injector{vars: vars}
}

// BuildEnv merges the provided vars with the current process environment.
// Values in vars take precedence over existing environment variables.
func (i *Injector) BuildEnv() []string {
	existing := os.Environ()
	overridden := make(map[string]bool)

	result := make([]string, 0, len(existing)+len(i.vars))

	for _, e := range existing {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			if val, ok := i.vars[parts[0]]; ok {
				result = append(result, parts[0]+"="+val)
				overridden[parts[0]] = true
				continue
			}
		}
		result = append(result, e)
	}

	for k, v := range i.vars {
		if !overridden[k] {
			result = append(result, k+"="+v)
		}
	}

	return result
}

// Run executes the given command with the injected environment variables.
func (i *Injector) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		return fmt.Errorf("command not found: %s", args[0])
	}

	cmd := exec.Command(path, args[1:]...)
	cmd.Env = i.BuildEnv()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
