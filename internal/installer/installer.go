// Package installer provides functionality to install OpenCode agent configuration.
package installer

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:templates
var templates embed.FS

// Result represents the outcome of installing a single file.
type Result struct {
	Path    string
	Created bool
	Skipped bool
	Error   error
}

// Install copies the embedded templates to the target directory.
// If force is true, existing files will be overwritten.
func Install(targetDir string, force bool) ([]Result, error) {
	var results []Result

	// Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("creating target directory: %w", err)
	}

	// Walk through embedded templates
	err := fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root "templates" directory
		if path == "templates" {
			return nil
		}

		// Get relative path (remove "templates/" prefix)
		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return fmt.Errorf("getting relative path: %w", err)
		}

		// Determine destination path
		destPath := filepath.Join(targetDir, relPath)

		// Handle special case: config.json and agents/ go to .opencode/
		if relPath == "config.json" || relPath == "plan.md" {
			destPath = filepath.Join(targetDir, ".opencode", relPath)
		} else if filepath.Dir(relPath) == "agents" {
			destPath = filepath.Join(targetDir, ".opencode", relPath)
		}

		// Create directories
		if d.IsDir() {
			if relPath == "agents" {
				// agents/ goes to .opencode/agents/
				destPath = filepath.Join(targetDir, ".opencode", "agents")
			}
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", destPath, err)
			}
			return nil
		}

		// Read file content from embedded FS
		content, err := templates.ReadFile(path)
		if err != nil {
			results = append(results, Result{Path: destPath, Error: err})
			return nil
		}

		// Check if file exists
		if _, err := os.Stat(destPath); err == nil {
			if !force {
				results = append(results, Result{Path: destPath, Skipped: true})
				return nil
			}
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			results = append(results, Result{Path: destPath, Error: err})
			return nil
		}

		// Write file
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			results = append(results, Result{Path: destPath, Error: err})
			return nil
		}

		results = append(results, Result{Path: destPath, Created: true})
		return nil
	})

	if err != nil {
		return results, fmt.Errorf("walking templates: %w", err)
	}

	return results, nil
}
