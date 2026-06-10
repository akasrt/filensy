package cmd

import (
	"fmt"
	"time"

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
		fmt.Println("Size: ", formatSize(fileData.Size))
		fmt.Println("Is_Encrypted: ", fileData.Is_Encrypted)
		fmt.Println("Visibility: ", fileData.Visibility)
		fmt.Println("Uploaded At: ", formatTime(fileData.CreatedAt))
		fmt.Println("Expires At: ", formatTime(fileData.ExpiresAt))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().StringVarP(&fileToken, "token", "t", "", "File token for private file")
}

func formatTime(t time.Time) string {
	t = t.Local()
	return t.Format("2006-01-02 03:04:05 PM MST")
}

func formatSize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n = n / unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGT"[exp])
}
