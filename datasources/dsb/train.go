package dsb

import (
	"encoding/json"
	"fmt"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/utils"
	"io"
	"net/url"
	"strings"
)

func (api *DSBApi) GetTrains(key string, value string) (chan *models.TrainList, chan error) {
	if len(key) == 0 && len(value) == 0 {
		return api.getTrainsByQuery("")
	} else {
		filter := fmt.Sprintf("%s eq '%s'", key, value)
		// https://github.com/golang/go/issues/4013
		escapedFilter := fmt.Sprintf("?$filter=%s", strings.Replace(url.QueryEscape(filter), "+", "%20", -1))
		return api.getTrainsByQuery(escapedFilter)
	}
}

func (api *DSBApi) getTrainsByQuery(query string) (chan *models.TrainList, chan error) {
	success, failure := make(chan *models.TrainList), make(chan error)

	request, err := api.buildRequest(httpGET, "/Queue()"+query)

	if err != nil {
		utils.SubmitAsync(err, failure)
	} else {
		q := NewQuery(failure, request, func(body io.Reader) {
			var container map[string][]json.RawMessage
			if err := json.NewDecoder(body).Decode(&container); err != nil {
				failure <- err
			} else {
				success <- convertTrainJSONList(container["d"])
				logger.Println("Completed request")
			}
		})
		api.DoAsync(q)
	}
	return success, failure
}

func convertTrainJSONList(items []json.RawMessage) *models.TrainList {
	result := []models.Train{}
	for _, item := range items {
		m := models.Train{}
		if err := json.Unmarshal(item, &m); err == nil {
			result = append(result, m)
		} else {
			logger.Println(err)
		}
	}
	return &models.TrainList{result}
}
