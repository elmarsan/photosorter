package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of photosorter",
	Long:  "Print the version number of photosorter",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("photosorter v1.0.0 ")
	},
}
