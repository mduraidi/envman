package cmd

import (
	"fmt"
	"os"
	"envman/internal"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available SDKs and versions",
	Run: func(cmd *cobra.Command, args []string) {
		root := sdkRoot
		if root == "" {
			root = getDefaultSDKRoot()
		}
		sdks, err := internal.DiscoverSDKs(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error discovering SDKs: %v\n", err)
			os.Exit(1)
		}
		if len(sdks) == 0 {
			fmt.Fprintln(os.Stderr, "No SDKs found in the specified sdkRoot.")
			os.Exit(1)
		}
		for lang, infos := range sdks {
			fmt.Printf("%s:\n", lang)
			if len(infos) == 0 {
				fmt.Println("  (none found)")
			}
			for _, info := range infos {
				fmt.Printf("  %s\n", info.Version)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
