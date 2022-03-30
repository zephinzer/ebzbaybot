package collection

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type LoadOpts struct {
	Address    string
	Connection *sql.DB
}

func LoadAll(opts LoadOpts) (Collections, error) {
	connection := opts.Connection
	selectedColumns := []string{
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
	query := fmt.Sprintf("SELECT %s FROM collections", strings.Join(selectedColumns, ", "))

	log.Debug(query)
	rows, err := connection.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %s", err)
	}
	defer rows.Close()
	collections := Collections{}
	for rows.Next() {
		c := Collection{}
		if err := rows.Scan(
			&c.ID,
			&c.Label,
			&c.Aliases,
			&c.AverageSalePrice,
			&c.FloorPrice,
			&c.NumberActive,
			&c.NumberOfSales,
			&c.TotalRoyalties,
			&c.TotalVolume,
			&c.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row[%v]", len(collections))
		}
		collections = append(collections, c)
	}
	return collections, nil
}
