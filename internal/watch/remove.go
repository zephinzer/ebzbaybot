package watch

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type RemoveOpts struct {
	Watches    Watches
	Connection *sql.DB
}

func Remove(opts RemoveOpts) error {
	if len(opts.Watches) == 0 {
		return nil
	}
	connection := opts.Connection
	deletionCriteria := []string{}
	parameterCount := 0
	parameterValues := []interface{}{}
	for _, watchToDelete := range opts.Watches {
		deletionCriteria = append(
			deletionCriteria,
			fmt.Sprintf("(chat_id = $%v AND collection_id = $%v)", parameterCount+1, parameterCount+2),
		)
		parameterValues = append(parameterValues, watchToDelete.ChatID, watchToDelete.CollectionID)
		parameterCount += 2
	}
	deletionSelection := strings.Join(deletionCriteria, " OR ")
	var combinedTransactions strings.Builder
	combinedTransactions.WriteString("DELETE FROM watches WHERE " + deletionSelection + ";")

	log.Debug(combinedTransactions.String())
	_, err := connection.Exec(
		combinedTransactions.String(),
		parameterValues...,
	)
	if err != nil {
		return fmt.Errorf("failed to remove: %s", err)
	}

	return nil
}
