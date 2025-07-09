package cmd

import (
	"fmt"
	"os"
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
		m := map[string]string{
			"dotnet": cfg.Dotnet,
			"python": cfg.Python,
			"nodejs": cfg.Nodejs,
			"golang": cfg.Golang,
			"java":   cfg.Java,
		}
		hasErr := false
		if err := internal.GenerateActivateBat(m, sdkRoot, "."); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate activate.bat: %v\n", err)
			hasErr = true
		}
		if err := internal.GenerateActivatePs1(m, sdkRoot, "."); err != nil {
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
