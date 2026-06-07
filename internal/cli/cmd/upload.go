package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	filePath     string
	filePassword string
)

var uploadCmd = &cobra.Command{
	Use:   "upload [file_path]",
	Short: "upload files to the server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing file path")
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
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVar(&filePath, "path", "", "path to file")
	uploadCmd.Flags().StringVarP(&filePassword, "password", "p", "", "Password to encrypt the file")
}
