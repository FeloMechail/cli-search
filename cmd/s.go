/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// sCmd represents the s command
var sCmd = &cobra.Command{
	Use:   "s",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := LoadConfig()
		if err != nil {
			log.Fatal(err)
		}

		var flags []string
		searchQuery := strings.Join(args, " ")
		urlf, _ := cmd.Flags().GetBool("url")
		enginef, _ := cmd.Flags().GetString("engine")

		if urlf && enginef != "" {
			log.Fatal("You cannot use the -u and -e flags at the same time")
		}
		if urlf {
			flags = append(flags, "u")
		} else if enginef != "" {
			flags = append(flags, "e")
		}

		output, _ := PerformSearch(searchQuery, flags)
		fmt.Println(string(output))
		fmt.Printf("Searching for \"%s\"\n", searchQuery)
	},
}

func init() {
	rootCmd.AddCommand(sCmd)
	sCmd.Flags().BoolP("url", "u", false, "Go to url instead of searching")
	sCmd.Flags().StringVarP(&engine, "engine", "e", config.DefaultSearch, "Use different search engine")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
