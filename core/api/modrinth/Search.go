package modrinth

import (
	"encoding/json"
	"strconv"

	"github.com/NEXORA-Studios/Nova.ModDeps/core/api"
)

func Search(query string, offset int) (IMrSearchResponse, error) {
	requester := api.Requester{}
	queryMap := map[string]string{
		"query":  query,
		"offset": strconv.Itoa(offset),
		"limit":  strconv.Itoa(10),
		"facets": "[[\"project_type:mod\"]]",
	}
	response, err := requester.Get("/search", queryMap)
	if err != nil {
		return IMrSearchResponse{}, err
	}
	var searchResponse IMrSearchResponse
	err = json.Unmarshal([]byte(response), &searchResponse)
	if err != nil {
		return IMrSearchResponse{}, err
	}
	return searchResponse, nil
}
