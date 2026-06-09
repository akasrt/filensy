package cmd

import (
	"fmt"

	"github.com/akasrt/filensy/internal/cli/file"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find [file_code]",
	Short: "Find a file and retrieve metadata with file code",
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

		fileData, err := fileService.FindFile(fileCode, fileToken)
		if err != nil {
			return err
		}

		fmt.Println("File MetaData")
		fmt.Println("Name: ", fileData.Name)
		fmt.Println("Size: ", fileData.Size)
		fmt.Println("Is_Encrypted: ", fileData.Is_Encrypted)
		fmt.Println("Visibility: ", fileData.Visibility)
		fmt.Println("Uploaded At: ", fileData.CreatedAt)
		fmt.Println("Expires At: ", fileData.ExpiresAt)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().StringVarP(&fileToken, "token", "t", "", "File token for private file")
}
