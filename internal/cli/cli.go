package cli

import (
	"fmt"
	"io"
)

// Run is kept for backward compatibility with the previous single-command entrypoint.
// The new CLI uses Execute() with subcommands.
func Run(_ []string, _ io.Writer, stderr io.Writer) int {
	if err := Execute(); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}
