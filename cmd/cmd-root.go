package cmd

import (
	"github.com/spf13/cobra"
)

func generateRootCmd() *cobra.Command {
	var confPath string

	var rootCmd = &cobra.Command{
		Use:   "ferleasy",
		Short: "Manages ferlease releases declaratively from an s3 store",
	}

	rootCmd.PersistentFlags().StringVarP(&confPath, "config", "c", "config.yml", "Path to a yaml configuration file")
	rootCmd.MarkPersistentFlagFilename("config")

	rootCmd.AddCommand(generateAddCmd(&confPath))
	rootCmd.AddCommand(generateRemoveCmd(&confPath))
	rootCmd.AddCommand(generateListCmd(&confPath))
	rootCmd.AddCommand(generateSyncCmd(&confPath))

	return rootCmd
}

func Execute() error {
	return generateRootCmd().Execute()
}