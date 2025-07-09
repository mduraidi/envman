package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envman",
	Short: "Portable SDK environment manager",
	Long:  `Manage and activate portable SDK environments for dotnet, java, python, nodejs, golang.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
