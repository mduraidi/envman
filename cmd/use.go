package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"envman/internal"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [env]",
	Short: "Switch between local/global environments",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var cfgPath string
		var envType string
		if len(args) == 0 || args[0] == "local" {
			cfgPath = filepath.Join(".", "envman.json")
			envType = "local"
		} else if args[0] == "global" {
			cfgPath = filepath.Join(os.Getenv("USERPROFILE"), ".envman", "envman.json")
			envType = "global"
		} else {
			fmt.Fprintf(os.Stderr, "Unknown environment: %s. Use 'local' or 'global'.\n", args[0])
			os.Exit(1)
		}
		cfg, err := internal.LoadEnvConfig(cfgPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not load %s environment: %v\n", envType, err)
			os.Exit(1)
		}
		if err := internal.SaveEnvConfig("envman.json", cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to switch environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Switched to %s environment.\n", envType)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
