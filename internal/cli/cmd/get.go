package cmd

import (
	"fmt"

	"github.com/akasrt/filensy/internal/cli/file"
	"github.com/spf13/cobra"
)

var (
	getOutputDir string
	getPassword  string
	fileToken    string
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
		fileCode := args[0]

		fileService, err := file.NewFileService()
		if err != nil {
			return err
		}

		err = fileService.GetFile(getOutputDir, fileCode, fileToken, getPassword)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&getOutputDir, "dir", "d", "", "Directory to store the file")
	getCmd.Flags().StringVarP(&getPassword, "password", "p", "", "Password to decrypt file")
	getCmd.Flags().StringVarP(&fileToken, "token", "t", "", "File token for private file")
}
