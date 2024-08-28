/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
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
	Run: func(cmd *cobra.Command, args []string) {
		file, err := loadConfig("cmd/config.yaml")
		if err != nil {
			log.Fatalf("Could not open config file %v", err)
		}
		pathf, _ := cmd.Flags().GetBool("showpath")
		showConfigf, _ := cmd.Flags().GetBool("showconfig")
		if pathf {
			showConfigPath()
		}
		if showConfigf {
			showConfig(file)
		}
	},
}

func showConfigPath() {
	fmt.Print("PATH: cmd/config.yaml\n")
}

func showConfig(file *Config) {
	for _, engine := range file.SearchEngines {
		fmt.Printf(
			"Name: %s, Shortcut: %s, URL: %s\n",
			engine.Name,
			engine.Shortcut,
			engine.URL,
		)
	}
	fmt.Printf("Default Engine: %s, Default Browser: %s", file.DefaultSearch, file.DefaultBrowser)
}

func init() {
	configCmd.Flags().Bool("showpath", false, "Path to config file")
	configCmd.Flags().Bool("showconfig", false, "Show Current config")

	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
