package ebzbay

type ListingResponse struct {
	Listings   Listings `json:"listings"`
	TotalCount int64    `json:"totalCount"`
}

type Listings []Listing

type Listing struct {
	ListingID   int64  `json:"listingId`
	ListingTime int64  `json:"listingTime"`
	NFT         NFT    `json:"nft"`
	NftAddress  string `json:"nftAddress"`
	Price       string `json:"price"`
	Seller      string `json:"seller"`
}

type NFT struct {
	Name                 string            `json:"name"`
	Edition              int64             `json:"edition"`
	NftAddress           string            `json:"nftAddress"`
	NftID                string            `json:"nftId"`
	OriginalImage        string            `json:"original_image"`
	Score                float64           `json:"score"`
	SimplifiedAttributes map[string]string `json:"simplifiedAttributes"`
}
