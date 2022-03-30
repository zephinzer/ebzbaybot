package watch

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type LoadOpts struct {
	Connection *sql.DB
	// Only defines a chat ID, when this is
	// defined, the returned Watches only
	// contains Watches by the chat owner
	OnlyFor int64
}

func Load(opts LoadOpts) (Watches, error) {
	connection := opts.Connection
	selectedColumns := []string{
		"chat_id",
		"collection_id",
		"last_updated",
	}
	query := fmt.Sprintf(
		"SELECT %s FROM watches",
		strings.Join(selectedColumns, ", "),
	)
	log.Debug(query)
	queryParams := []interface{}{}
	if opts.OnlyFor != 0 {
		query += " WHERE chat_id = $1;"
		queryParams = append(queryParams, opts.OnlyFor)
	}
	rows, err := connection.Query(query, queryParams...)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %s", err)
	}
	defer rows.Close()
	watches := Watches{}
	for rows.Next() {
		w := Watch{}
		if err := rows.Scan(
			&w.ChatID,
			&w.CollectionID,
			&w.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row[%v]", len(watches))
		}
		watches = append(watches, w)
	}
	return watches, nil
}
