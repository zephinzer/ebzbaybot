package collection

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zephinzer/ebzbaybot/pkg/ebzbay"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:     "collection",
		Short:   "Gets statistics on a collection",
		Aliases: []string{"c"},
		RunE:    runE,
	}
	return &command
}

func runE(command *cobra.Command, args []string) error {
	collectionStats := ebzbay.GetCollectionStats(ebzbay.CollectionMadMeerkays)
	fmt.Printf("average lowest price: %v\n", collectionStats.AverageLowestTenPrice)
	fmt.Printf("floor price         : %v\n", collectionStats.FloorPrice)
	fmt.Printf("number of listings  : %v\n", collectionStats.Listings)
	return nil
}
