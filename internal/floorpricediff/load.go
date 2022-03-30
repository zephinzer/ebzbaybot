package floorpricediff

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type LoadOpts struct {
	Connection *sql.DB
}

func Load(opts LoadOpts) (FloorPriceDiffs, error) {
	connection := opts.Connection
	selectedColumns := []string{
		"collection_id",
		"previous_price",
		"current_price",
		"last_updated",
	}
	query := fmt.Sprintf("SELECT %s FROM floor_price_diffs", strings.Join(selectedColumns, ", "))
	log.Debug(query)
	rows, err := connection.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %s", err)
	}
	defer rows.Close()
	floorPriceDiffs := FloorPriceDiffs{}
	for rows.Next() {
		fpd := FloorPriceDiff{}
		if err := rows.Scan(
			&fpd.CollectionID,
			&fpd.PreviousPrice,
			&fpd.CurrentPrice,
			&fpd.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row[%v]", len(floorPriceDiffs))
		}
		floorPriceDiffs = append(floorPriceDiffs, fpd)
	}
	return floorPriceDiffs, nil
}
