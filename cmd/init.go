package cmd

import (
	"fmt"
	"envman/internal"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var sdkRoot string

func getDefaultSDKRoot() string {
	exePath, err := os.Executable()
	if err != nil {
		return "./envman_sdks" // fallback
	}
	dir := filepath.Dir(exePath)
	return filepath.Join(dir, "envman_sdks")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new environment in the current folder (dynamic, YAML DSL)",
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := filepath.Join(".", "envman.json")
		if _, err := os.Stat(cfgPath); err == nil {
			fmt.Println("envman.json already exists in this folder.")
			return
		}
		root := sdkRoot
		if root == "" {
			root = getDefaultSDKRoot()
		}
		sdks, err := internal.DiscoverSDKs(root)
		if err != nil {
			fmt.Println("Error discovering SDKs:", err)
			return
		}
		// Dynamically list toolchains from toolchains folder
		toolchainFiles, err := os.ReadDir("toolchains")
		if err != nil {
			fmt.Println("Failed to read toolchains folder:", err)
			return
		}
		toolchains := make([]string, 0)
		for _, f := range toolchainFiles {
			if f.IsDir() || !strings.HasPrefix(f.Name(), "envman_") || !strings.HasSuffix(f.Name(), ".yaml") {
				continue
			}
			name := strings.TrimSuffix(strings.TrimPrefix(f.Name(), "envman_"), ".yaml")
			toolchains = append(toolchains, name)
		}
		fmt.Println("Select a toolchain to initialize:")
		for i, t := range toolchains {
			fmt.Printf("  [%d] %s\n", i+1, t)
		}
		fmt.Print("Enter number [1]: ")
		var sel int
		_, err = fmt.Scanf("%d\n", &sel)
		if err != nil || sel < 1 || sel > len(toolchains) {
			sel = 1
		}
		selected := toolchains[sel-1]
		infos := sdks[selected]
		if len(infos) == 0 {
			fmt.Printf("No SDKs found for %s in %s\n", selected, root)
			return
		}
		ver := infos[len(infos)-1].Version
		fmt.Printf("Using %s version: %s\n", selected, ver)
		cfgFile := filepath.Join("toolchains", "envman_"+selected+".yaml")
		dsl, err := internal.LoadToolchainConfig(cfgFile)
		if err != nil {
			fmt.Println("Failed to load toolchain config:", err)
			return
		}
		cwd, _ := os.Getwd()
		vars := map[string]string{
			"version": ver,
			"cwd_basename": filepath.Base(cwd),
			"sdk_root": root,
		}
		if err := internal.RunToolchainSteps(dsl, vars); err != nil {
			fmt.Println("Error during toolchain init:", err)
			return
		}
		cfg := make(map[string]string)
		cfg[selected] = ver
		if err := internal.SaveEnvConfig(cfgPath, cfg); err != nil {
			fmt.Println("Failed to write envman.json:", err)
			return
		}
		fmt.Printf("Initialized new %s environment with version %s.\n", selected, ver)
		// Optionally copy templates
		templates := []string{".gitignore", "README.md"}
		for _, t := range templates {
			src := filepath.Join(root, "..", "templates", t)
			dst := filepath.Join(".", t)
			_ = internal.CopyTemplate(src, dst)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&sdkRoot, "sdk-root", "", "Path to SDKs root folder (default: next to envman.exe)")
}
