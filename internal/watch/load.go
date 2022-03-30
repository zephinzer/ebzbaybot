package watch

import (
	"database/sql"
	"fmt"
	"strings"
)

type LoadOpts struct {
	Connection *sql.DB
}

func Load(opts LoadOpts) (Watches, error) {
	connection := opts.Connection
	selectedColumns := []string{
		"chat_id",
		"collection_id",
		"last_updated",
	}
	rows, err := connection.Query(
		fmt.Sprintf(
			"SELECT %s FROM watches",
			strings.Join(selectedColumns, ", "),
		),
	)
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
