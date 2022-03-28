package collection

import "github.com/spf13/cobra"

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:     "collections",
		Short:   "Retrieves collections available on Ebisus Bay",
		Aliases: []string{"c"},
		RunE:    runE,
	}
	return &command
}

func runE(command *cobra.Command, args []string) error {
	return nil
}
