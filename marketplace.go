package discogs

import (
	"context"
	"net/url"
	"strconv"
)

const (
	priceSuggestionsURI = "/price_suggestions/"
	releaseStatsURI     = "/stats/"
)

type marketPlaceService struct {
	request  requestFunc
	url      string
	currency string
}

type MarketPlaceService interface {
	// The best price suggestions according to grading
	// Authentication is required.
	PriceSuggestions(ctx context.Context, releaseID int) (*PriceListing, error)
	// Short summary of marketplace listings
	// Authentication is optional.
	ReleaseStatistics(ctx context.Context, releaseID int) (*Stats, error)
}

func newMarketPlaceService(req requestFunc, url string, currency string) MarketPlaceService {
	return &marketPlaceService{
		request:  req,
		url:      url,
		currency: currency,
	}
}

// Listing is a marketplace listing with the user's currency and a price value
type Listing struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

// PriceListings are Listings per grading quality
type PriceListing struct {
	VeryGood     *Listing `json:"Very Good (VG),omitempty"`
	GoodPlus     *Listing `json:"Good Plus (G+),omitempty"`
	NearMint     *Listing `json:"Near Mint (NM or M-)"`
	Good         *Listing `json:"Good (G),omitempty"`
	VeryGoodPlus *Listing `json:"Very Good Plus (VG+),omitempty"`
	Mint         *Listing `json:"Mint (M),omitempty"`
	Fair         *Listing `json:"Fair (F),omitempty"`
	Poor         *Listing `json:"Poor (P),omitempty"`
}

// Stats returns the marketplace stats summary for a release containing
type Stats struct {
	LowestPrice *Listing `json:"lowest_price"`
	ForSale     int      `json:"num_for_sale"`
	Blocked     bool     `json:"blocked_from_sale"`
}

func (s *marketPlaceService) ReleaseStatistics(ctx context.Context, releaseID int) (*Stats, error) {
	params := url.Values{}
	params.Set("curr_abbr", s.currency)

	var stats *Stats
	err := s.request(ctx, s.url+releaseStatsURI+strconv.Itoa(releaseID), params, &stats)
	return stats, err
}

func (s *marketPlaceService) PriceSuggestions(ctx context.Context, releaseID int) (*PriceListing, error) {
	var listings *PriceListing
	err := s.request(ctx, s.url+priceSuggestionsURI+strconv.Itoa(releaseID), nil, &listings)
	return listings, err
}
