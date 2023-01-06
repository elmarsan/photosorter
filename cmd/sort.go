package cmd

import (
	"fmt"
	"os"

	"github.com/elmarsan/photosorter/pkg/photosorter"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sortCmd)
	sortCmd.PersistentFlags().String("format", "month", "The folder structure used for organise the photos")
}

var sortCmd = &cobra.Command{
	Use:   "sort [#src #dst]",
	Short: "Sort photos contained in a directory by it's original creation date",
	Long:  "Sort photos contained in a directory by it's original creation date",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		format := cmd.Flag("format").Value.String()
		err := validateFormatFlag(format)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		photosorter.SortDir(args[0], args[1], format)
	},
}

func validateFormatFlag(format string) error {
	validFormats := []string{"year", "month"}

	for _, valid := range validFormats {
		if format == valid {
			return nil
		}
	}

	return fmt.Errorf("Invalid format flag provided: '%s'. \nValid formats: 'year', 'month'.", format)
}
