package cmd

import (
	"fmt"
	"envman/internal"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var sdkRoot string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new environment in the current folder",
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := filepath.Join(".", "envman.json")
		if _, err := os.Stat(cfgPath); err == nil {
			fmt.Println("envman.json already exists in this folder.")
			return
		}
		sdks, err := internal.DiscoverSDKs(sdkRoot)
		if err != nil {
			fmt.Println("Error discovering SDKs:", err)
			return
		}
		cfg := &internal.EnvConfig{}
		// For each SDK, pick the latest version
		for lang, infos := range sdks {
			if len(infos) > 0 {
				cfgVer := infos[len(infos)-1].Version // assume sorted, pick last
				switch lang {
				case "dotnet":
					cfg.Dotnet = cfgVer
				case "python":
					cfg.Python = cfgVer
				case "nodejs":
					cfg.Nodejs = cfgVer
				case "golang":
					cfg.Golang = cfgVer
				case "java":
					cfg.Java = cfgVer
				}
			}
		}
		if err := internal.SaveEnvConfig(cfgPath, cfg); err != nil {
			fmt.Println("Failed to write envman.json:", err)
			return
		}
		fmt.Println("Initialized new environment with latest SDKs.")
		// Optionally copy templates
		templates := []string{".gitignore", "README.md"}
		for _, t := range templates {
			src := filepath.Join(sdkRoot, "..", "templates", t)
			dst := filepath.Join(".", t)
			_ = internal.CopyTemplate(src, dst)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&sdkRoot, "sdk-root", "../envman_sdks", "Path to SDKs root folder")
}
