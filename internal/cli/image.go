package cli

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	filePath  string
	convRatio float32
	newType   string
)

// imageCmd represents the image command.
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Add request to convolute image",
	Long: `Add request for the image conversion to the server.
To add reqeust you should provide image by it's path in -p flag.
You should also provide type to convErted image in -t flag
Also you can provide convolution ratio using -r flag`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("image called")
		return addRequest(filePath, convRatio, newType)
	},
}

func addRequest(path string, ratio float32, newType string) error {
	req, err := createMultipartRequest(path, ratio, newType)
	if err != nil {
		return fmt.Errorf("add request: %w", err)
	}

	err = auth(req)
	if err != nil {
		return fmt.Errorf("add request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("add request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("add request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(respBody))
		return fmt.Errorf("wrong status code %s", resp.Status)
	}

	fmt.Println(string(respBody))

	return nil
}

func createMultipartRequest(path string, ratio float32, fileType string) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("create multipart request: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("Image", path)
	if err != nil {
		return nil, fmt.Errorf("create multipart request: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("create multipart request: %w", err)
	}

	err = writer.WriteField("CompressionInfo",
		`{
			"ratio": `+fmt.Sprintf("%v", ratio)+`,
			"newType": "`+fileType+`"
			}`)
	if err != nil {
		return nil, fmt.Errorf("create multipart request: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("create multipart request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url+"/requests/image", body)
	if err != nil {
		return nil, fmt.Errorf("create multipart request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func init() {
	requestsCmd.AddCommand(imageCmd)

	imageCmd.Flags().StringVarP(&filePath, "path", "p", "", "path to the converted image")
	imageCmd.Flags().Float32VarP(&convRatio, "ratio", "r", 1, "convolution ratio")
	imageCmd.Flags().StringVarP(&newType, "type", "t", "", "type of the converted image")

	if err := imageCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println("flag path is not provided")
	}

	if err := imageCmd.MarkFlagRequired("type"); err != nil {
		fmt.Println("flag type is not provided")
	}
}
