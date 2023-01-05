package cmd

import (
	"photosorter/pkg/photosorter"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sortCmd)
	sortCmd.PersistentFlags().String("format", "YYYY/MM", "The folder structure used for organise the photos")
}

var sortCmd = &cobra.Command{
	Use:   "sort [#src #dst]",
	Short: "Sort photos contained in a directory by it's original creation date",
	Long:  "Sort photos contained in a directory by it's original creation date",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		photosorter.SortDir(args[0], args[1])
	},
}
