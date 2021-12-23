package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command.
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout the user",
	Long: `This command is used to logout user from the server.
	
	It is simply made by deleting file with jwt token
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("logout called")

		err := deleteJWT()
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Print("you aren't login")
				return nil
			}

			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
