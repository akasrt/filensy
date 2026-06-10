package cmd

import (
	"fmt"

	"github.com/akasrt/filensy/internal/cli/file"
	"github.com/spf13/cobra"
)

var (
	filePassword string
	isPublic     bool
	fileTTL      string
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
		targetPath := args[0]

		fileService, err := file.NewFileService()
		if err != nil {
			return err
		}

		opts := file.FileOptions{
			Password:    filePassword,
			IsEncrypted: filePassword != "",
			IsPublic:    isPublic,
			TTL:         fileTTL,
		}
		fileData, err := fileService.UploadFile(targetPath, opts)
		if err != nil {
			return err
		}

		fmt.Println("Upload Complete")
		fmt.Println("File Code: ", fileData.Code)
		fmt.Println("File Access Token: ", *fileData.Token)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVarP(&filePassword, "password", "p", "", "Password to encrypt the file")
	uploadCmd.Flags().BoolVar(&isPublic, "public", false, "Set file visibility to public")
	uploadCmd.Flags().StringVar(&fileTTL, "ttl", "", "Time to live for the file (e.g., '1h', '7d')")
}
