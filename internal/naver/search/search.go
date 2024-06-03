package search

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
)

const NaverAPIBaseUrl = "https://openapi.naver.com"
const NaverLocalSearchPath = "/v1/search/local.json"

type Client interface {
	Query(ctx context.Context, queryString string) (*QueryResult, error)
}

type client struct {
	resty *resty.Client
}

type QueryResult struct {
	LastBuildDate string `json:"lastBuildDate"`
	Total         int    `json:"total"`
	Start         int    `json:"start"`
	Display       int    `json:"display"`
	Items         []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Telephone   string `json:"telephone"`
		Address     string `json:"address"`
		RoadAddress string `json:"roadAddress"`
		Mapx        string `json:"mapx"`
		Mapy        string `json:"mapy"`
	} `json:"items"`
}

func New(clientId, clientSecret string) Client {
	cli := resty.New().
		SetBaseURL(NaverAPIBaseUrl).
		SetHeaders(map[string]string{
			"X-Naver-Client-Id":     clientId,
			"X-Naver-Client-Secret": clientSecret,
		})
	return &client{cli}
}

func (c *client) Query(ctx context.Context, queryString string) (*QueryResult, error) {
	req := c.resty.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"query":   queryString,
			"display": "5",
			"sort":    "random",
		})
	res, err := req.Get(NaverLocalSearchPath)
	if err != nil {
		return nil, errors.Join(errors.New("fail to query naver local"), err)
	}
	if res.StatusCode() >= 400 {
		return nil, errors.Join(errors.New("status is not 200"), errors.New(string(res.Body())))
	}
	result := new(QueryResult)
	err = json.Unmarshal(res.Body(), result)
	if err != nil {
		return nil, errors.Join(errors.New("fail to unmarshal naver local result"), err)
	}

	return result, nil
}
