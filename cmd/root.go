package cmd

import (
	"github.com/muesli/coral"
)

var rootCmd = &coral.Command{
	Use:   "meterbot",
	Short: "Prepaid meter balance alert bot",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// coral already prints the error; non-zero exit is handled by Execute.
		return
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
}
