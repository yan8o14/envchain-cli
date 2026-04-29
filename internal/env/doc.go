// Package env provides functionality for injecting environment variables
// into child processes. It merges vault-managed secrets with the current
// process environment and executes commands with the combined environment.
//
// Usage:
//
//	vars := map[string]string{"API_KEY": "secret"}
//	injector := env.New(vars)
//	err := injector.Run([]string{"myapp", "--serve"})
package env
