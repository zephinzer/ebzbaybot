package ebzbaybot

import (
	"github.com/spf13/cobra"
	"github.com/zephinzer/ebzbaybot/cmd/ebzbaybot/get"
	"github.com/zephinzer/ebzbaybot/cmd/ebzbaybot/start"
	"github.com/zephinzer/ebzbaybot/internal/constants"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:     "ebzbaybot",
		Long:    "This is an unofficial bot for monitoring NFTs on Ebisus's Bay (https://app.ebisusbay.com)",
		Version: constants.Version,
	}
	command.AddCommand(get.GetCommand())
	command.AddCommand(start.GetCommand())
	return &command
}
