package list

import (
	"github.com/spf13/cobra"
	"github.com/zephinzer/ebzbaybot/cmd/ebzbaybot/list/collection"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "list",
		Short: "Queries information from Ebisus's Bay",
	}
	command.AddCommand(collection.GetCommand())
	return &command
}
