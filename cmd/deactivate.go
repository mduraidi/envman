package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Revert environment variables to previous state (manual step)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "To deactivate, close your shell or manually restore your previous PATH and environment variables.")
		fmt.Fprintln(os.Stderr, "(Optional: implement a more advanced deactivate script if needed.)")
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)
}
