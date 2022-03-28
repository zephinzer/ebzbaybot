package ebzbay

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"gitlab.com/zephinzer/go-devops"
)

func getDefaultParams(page, collection, sortBy, direction string) map[string]string {
	return map[string]string{
		"state":      "0",
		"page":       page,
		"pageSize":   "10",
		"sortBy":     sortBy,
		"direction":  direction,
		"collection": collection,
	}
}

func GetCollectionStats(collection string) CollectionStats {
	params := getDefaultParams("1", collection, SortByPrice, DirectionAscending)
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
	for _, listing := range listingResponse.Listings {
		// timeNow := time.Now()
		// timeListed := time.Unix(listing.ListingTime, 0)
		// timeSinceInMinutes := timeNow.Sub(timeListed).Minutes()
		// hoursSince := int64(timeSinceInMinutes) / 60
		// minutesSince := int64(timeSinceInMinutes) - hoursSince*60
		formattedPrice, _ := strconv.ParseFloat(listing.Price, 64)
		totalPrice += formattedPrice
	}
	averageLowestPrice := totalPrice / float64(len(listingResponse.Listings))
	floorPrice, _ := strconv.ParseFloat(listingResponse.Listings[0].Price, 64)
	collectionStats := CollectionStats{
		AverageLowestTenPrice: averageLowestPrice,
		FloorPrice:            floorPrice,
		Listings:              listingResponse.TotalCount,
	}
	return collectionStats
}