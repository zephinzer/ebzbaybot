package watch

import (
	"database/sql"
	"fmt"

	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type ExistsOpts struct {
	Connection   *sql.DB
	ChatID       int64
	CollectionID string
}

func Exists(opts ExistsOpts) (bool, error) {
	connection := opts.Connection
	query := "SELECT COUNT(*) FROM watches"
	queryParams := []interface{}{}
	query += " WHERE chat_id = $1 AND collection_id = $2;"
	queryParams = append(queryParams, opts.ChatID, opts.CollectionID)
	log.Debug(query)
	rows, err := connection.Query(query, queryParams...)
	if err != nil {
		return false, fmt.Errorf("failed to get collections: %s", err)
	}
	defer rows.Close()
	watchExists := 0
	for rows.Next() {
		var watchCount int
		rows.Scan(&watchCount)
		watchExists += watchCount
	}
	return watchExists > 0, nil
}
