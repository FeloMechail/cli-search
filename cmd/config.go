/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	browser string
	engine  string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pathf, _ := cmd.Flags().GetBool("showpath")
		showConfigf, _ := cmd.Flags().GetBool("showconfig")

		if browser != "" {
			err := SetDefaultBrowser(browser)
			if err != nil {
				return fmt.Errorf("could not set default browser %v", err)
			}
		}

		if engine != "" {
			err := SetDefaultSearchEngine(engine)
			if err != nil {
				return fmt.Errorf("Could not set default search engine %v", err)
			}
		}

		if pathf {
			showConfigPath()
		}

		if showConfigf {
			showConfig()
		}

		return nil
	},
}

func init() {
	configCmd.Flags().Bool("showpath", false, "Path to config file")
	configCmd.Flags().Bool("showconfig", false, "Show current config")
	configCmd.Flags().
		StringVar(&browser, "set-default-browser", "", "Set default browser")
	configCmd.Flags().
		StringVar(&engine, "set-default-engine", "", "Set default search engine")

	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
