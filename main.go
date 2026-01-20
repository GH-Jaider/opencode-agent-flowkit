package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/inalambria/opencode-agent-flowkit/internal/installer"
)

var version = "0.1.0"

func main() {
	// Define flags
	force := flag.Bool("force", false, "Overwrite existing files")
	forceShort := flag.Bool("f", false, "Overwrite existing files (shorthand)")
	showVersion := flag.Bool("version", false, "Show version")
	versionShort := flag.Bool("v", false, "Show version (shorthand)")
	help := flag.Bool("help", false, "Show help")
	helpShort := flag.Bool("h", false, "Show help (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `OpenCode Agent FlowKit v%s

Install OpenCode agent workflow configuration in any repository.

Usage: opencode-agent-flowkit [OPTIONS] [TARGET_DIR]

Arguments:
  TARGET_DIR    Directory to install to (default: current directory)

Options:
  -f, --force   Overwrite existing files
  -v, --version Show version
  -h, --help    Show help

Examples:
  opencode-agent-flowkit              # Install in current dir
  opencode-agent-flowkit ./my-repo    # Install in specific dir
  opencode-agent-flowkit --force      # Overwrite existing files
`, version)
	}

	flag.Parse()

	// Handle version flag
	if *showVersion || *versionShort {
		fmt.Printf("opencode-agent-flowkit v%s\n", version)
		os.Exit(0)
	}

	// Handle help flag
	if *help || *helpShort {
		flag.Usage()
		os.Exit(0)
	}

	// Determine target directory
	targetDir := "."
	if flag.NArg() > 0 {
		targetDir = flag.Arg(0)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	// Check if target directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: directory does not exist: %s\n", absPath)
		os.Exit(1)
	}

	// Combine force flags
	forceInstall := *force || *forceShort

	// Print header
	fmt.Printf("\n📦 OpenCode Agent FlowKit v%s\n\n", version)
	fmt.Printf("Installing to: %s\n\n", absPath)

	// Run installer
	results, err := installer.Install(absPath, forceInstall)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print results
	var created, skipped, errors int
	for _, r := range results {
		relPath, _ := filepath.Rel(absPath, r.Path)
		if r.Error != nil {
			fmt.Printf("❌ Error: %s - %v\n", relPath, r.Error)
			errors++
		} else if r.Skipped {
			fmt.Printf("⏭️  Skipped: %s (already exists)\n", relPath)
			skipped++
		} else if r.Created {
			fmt.Printf("✅ Created: %s\n", relPath)
			created++
		}
	}

	// Print summary
	fmt.Println()
	if errors > 0 {
		fmt.Printf("⚠️  Completed with errors: %d created, %d skipped, %d errors\n", created, skipped, errors)
		os.Exit(1)
	}

	if created == 0 && skipped > 0 {
		fmt.Printf("ℹ️  No changes: all %d files already exist. Use --force to overwrite.\n", skipped)
	} else {
		fmt.Printf("🎉 Done! %d files installed", created)
		if skipped > 0 {
			fmt.Printf(", %d skipped", skipped)
		}
		fmt.Println()
	}

	// Print next steps
	if created > 0 {
		fmt.Println(`
💡 Next steps:
   1. Edit .opencode/config.json to configure models (optional)
   2. Create AGENTS.md with your project's code style (optional)
   3. Start planning: opencode -> select "plan" mode`)
	}
	fmt.Println()
}
