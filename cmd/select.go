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
	Short: "Interactively select SDKs for the environment",
	Run: func(cmd *cobra.Command, args []string) {
		sdks, err := internal.DiscoverSDKs(sdkRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error discovering SDKs: %v\n", err)
			os.Exit(1)
		}
		cfg := &internal.EnvConfig{}
		reader := bufio.NewReader(os.Stdin)
		for _, lang := range []string{"dotnet", "python", "nodejs", "golang", "java"} {
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
				} else {
					fmt.Fprintf(os.Stderr, "Invalid selection for %s, using default.\n", lang)
				}
			}
			switch lang {
			case "dotnet":
				cfg.Dotnet = ver
			case "python":
				cfg.Python = ver
			case "nodejs":
				cfg.Nodejs = ver
			case "golang":
				cfg.Golang = ver
			case "java":
				cfg.Java = ver
			}
		}
		if err := internal.SaveEnvConfig("envman.json", cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write envman.json: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("envman.json updated with selected SDK versions.")
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
