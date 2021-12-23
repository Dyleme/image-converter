package cli

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	registerSourceFile string
	registerNickname   string
	registerPassword   string
	registerEmail      string
)

// registerCmd represents the register command.
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registration to the server",
	Long: `This command can be used to register new users to the server
You should provide your nickname and password as the flags to the command.
If you want, yod can also provide your email as flag.

Login is realized by the saving jwt token, gotten from the server, to the local file.

If you have json file with the password and nickname you can use -s flag to 
connect using this file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("register called")

		var (
			body []byte
			err  error
		)

		if loginSourceFile != "" {
			body, err = credentialsFromFile(registerSourceFile)

		} else {
			if registerNickname != "" && registerPassword != "" {
				body, err = credentialsFromArgs(registerNickname, registerPassword, registerEmail)

			} else {
				return fmt.Errorf("you should specify the source file or nickname and password")
			}
		}
		if err != nil {
			return err
		}

		resp, err := http.Post(url+"/auth/register", "application/json", bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Successful registration")
		} else {
			fmt.Println(resp.Status)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVarP(&registerSourceFile, "source", "s", "", "-s [file]")
	registerCmd.Flags().StringVarP(&registerNickname, "nickname", "n", "", "-n [nickname]")
	registerCmd.Flags().StringVarP(&registerPassword, "password", "p", "", "-p [password]")
	registerCmd.Flags().StringVar(&registerEmail, "email", "", "--email [email]")
}
