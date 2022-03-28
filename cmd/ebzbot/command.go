package ebzbot

import (
	"github.com/spf13/cobra"
	"github.com/zephinzer/ebzbaybot/cmd/ebzbot/get"
	"github.com/zephinzer/ebzbaybot/cmd/ebzbot/start"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:  "ebzbot",
		Long: "This is an unofficial bot for monitoring NFTs on Ebisus's Bay (https://app.ebisusbay.com)",
	}
	command.AddCommand(get.GetCommand())
	command.AddCommand(start.GetCommand())
	return &command
}
