package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envman",
	Short: "Portable SDK environment manager",
	Long:  getDynamicLongDescription(),
}

func getDynamicLongDescription() string {
	files, err := os.ReadDir("toolchains")
	if err != nil {
		return "Manage and activate portable SDK environments."
	}
	names := []string{}
	for _, f := range files {
		if f.IsDir() || !strings.HasPrefix(f.Name(), "envman_") || !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		name := strings.TrimSuffix(strings.TrimPrefix(f.Name(), "envman_"), ".yaml")
		names = append(names, name)
	}
	if len(names) == 0 {
		return "Manage and activate portable SDK environments."
	}
	return "Manage and activate portable SDK environments for: " + strings.Join(names, ", ") + "."
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
