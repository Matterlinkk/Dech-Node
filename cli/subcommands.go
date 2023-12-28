package cli

import (
	"fmt"
	"github.com/Matterlinkk/Dech-Node/routes"
	"github.com/spf13/cobra"
)

var listUsersCmd = &cobra.Command{
	Use:   "user/list",
	Short: "Show all users",
	Run: func(cmd *cobra.Command, args []string) {
		status, text := routes.CallHandler("http://localhost:8080/user/list")
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var blockchainCmd = &cobra.Command{
	Use:   "blockchain/show",
	Short: "Show array of blocks",
	Run: func(cmd *cobra.Command, args []string) {
		status, text := routes.CallHandler("http://localhost:8080/blockchain/show")
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var createUserCmd = &cobra.Command{
	Use:   "user/create",
	Short: "Create user",
	Run: func(cmd *cobra.Command, args []string) {
		pK, _ := cmd.Flags().GetString("pk")
		nickname, _ := cmd.Flags().GetString("nickname")
		address := fmt.Sprintf("http://localhost:8080/user/create?pk=%s&nickname=%s", pK, nickname)
		status, text := routes.CallHandler(address)
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var findUserByNicknameCmd = &cobra.Command{
	Use:   "user/find",
	Short: "Find user by nickname",
	Run: func(cmd *cobra.Command, args []string) {
		nickname, _ := cmd.Flags().GetString("nickname")
		address := fmt.Sprintf("http://localhost:8080/user/find/%s", nickname)
		status, text := routes.CallHandler(address)
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var createTxCmd = &cobra.Command{
	Use:   "tx/create",
	Short: "Create transaction by 2 nicknames sender and receiver",
	Run: func(cmd *cobra.Command, args []string) {
		pK, _ := cmd.Flags().GetString("sender")
		nickname, _ := cmd.Flags().GetString("receiver")
		data, _ := cmd.Flags().GetString("data")
		address := fmt.Sprintf("http://localhost:8080/tnx/create/%s/%s/message?data=%s", pK, nickname, data)
		status, text := routes.CallHandler(address)
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}

var showMessage = &cobra.Command{
	Use:   "message/show",
	Short: "Show all users between 2 nicknames",
	Run: func(cmd *cobra.Command, args []string) {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		address := fmt.Sprintf("http://localhost:8080/message/show/%s/%s", from, to)
		status, text := routes.CallHandler(address)
		fmt.Printf("Status: %d\nBody:%s", status, text)
	},
}
