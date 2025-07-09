package cmd

import (
	"bufio"
	"fmt"
	"envman/internal"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Interactively select SDKs for the environment (dynamic)",
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
		// Dynamically list toolchains from toolchains folder
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
		cfg := make(map[string]string)
		reader := bufio.NewReader(os.Stdin)
		for _, lang := range toolchains {
			infos := sdks[lang]
			if len(infos) == 0 {
				fmt.Fprintf(os.Stderr, "Warning: No %s SDKs found.\n", lang)
				continue
			}
			fmt.Printf("Available %s versions:\n", lang)
			for i, info := range infos {
				fmt.Printf("  [%d] %s\n", i+1, info.Version)
			}
			fmt.Printf("Select %s version [default: %s]: ", lang, infos[len(infos)-1].Version)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			ver := infos[len(infos)-1].Version
			if input != "" {
				idx := 0
				fmt.Sscanf(input, "%d", &idx)
				if idx > 0 && idx <= len(infos) {
					ver = infos[idx-1].Version
				}
			}
			cfg[lang] = ver
		}
		if err := internal.SaveEnvConfig("envman.json", cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save envman.json: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Updated envman.json with selected SDK versions.")
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
