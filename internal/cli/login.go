package cli

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	loginSourceFile string
	loginNickname   string
	loginPassword   string
	loginEmail      string
)

// loginCmd represents the login command.
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login user to the server",
	Long: `This command can be used to login to the server.
	You should provide your nickname and password as the flags to the command.
	If you want, yod can also provide your email as flag.

	Login is realized by the saving jwt token, gotten from the server, to the local file.

	If you have json file with the password and nickname you can use -s flag to 
	connect using this file.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("login called")

		var (
			body []byte
			err  error
		)

		if loginSourceFile != "" {
			body, err = credentialsFromFile(loginSourceFile)

		} else {
			if loginNickname != "" && loginPassword != "" {
				body, err = credentialsFromArgs(loginNickname, loginPassword, loginEmail)

			} else {
				return fmt.Errorf("you should specify the source file or nickname and password")
			}
		}
		if err != nil {
			return err
		}

		resp, err := http.Post(url+"/auth/login", "application/json", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		_, err = getToken(b)
		if err != nil {
			return err
		}

		err = saveJWT(b)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&loginSourceFile, "source", "s", "", "-s [file]")
	loginCmd.Flags().StringVarP(&loginNickname, "nickname", "n", "", "-n [nickname]")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "-p [password]")
	loginCmd.Flags().StringVar(&loginEmail, "email", "", "--email [email]")
}
