package discogs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	discogsAPI = "https://api.discogs.com"
)

// Options is a set of options to use discogs API client
type Options struct {
	// Discogs API endpoint (optional).
	URL string
	// Currency to use (optional, default is USD).
	Currency string
	// UserAgent to to call discogs api with.
	UserAgent string
	// Token provided by discogs (optional).
	Token string
	// HTTP client instance to use for HTTP requests
	Client *http.Client
	// Rate limit instance to track request rates
	RateLimit *RateLimit
}

// Discogs is an interface for making Discogs API requests.
type Discogs interface {
	CollectionService
	DatabaseService
	MarketPlaceService
	SearchService
}

type discogs struct {
	CollectionService
	DatabaseService
	SearchService
	MarketPlaceService
}

type requestFunc func(ctx context.Context, path string, params url.Values, resp interface{}) error

// New returns a new discogs API client.
func New(o *Options) (Discogs, error) {
	header := &http.Header{}

	if o == nil || o.UserAgent == "" {
		return nil, ErrUserAgentInvalid
	}

	header.Add("User-Agent", o.UserAgent)

	cur, err := currency(o.Currency)
	if err != nil {
		return nil, err
	}

	// set token, it's required for some queries like search
	if o.Token != "" {
		header.Add("Authorization", "Discogs token="+o.Token)
	}

	if o.URL == "" {
		o.URL = discogsAPI
	}

	client := o.Client
	if client == nil {
		client = &http.Client{}
	}
	req := func(ctx context.Context, path string, params url.Values, resp interface{}) error {
		return request(ctx, client, header, o.RateLimit, path, params, resp)
	}

	return discogs{
		newCollectionService(req, o.URL+"/users"),
		newDatabaseService(req, o.URL, cur),
		newSearchService(req, o.URL+"/database/search"),
		newMarketPlaceService(req, o.URL+"/marketplace", cur),
	}, nil
}

// currency validates currency for marketplace data.
// Defaults to the authenticated users currency. Must be one of the following:
// USD GBP EUR CAD AUD JPY CHF MXN BRL NZD SEK ZAR
func currency(c string) (string, error) {
	switch c {
	case "USD", "GBP", "EUR", "CAD", "AUD", "JPY", "CHF", "MXN", "BRL", "NZD", "SEK", "ZAR":
		return c, nil
	case "":
		return "USD", nil
	default:
		return "", ErrCurrencyNotSupported
	}
}

func request(ctx context.Context, client *http.Client, header *http.Header, rl *RateLimit, path string, params url.Values, resp interface{}) error {
	r, err := http.NewRequestWithContext(ctx, "GET", path+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}
	r.Header = *header

	response, err := client.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if rl != nil {
		total, _ := strconv.Atoi(response.Header.Get("X-Discogs-Ratelimit"))               // The total number of requests you can make in a one minute window.
		used, _ := strconv.Atoi(response.Header.Get("X-Discogs-Ratelimit-Used"))           // The number of requests you’ve made in your existing rate limit window.
		remaining, _ := strconv.Atoi(response.Header.Get("X-Discogs-Ratelimit-Remaining")) // The number of remaining requests you are able to make in the existing rate limit window.
		rl.Update(total, used, remaining)
	}

	if response.StatusCode != http.StatusOK {
		switch response.StatusCode {
		case http.StatusUnauthorized:
			return ErrUnauthorized
		case http.StatusTooManyRequests:
			return ErrTooManyRequests
		default:
			return fmt.Errorf("unknown error: %s", response.Status)
		}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &resp)
}
