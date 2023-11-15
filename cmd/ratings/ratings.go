package ratings

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

const (
	catalogItemRating = "/api/ratings/v1/catalogitem/"
	provisionsRatings = "/api/ratings/v1/request/"
)

type RatingClient struct {
	ratingsAPI string
}

type NewRating struct {
	RequestID string `json:"request_id"`
	Email   string `json:"email"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment,omitempty"`
	Useful  string `json:"useful,omitempty"`
}

type Rating struct {
	RatingScore  float32 `json:"rating_score"`
	TotalRatings int     `json:"total_ratings"`
}

func NewClient(ratingsAPI string) (*RatingClient, error) {
	// Check it's an Absolute URL or absolute path
	uri, err := url.ParseRequestURI(ratingsAPI)
	if err != nil {
		return nil, err
	}

	// Check it's an acceptable scheme
	switch uri.Scheme {
	case "http":
	case "https":
	default:
		return nil, errors.New("invalid scheme")
	}

	// Check it's a valid domain name
	_, err = net.LookupHost(uri.Hostname())
	if err != nil {
		return nil, err
	}

	return &RatingClient{
		ratingsAPI: strings.TrimRight(ratingsAPI, "/"),
	}, nil
}

func (c *RatingClient) GetRatings(catalogItemName string) (*Rating, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s%s", c.ratingsAPI, catalogItemRating, catalogItemName))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rating := Rating{}
	err = json.Unmarshal(bodyBytes, &rating)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (c *RatingClient) SetRating(rating NewRating) (bool, error) {
	byteRating, err := json.Marshal(rating)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(
		fmt.Sprintf(
			"%s%s%s",
			c.ratingsAPI,
			provisionsRatings,
			rating.RequestID,
		),
		"application/json",
		bytes.NewBuffer(byteRating),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}

	return true, nil
}
