package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var url = "http://localhost:8080"

type WrongStatusError struct {
	Status int
}

func (e WrongStatusError) Error() string {
	return fmt.Sprintf("error response status: %v", e.Status)
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Application allows compress image and convert it to another type",
	Long: `This cli is written on Go. This app communicate with the server and 
	convert images on the server side`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
