// Package cli provides access to the documentation for the CLI commands
package cli

import (
	"embed"
	"fmt"
	"log/slog"
	"strings"
)

//go:embed docs
var docs embed.FS

// Doc retrieves the documentation for a specific command
func Doc(cmd string) (string, error) {
	if cmd == "" {
		cmd = "viv"
	}
	if !strings.HasPrefix(cmd, "viv") {
		cmd = "viv " + cmd
	}

	doc := strings.ReplaceAll(cmd, " ", "_")
	doc = fmt.Sprintf("docs/%s.md", doc)
	content, err := docs.ReadFile(doc)
	if err != nil {
		return "", fmt.Errorf("failed to read documentation for command %s: %w", cmd, err)
	}
	// add example.md content if the command is "viv"
	if cmd == "viv" {
		doc = "docs/example.md"
		example, err := docs.ReadFile(doc)
		if err != nil {
			return "", fmt.Errorf("failed to read root documentation: %w", err)
		}
		content = append(content, example...)
	}
	slog.Info("viv documentation", "command", cmd, "doc", doc)
	return string(content), nil
}
