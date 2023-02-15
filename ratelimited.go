package discogs

import (
	"context"
)

// RateLimited returns d with all functions replaced with versions that honor rate limiting per rl.
func RateLimited(d Discogs, rl *RateLimit) Discogs {
	return &ratelimitedDiscogs{
		ratelimitedCollectionService:  ratelimitedCollectionService{d: d, rl: rl},
		ratelimitedDatabaseService:    ratelimitedDatabaseService{d: d, rl: rl},
		ratelimitedSearchService:      ratelimitedSearchService{d: d, rl: rl},
		ratelimitedMarketPlaceService: ratelimitedMarketPlaceService{d: d, rl: rl},
	}
}

// ratelimitedDiscogs implements Discogs with rate limiting
type ratelimitedDiscogs struct {
	ratelimitedCollectionService
	ratelimitedDatabaseService
	ratelimitedSearchService
	ratelimitedMarketPlaceService
}

type ratelimitedDatabaseService struct {
	d  Discogs
	rl *RateLimit
}

func (r ratelimitedDatabaseService) Artist(ctx context.Context, artistID int) (v *Artist, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.Artist(ctx, artistID)
		return err
	})
	return
}

func (r ratelimitedDatabaseService) ArtistReleases(ctx context.Context, artistID int, pagination *Pagination) (v *ArtistReleases, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.ArtistReleases(ctx, artistID, pagination)
		return err
	})
	return
}

func (r ratelimitedDatabaseService) Label(ctx context.Context, labelID int) (v *Label, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.Label(ctx, labelID)
		return err
	})
	return
}

func (r ratelimitedDatabaseService) LabelReleases(ctx context.Context, labelID int, pagination *Pagination) (v *LabelReleases, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.LabelReleases(ctx, labelID, pagination)
		return err
	})
	return
}

func (r ratelimitedDatabaseService) Master(ctx context.Context, masterID int) (v *Master, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.Master(ctx, masterID)
		return err
	})
	return
}

func (r ratelimitedDatabaseService) MasterVersions(ctx context.Context, masterID int, pagination *Pagination) (v *MasterVersions, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.MasterVersions(ctx, masterID, pagination)
		return err
	})
	return
}

func (r ratelimitedDatabaseService) Release(ctx context.Context, releaseID int) (v *Release, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.Release(ctx, releaseID)
		return err
	})
	return
}

func (r ratelimitedDatabaseService) ReleaseRating(ctx context.Context, releaseID int) (v *ReleaseRating, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.ReleaseRating(ctx, releaseID)
		return err
	})
	return
}

type ratelimitedMarketPlaceService struct {
	d  Discogs
	rl *RateLimit
}

func (r ratelimitedMarketPlaceService) PriceSuggestions(ctx context.Context, releaseID int) (v *PriceListing, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.PriceSuggestions(ctx, releaseID)
		return err
	})
	return
}

func (r ratelimitedMarketPlaceService) ReleaseStatistics(ctx context.Context, releaseID int) (v *Stats, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.ReleaseStatistics(ctx, releaseID)
		return err
	})
	return
}

type ratelimitedCollectionService struct {
	d  Discogs
	rl *RateLimit
}

func (r ratelimitedCollectionService) CollectionFolders(ctx context.Context, username string) (v *CollectionFolders, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.CollectionFolders(ctx, username)
		return err
	})
	return
}

func (r ratelimitedCollectionService) CollectionItemsByFolder(ctx context.Context, username string, folderID int, pagination *Pagination) (v *CollectionItems, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.CollectionItemsByFolder(ctx, username, folderID, pagination)
		return err
	})
	return
}

func (r ratelimitedCollectionService) CollectionItemsByRelease(ctx context.Context, username string, releaseID int) (v *CollectionItems, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.CollectionItemsByRelease(ctx, username, releaseID)
		return err
	})
	return
}

func (r ratelimitedCollectionService) Folder(ctx context.Context, username string, folderID int) (v *Folder, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.Folder(ctx, username, folderID)
		return err
	})
	return
}

type ratelimitedSearchService struct {
	d  Discogs
	rl *RateLimit
}

func (r ratelimitedSearchService) Search(ctx context.Context, req SearchRequest) (v *Search, e error) {
	e = r.rl.Call(ctx, func() error {
		var err error
		v, err = r.d.Search(ctx, req)
		return err
	})
	return
}
