/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/Bass-Peerapon/innoctl/cmd/create"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "innoctl",
	Short: "innovasive control tool",
	Long:  `innovasive control tool`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(create.CreateCmd)
}
