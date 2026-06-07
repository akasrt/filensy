package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	getOutputDir string
	getPassword  string
	getFileToken string
)

var getCmd = &cobra.Command{
	Use:   "get [file_code]",
	Short: "Download file with the file code",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing file code")
		}
		if len(args) > 1 {
			return fmt.Errorf("too many args")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&getOutputDir, "dir", "d", "", "Directory to store the file")
	getCmd.Flags().StringVarP(&getPassword, "password", "p", "", "Password to decrypt file")
	getCmd.Flags().StringVarP(&getFileToken, "token", "t", "", "File token for private file")
}
