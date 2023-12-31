package cli

import (
	"fmt"
	"github.com/Matterlinkk/Dech-Node/block"
	"github.com/Matterlinkk/Dech-Node/message"
	"github.com/Matterlinkk/Dech-Node/routes"
	"github.com/Matterlinkk/Dech-Node/transaction"
	"github.com/Matterlinkk/Dech-Node/transportchan"
	"github.com/Matterlinkk/Dech-Node/user"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:   "node",
	Short: "Starting this application",
	Long:  "Launching processes, transactions/blocks, creating all the arrays needed for the code",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Root command executed")

		r := chi.NewRouter()

		var loggedUser user.User

		messageMap := message.CreateMessageMap()

		db := block.CreateBlockchain()

		channelTnx := make(chan transaction.Transaction)
		channelBlock := make(chan block.Block)

		go func() {
			transportchan.ProcessBlockchain(channelBlock, db)
		}()

		go func() {
			transportchan.ProcessBlock(channelBlock, channelTnx, db, messageMap)
		}()

		routes.RegisterRoutes(r, db, channelTnx, messageMap, &loggedUser)

		go func() {
			if err := http.ListenAndServe(":8080", r); err != nil {
				fmt.Printf("Error in ListenAndServe: %s", err)
			}
		}()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		fmt.Println("Shutting down the server...")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error in execute: %s", err)
	}
}

func init() {
	rootCmd.AddCommand(blockchainCmd)

	createUserCmd.Flags().String("pk", "", "user's private key")
	createUserCmd.Flags().String("nickname", "unknown", "user's nickname")
	createUserCmd.Flags().String("password", "", "user's password")
	rootCmd.AddCommand(createUserCmd)

	findUserByNicknameCmd.Flags().String("nickname", "", "user's nickname")
	rootCmd.AddCommand(findUserByNicknameCmd)

	createTxCmd.Flags().String("receiver", "", "receiver's nickname")
	createTxCmd.Flags().String("data", "", "some data(only string now)")
	rootCmd.AddCommand(createTxCmd)

	showMessage.Flags().String("from", "", "sender's nickname(from is alise)")
	showMessage.Flags().String("to", "", "receiver's nickname(from is alise)")
	rootCmd.AddCommand(showMessage)

	rootCmd.AddCommand(addressBookCmd)

	rootCmd.AddCommand(UserProfileCmd)

	loginUserCmd.Flags().String("nickname", "", "nickname for login")
	loginUserCmd.Flags().String("password", "", "user's password")
	rootCmd.AddCommand(loginUserCmd)
}
