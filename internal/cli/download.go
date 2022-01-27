package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	destPath string
	imageID  int
)

const modeWriteReadExecute = 0o755

// downloadCmd represents the download command.
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads file from server",
	Long: `This command provide you to download image
using it's id on server. You get get this is from the conversation reqeust.
To get in you can run "requests" command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("download called")

		if destPath == "" {
			destPath = ""
		}

		err := downloadFile(imageID, destPath)
		if err != nil {
			return err
		}

		return nil
	},
}

func downloadFile(id int, _ string) error {
	req, err := http.NewRequest(http.MethodGet, url+"/download/image/"+strconv.Itoa(id), http.NoBody)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	err = auth(req)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}
	defer resp.Body.Close()

	values := resp.Header.Values("Content-Disposition")

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	reg := regexp.MustCompile(`filename=".*"`)

	var path string

	for _, val := range values {
		loc := reg.FindStringIndex(val)
		if loc != nil {
			path = val[loc[0]+len(`filename="`) : loc[1]-1]
			break
		}
	}

	file, err := os.OpenFile(dir+"/../../"+path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, modeWriteReadExecute)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().IntVarP(&imageID, "id", "i", defaultID, "--id [download image id]")
	downloadCmd.Flags().StringVarP(&destPath, "destination", "d", "", "-d [path to the file]")

	if err := downloadCmd.MarkFlagRequired("id"); err != nil {
		fmt.Println("flag id was not provided")
	}
}
