package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zhanjie/jsondiff/internal/diff"
	"github.com/zhanjie/jsondiff/internal/input"
	"github.com/zhanjie/jsondiff/internal/jsonpath"
	"github.com/zhanjie/jsondiff/internal/render"
)

var (
	pathFlag  string
	colorFlag string
	keyFlag   string
)

var rootCmd = &cobra.Command{
	Use:   "jsondiff <source> <target>",
	Short: "Compare two JSON values and display differences",
	Long: `jsondiff compares two JSON values (files or inline strings) and displays
the differences with colored terminal output.

Examples:
  jsondiff old.json new.json
  jsondiff '{"a":1}' '{"a":2}'
  jsondiff config.json '{"debug":true}'
  jsondiff server_old.json server_new.json --path "database.connection"
  jsondiff old.json new.json --key "category"    # Match array objects by "category" field`,
	Args: cobra.ExactArgs(2),
	RunE: runDiff,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}

func init() {
	rootCmd.Flags().StringVar(&pathFlag, "path", "", "Only compare the specified JSON sub-path (e.g. \"users.0.name\")")
	rootCmd.Flags().StringVar(&colorFlag, "color", "auto", "Color mode: auto | always | never")
	rootCmd.Flags().StringVar(&keyFlag, "key", "", "Field name to match array objects (e.g. \"id\", \"name\", \"category\")")
}

func runDiff(cmd *cobra.Command, args []string) error {
	// Configure color output
	configureColor(colorFlag)

	// Resolve inputs
	oldVal, err := input.Resolve(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse source: %w", err)
	}
	newVal, err := input.Resolve(args[1])
	if err != nil {
		return fmt.Errorf("failed to parse target: %w", err)
	}

	// Apply JSON path filter if specified
	if pathFlag != "" {
		oldVal, err = jsonpath.Extract(oldVal, pathFlag)
		if err != nil {
			return fmt.Errorf("failed to extract path from source: %w", err)
		}
		newVal, err = jsonpath.Extract(newVal, pathFlag)
		if err != nil {
			return fmt.Errorf("failed to extract path from target: %w", err)
		}
	}

	// Compare
	opts := diff.Options{KeyField: keyFlag}
	diffs := diff.CompareWithOpts(oldVal, newVal, "$", opts)

	// Render
	render.Render(diffs, os.Stdout)

	// Exit code: 0 if equal, 1 if different
	if len(diffs) > 0 {
		os.Exit(1)
	}
	return nil
}

func configureColor(mode string) {
	switch mode {
	case "always":
		color.NoColor = false
	case "never":
		color.NoColor = true
	case "auto":
		// fatih/color auto-detects TTY by default
	default:
		fmt.Fprintf(os.Stderr, "Warning: unknown color mode %q, using auto\n", mode)
	}
}
