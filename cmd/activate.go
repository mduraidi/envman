package cmd

import (
	"fmt"
	"os"
	"strings"
	"envman/internal"
	"github.com/spf13/cobra"
)

var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Set up environment and generate activation scripts",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := internal.LoadEnvConfig("envman.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: envman.json not found or invalid. Run 'envman init' or 'envman select' first.\nDetails: %v\n", err)
			os.Exit(1)
		}
		// Dynamically get supported toolchains from toolchains folder
		toolchainFiles, err := os.ReadDir("toolchains")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read toolchains folder: %v\n", err)
			os.Exit(1)
		}
		toolchains := make([]string, 0)
		for _, f := range toolchainFiles {
			if f.IsDir() || !strings.HasPrefix(f.Name(), "envman_") || !strings.HasSuffix(f.Name(), ".yaml") {
				continue
			}
			name := strings.TrimSuffix(strings.TrimPrefix(f.Name(), "envman_"), ".yaml")
			toolchains = append(toolchains, name)
		}
		m := map[string]string{}
		for _, name := range toolchains {
			if v, ok := cfg[name]; ok && v != "" {
				m[name] = v
			}
		}
		root := sdkRoot
		if root == "" {
			root = getDefaultSDKRoot()
		}
		hasErr := false
		if err := internal.GenerateActivateBat(m, root, "."); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate activate.bat: %v\n", err)
			hasErr = true
		}
		if err := internal.GenerateActivatePs1(m, root, "."); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate activate.ps1: %v\n", err)
			hasErr = true
		}
		if hasErr {
			os.Exit(1)
		}
		fmt.Println("Activation scripts generated: activate.bat, activate.ps1")
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
