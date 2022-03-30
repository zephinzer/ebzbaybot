package collection

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type SaveOpts struct {
	Collections Collections
	Connection  *sql.DB
}

func Save(opts SaveOpts) error {
	connection := opts.Connection
	insertColumns := []string{
		"id",
		"label",
		"aliases",
		"average_sale_price",
		"floor_price",
		"number_active",
		"number_of_sales",
		"total_royalties",
		"total_volume",
		"last_updated",
	}
	onConflictUpdateColumns := []string{
		"average_sale_price",
		"floor_price",
		"number_active",
		"number_of_sales",
		"total_royalties",
		"total_volume",
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
		"INSERT INTO collections " +
			"(" + strings.Join(insertColumns, ",") + ") " +
			"VALUES ",
	)

	// populate values
	valueSets := []string{}
	parameterValues := []interface{}{}
	for _, collectionInstance := range opts.Collections {
		parameters := []string{}
		for _ = range insertColumns {
			parameters = append(parameters, fmt.Sprintf("$%v", globalParameterCount))
			globalParameterCount += 1
		}
		valueSets = append(valueSets, "("+strings.Join(parameters, ",")+")")
		parameterValues = append(
			parameterValues,
			collectionInstance.ID,
			collectionInstance.Label,
			collectionInstance.Aliases,
			collectionInstance.AverageSalePrice,
			collectionInstance.FloorPrice,
			collectionInstance.NumberActive,
			collectionInstance.NumberOfSales,
			collectionInstance.TotalRoyalties,
			collectionInstance.TotalVolume,
			"NOW()",
		)
	}
	combinedTransactions.WriteString(
		strings.Join(valueSets, ",") +
			" ON CONFLICT (id) DO UPDATE SET " +
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
