package cmd

import (
	"fmt"

	"github.com/akasrt/filensy/internal/cli/file"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [file_code]",
	Short: "Delete a file from the server",
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
		fileCode := args[0]

		fileService, err := file.NewFileService()
		if err != nil {
			return err
		}

		err = fileService.DeleteFile(fileCode, fileToken)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&fileToken, "token", "t", "", "File token")
}
