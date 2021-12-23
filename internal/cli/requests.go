package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

const defaultID = -1

var (
	requestID       int
	deleteRequestID int
)

// requestsCmd represents the requests command.
var requestsCmd = &cobra.Command{
	Use:   "requests",
	Short: "Working with requests",
	Long:  `This command is used to work with the requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("requests called")
		if requestID != defaultID {
			return oneRequest(requestID)
		}
		if deleteRequestID != defaultID {
			return deleteRequest(deleteRequestID)
		}
		return allRequests()
	},
}

func allRequests() error {
	req, err := http.NewRequest(http.MethodGet, url+"/requests", nil)
	if err != nil {
		return err
	}

	err = auth(req)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var prettyJSON bytes.Buffer

	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON.String())

	return nil
}

func oneRequest(id int) error {
	req, err := http.NewRequest(http.MethodGet, url+"/requests/"+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	err = auth(req)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var prettyJSON bytes.Buffer

	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON.String())

	return nil
}

func deleteRequest(id int) error {
	req, err := http.NewRequest(http.MethodDelete, url+"/requests/"+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	err = auth(req)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status code: %s", resp.Status)
	}

	fmt.Println("Successful delete")

	return nil
}

func init() {
	rootCmd.AddCommand(requestsCmd)

	requestsCmd.Flags().IntVar(&requestID, "id", defaultID, "get reqeust with provided id")
	requestsCmd.Flags().IntVarP(&deleteRequestID, "delete", "d", defaultID, "delete request with provided id")
}
