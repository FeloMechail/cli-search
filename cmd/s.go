/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		var flags []string
		searchQuery := strings.Join(args, " ")
		urlf, err := cmd.Flags().GetBool("url")
		if err != nil {
			return err
		}

		enginef, err := cmd.Flags().GetString("engine")
		if err != nil {
			return err
		}

		if urlf && enginef != "" {
			return errors.New(
				"You cannot use the -u and -e flags at the same time",
			)
		}
		if urlf {
			flags = append(flags, "u")
		} else if enginef != "" {
			flags = append(flags, "e")
		}

		output, err := PerformSearch(searchQuery, flags)
		if err != nil {
			return err
		}

		fmt.Print("Command output\n\n")
		fmt.Println(string(output))
		fmt.Print("End of Command output\n\n")
		fmt.Printf("Searching for \"%s\"\n", searchQuery)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sCmd)
	sCmd.Flags().BoolP("url", "u", false, "Go to url instead of searching")
	sCmd.Flags().
		StringVarP(&engine, "engine", "e", config.DefaultSearch, "Use different search engine")
	sCmd.MarkFlagsMutuallyExclusive("url", "engine")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
