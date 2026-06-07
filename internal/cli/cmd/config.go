package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "configure settings",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("requires exactly 2 arguments: [key] and [value]")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		switch key {
		case "auth":
		case "dir":
		default:
			return fmt.Errorf("invalid configuration key")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
