package floorpricediff

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type SaveOpts struct {
	FloorPriceDiffs FloorPriceDiffs
	Connection      *sql.DB
}

func Save(opts SaveOpts) error {
	if len(opts.FloorPriceDiffs) == 0 {
		return nil
	}
	connection := opts.Connection
	insertColumns := []string{
		"collection_id",
		"previous_price",
		"current_price",
		"listing_id",
		"image_url",
		"edition",
		"score",
		"rank",
		"last_updated",
	}
	onConflictUpdateColumns := []string{
		"previous_price",
		"current_price",
		"listing_id",
		"image_url",
		"edition",
		"score",
		"rank",
	}
	for i := 0; i < len(onConflictUpdateColumns); i++ {
		onConflictUpdateColumns[i] = fmt.Sprintf(
			"%s = EXCLUDED.%s",
			onConflictUpdateColumns[i],
			onConflictUpdateColumns[i],
		)
	}
	onConflictUpdateColumns = append(onConflictUpdateColumns, "last_updated = NOW()")

	// prepare main statement
	globalParameterCount := 1
	var combinedTransactions strings.Builder
	combinedTransactions.WriteString(
		"INSERT INTO floor_price_diffs " +
			"(" + strings.Join(insertColumns, ",") + ") " +
			"VALUES ",
	)

	// populate values
	valueSets := []string{}
	parameterValues := []interface{}{}
	for _, floorPriceDiff := range opts.FloorPriceDiffs {
		parameters := []string{}
		for _ = range insertColumns {
			parameters = append(parameters, fmt.Sprintf("$%v", globalParameterCount))
			globalParameterCount += 1
		}
		valueSets = append(valueSets, "("+strings.Join(parameters, ",")+")")
		parameterValues = append(
			parameterValues,
			floorPriceDiff.CollectionID,
			floorPriceDiff.PreviousPrice,
			floorPriceDiff.CurrentPrice,
			floorPriceDiff.ListingID,
			floorPriceDiff.ImageURL,
			floorPriceDiff.Edition,
			floorPriceDiff.Score,
			floorPriceDiff.Rank,
			"NOW()",
		)
	}
	combinedTransactions.WriteString(
		strings.Join(valueSets, ",") +
			" ON CONFLICT (collection_id) DO UPDATE SET " +
			strings.Join(onConflictUpdateColumns, ", "),
	)
	combinedTransactions.WriteString(";\n")

	log.Debug(combinedTransactions.String())
	_, err := connection.Exec(
		combinedTransactions.String(),
		parameterValues...,
	)
	if err != nil {
		return fmt.Errorf("failed to save: %s", err)
	}

	return nil
}
