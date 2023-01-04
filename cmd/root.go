package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "photosorter",
	Short: "photosorter is a tool for sorting your photos using it's metadata",
	Long:  "photosorter is a tool for sorting your photos using it's metadata",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
