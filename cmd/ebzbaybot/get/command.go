package get

import (
	"github.com/spf13/cobra"
	"github.com/zephinzer/ebzbaybot/cmd/ebzbaybot/get/collection"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "get",
		Short: "Retrieves information from Ebisus Bay",
	}
	command.AddCommand(collection.GetCommand())
	return &command
}
