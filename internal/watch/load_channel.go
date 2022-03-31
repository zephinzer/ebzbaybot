package watch

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type LoadChannelOpts struct {
	Connection *sql.DB
	// Only defines a chat ID, when this is
	// defined, the returned Watches only
	// contains Watches by the chat owner
	OnlyFor int64
}

func LoadChannel(opts LoadChannelOpts) (ChannelWatches, error) {
	connection := opts.Connection
	selectedColumns := []string{
		"chat_id",
		"collection_id",
		"last_updated",
	}
	query := fmt.Sprintf(
		"SELECT %s FROM watches_channel",
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
		return nil, fmt.Errorf("failed to get channel watches: %s", err)
	}
	defer rows.Close()
	watches := ChannelWatches{}
	for rows.Next() {
		cw := ChannelWatch{}
		if err := rows.Scan(
			&cw.ChatID,
			&cw.CollectionID,
			&cw.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("failed to scan channel watch row[%v]", len(watches))
		}
		watches = append(watches, cw)
	}
	return watches, nil
}
