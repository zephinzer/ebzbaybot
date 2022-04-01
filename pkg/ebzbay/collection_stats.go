package ebzbay

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/zephinzer/ebzbaybot/pkg/constants"
	"gitlab.com/zephinzer/go-devops"
)

type CollectionStats struct {
	FloorPrice            float64 `json:"floorPrice"`
	AverageLowestTenPrice float64 `json:"averageLowestTenPrice"`
	Listings              int64   `json:"listings"`
	FloorListingID        string  `json:"floorListingID"`
	FloorImageURL         string  `json:"floorImageURL"`
	FloorEdition          string  `json:"floorEdition"`
	FloorScore            string  `json:"floorScore"`
}

func getCollectionStatsParams(page, collection, sortBy, direction string) map[string]string {
	return map[string]string{
		"state":      "0",
		"page":       page,
		"pageSize":   "10",
		"sortBy":     sortBy,
		"direction":  direction,
		"collection": collection,
	}
}

func GetCollectionStats(collection string) *CollectionStats {
	if _, exists := constants.CollectionByAddress[collection]; !exists {
		return nil
	}
	params := getCollectionStatsParams("1", collection, SortByPrice, DirectionAscending)
	urlInstance, _ := url.Parse(API_BASE_URL)
	urlInstance.Path = API_PATH_LISTINGS
	query := urlInstance.Query()
	for paramKey, paramValue := range params {
		query.Add(paramKey, paramValue)
	}
	urlInstance.RawQuery = query.Encode()
	response, _ := devops.SendHTTPRequest(devops.SendHTTPRequestOpts{
		Headers: map[string][]string{"Referer": {REFERER}},
		URL:     urlInstance,
		Method:  http.MethodGet,
	})
	body, _ := ioutil.ReadAll(response.Body)
	var listingResponse ListingResponse
	json.Unmarshal(body, &listingResponse)
	var totalPrice float64

	// process floor item details
	var floorListingID string
	var floorImageURL string
	var floorEdition string
	var floorScore string
	if listingResponse.Listings != nil && len(listingResponse.Listings) > 0 {
		floorListing := listingResponse.Listings[0]
		floorListingID = strconv.FormatInt(floorListing.ListingID, 10)
		floorImageURL = floorListing.NFT.OriginalImage
		floorEdition = strconv.FormatInt(floorListing.NFT.Edition, 10)
		floorScore = strconv.FormatFloat(floorListing.NFT.Score, 'f', 2, 64)
	}
	for _, listing := range listingResponse.Listings {
		formattedPrice, _ := strconv.ParseFloat(listing.Price, 64)
		totalPrice += formattedPrice
	}
	averageLowestPrice := totalPrice / float64(len(listingResponse.Listings))
	floorPrice, _ := strconv.ParseFloat(listingResponse.Listings[0].Price, 64)
	collectionStats := CollectionStats{
		AverageLowestTenPrice: averageLowestPrice,
		FloorPrice:            floorPrice,
		Listings:              listingResponse.TotalCount,
		FloorListingID:        floorListingID,
		FloorImageURL:         floorImageURL,
		FloorEdition:          floorEdition,
		FloorScore:            floorScore,
	}
	return &collectionStats
}
