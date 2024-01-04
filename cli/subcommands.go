package cli

import (
	"fmt"
	"github.com/Matterlinkk/Dech-Node/routes"
	"github.com/spf13/cobra"
)

var blockchainCmd = &cobra.Command{
	Use:   "blockchain/show",
	Short: "Show array of blocks", //   .\node.exe blockchain/show
	Run: func(cmd *cobra.Command, args []string) {
		status, text := routes.CallHandler("http://localhost:8080/blockchain/show")
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var createUserCmd = &cobra.Command{
	Use:   "user/create",
	Short: "Create user", //  .\node.exe user/create --pk={privatekey} --nickname={username} --password={password}
	Run: func(cmd *cobra.Command, args []string) {
		pK, _ := cmd.Flags().GetString("pk")
		nickname, _ := cmd.Flags().GetString("nickname")
		password, _ := cmd.Flags().GetString("password")
		address := fmt.Sprintf("http://localhost:8080/user/create?pk=%s&nickname=%s&password=%s", pK, nickname, password)
		status, text := routes.CallHandler(address)
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var findUserByNicknameCmd = &cobra.Command{
	Use:   "user/find",
	Short: "Find user by nickname", //.\node.exe user/find --nickname={username}

	Run: func(cmd *cobra.Command, args []string) {
		nickname, _ := cmd.Flags().GetString("nickname")
		address := fmt.Sprintf("http://localhost:8080/user/find/%s", nickname)
		status, text := routes.CallHandler(address)
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var createTxCmd = &cobra.Command{
	Use:   "tx/create",
	Short: "Create transaction by 2 nicknames sender and receiver", // .\node tx/create --receiver={nickname} --data={some data(string only)}

	Run: func(cmd *cobra.Command, args []string) {
		receiver, _ := cmd.Flags().GetString("receiver")
		data, _ := cmd.Flags().GetString("data")
		address := fmt.Sprintf("http://localhost:8080/tnx/create/%s/message?data=%s", receiver, data)
		status, text := routes.CallHandler(address)
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var showMessage = &cobra.Command{
	Use:   "message/show",
	Short: "Show all users between 2 nicknames", // .\node message/show --from={sender's nickname} --to={receiver's nickname}

	Run: func(cmd *cobra.Command, args []string) {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		address := fmt.Sprintf("http://localhost:8080/message/show/%s/%s", from, to)
		status, text := routes.CallHandler(address)
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var addressBookCmd = &cobra.Command{
	Use:   "addressbook",
	Short: "Show map [nickname]address", //  .\node.exe addressbook
	Run: func(cmd *cobra.Command, args []string) {
		status, text := routes.CallHandler("http://localhost:8080/addressbook/show")
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var UserProfileCmd = &cobra.Command{
	Use:   "user/profile",
	Short: "Show user which logged", //  .\node user/profile

	Run: func(cmd *cobra.Command, args []string) {
		status, text := routes.CallHandler("http://localhost:8080/user/profile")
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var loginUserCmd = &cobra.Command{
	Use:   "user/login",
	Short: "Login user", // .\node user/login --nickname={nickname} --password={password}

	Run: func(cmd *cobra.Command, args []string) {
		nickname, _ := cmd.Flags().GetString("nickname")
		password, _ := cmd.Flags().GetString("password")
		address := fmt.Sprintf("http://localhost:8080/user/login/%s?password=%s", nickname, password)
		status, text := routes.CallHandler(address)
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}
