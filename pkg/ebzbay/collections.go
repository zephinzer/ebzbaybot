package ebzbay

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"gitlab.com/zephinzer/go-devops"
)

type CollectionsResponse struct {
	Status      int         `json:"status"`
	Collections Collections `json:"collections"`
}

type Collections []Collection

type Collection struct {
	AverageSalePrice string `json:"averageSalePrice"`
	Collection       string `json:"collection"`
	FloorPrice       string `json:"floorPrice"`
	NumberActive     string `json:"numberActive"`
	NumberOfSales    string `json:"numberOfSales"`
	TotalRoyalties   string `json:"totalRoyalties"`
	TotalVolume      string `json:"totalVolume"`
}

func getCollectionsParams(sortBy, direction string) map[string]string {
	return map[string]string{
		"sortBy":    sortBy,
		"direction": direction,
	}
}

func GetCollections() (Collections, error) {
	params := getCollectionsParams(SortByVolume, DirectionDescending)
	urlInstance, _ := url.Parse(API_BASE_URL)
	urlInstance.Path = API_PATH_COLLECTIONS
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
	var collectionsResponse CollectionsResponse
	if err := json.Unmarshal(body, &collectionsResponse); err != nil {
		fmt.Println(string(body))
		return nil, fmt.Errorf("failed to understand response for /collections: %s", err)
	}
	if collectionsResponse.Status != http.StatusOK {
		return nil, fmt.Errorf("failed to receive a success response (status: %v)", collectionsResponse.Status)
	}
	return collectionsResponse.Collections, nil
}
