package watch

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type SaveOpts struct {
	Watches    Watches
	Connection *sql.DB
}

func Save(opts SaveOpts) error {
	if len(opts.Watches) == 0 {
		return nil
	}
	connection := opts.Connection
	insertColumns := []string{
		"chat_id",
		"collection_id",
		"last_updated",
	}
	onConflictUpdateColumns := []string{
		"last_updated",
	}
	for i := 0; i < len(onConflictUpdateColumns); i++ {
		onConflictUpdateColumns[i] = fmt.Sprintf(
			"%s = EXCLUDED.%s",
			onConflictUpdateColumns[i],
			onConflictUpdateColumns[i],
		)
	}

	// prepare main statement
	globalParameterCount := 1
	var combinedTransactions strings.Builder
	combinedTransactions.WriteString(
		"INSERT INTO watches " +
			"(" + strings.Join(insertColumns, ",") + ") " +
			"VALUES ",
	)

	// populate values
	valueSets := []string{}
	parameterValues := []interface{}{}
	for _, floorPriceDiff := range opts.Watches {
		parameters := []string{}
		for _ = range insertColumns {
			parameters = append(parameters, fmt.Sprintf("$%v", globalParameterCount))
			globalParameterCount += 1
		}
		valueSets = append(valueSets, "("+strings.Join(parameters, ",")+")")
		parameterValues = append(
			parameterValues,
			floorPriceDiff.ChatID,
			floorPriceDiff.CollectionID,
			"NOW()",
		)
	}
	combinedTransactions.WriteString(
		strings.Join(valueSets, ",") +
			" ON CONFLICT (chat_id, collection_id) DO UPDATE SET " +
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
