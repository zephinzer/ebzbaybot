package ebzbay

const (
	API_BASE_URL         = "https://api.ebisusbay.com"
	API_PATH_LISTINGS    = "/listings"
	API_PATH_COLLECTIONS = "/collections"
	REFERER              = "https://app.ebisusbay.com/"

	CollectionMadMeerkays = "0x89dbc8bd9a6037cbd6ec66c4bf4189c9747b1c56"

	SortByPrice         = "price"
	SortByRank          = "rank"
	SortByScore         = "score"
	DirectionAscending  = "asc"
	DirectionDescending = "desc"
)

type CollectionStats struct {
	FloorPrice            float64
	AverageLowestTenPrice float64
	Listings              int64
}
