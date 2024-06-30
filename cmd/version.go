/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/y16ra/aws-bedrock-cli/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of aws-bedrock-cli",
	Long:  `This is aws-bedrock-cli's version command. It prints the version number of aws-bedrock-cli.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
