package cmd

import (
	"fmt"

	"github.com/akasrt/filensy/internal/cli/keyring"
	"github.com/spf13/cobra"
)

var passwordCmd = &cobra.Command{
	Use:   "local-password",
	Short: "Retrieve and display the local store encryption password",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := keyring.GetLocalPassword()
		if err != nil || password == "" {
			fmt.Println("Local password not set!")
			return nil
		}

		fmt.Printf("Local store password:\n%s\n", password)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(passwordCmd)
}
