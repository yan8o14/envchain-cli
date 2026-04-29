package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportFormat represents the output format for exporting variables.
type ExportFormat string

const (
	// FormatShell outputs variables as shell export statements.
	FormatShell ExportFormat = "shell"
	// FormatDotenv outputs variables in .env file format.
	FormatDotenv ExportFormat = "dotenv"
)

// Export writes the environment variables to w in the specified format.
func Export(w io.Writer, vars map[string]string, format ExportFormat) error {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := vars[k]
		var line string

		switch format {
		case FormatShell:
			line = fmt.Sprintf("export %s=%s\n", k, shellQuote(v))
		case FormatDotenv:
			line = fmt.Sprintf("%s=%s\n", k, v)
		default:
			return fmt.Errorf("unsupported export format: %s", format)
		}

		if _, err := fmt.Fprint(w, line); err != nil {
			return fmt.Errorf("failed to write variable %s: %w", k, err)
		}
	}

	return nil
}

// shellQuote wraps a value in single quotes, escaping any existing single quotes.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\'''") + "'"
}
