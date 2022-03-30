package collection

import "github.com/zephinzer/ebzbaybot/pkg/ebzbay"

type NewOpts struct {
	Aliases string
	Label   string
}

func New(from ebzbay.Collection, opts NewOpts) Collection {
	return Collection{
		ID:               from.Collection,
		Label:            opts.Label,
		Aliases:          opts.Aliases,
		AverageSalePrice: from.AverageSalePrice,
		FloorPrice:       from.FloorPrice,
		NumberActive:     from.NumberActive,
		NumberOfSales:    from.NumberOfSales,
		TotalRoyalties:   from.TotalRoyalties,
		TotalVolume:      from.TotalVolume,
	}
}
