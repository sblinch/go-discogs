package discogs

import (
	"context"
	"strconv"
)

// CollectionService is an interface to work with collection.
type CollectionService interface {
	// Retrieve a list of folders in a user’s collection.
	// If folder_id is not 0, authentication as the collection owner is required.
	CollectionFolders(ctx context.Context, username string) (*CollectionFolders, error)
	// Retrieve a list of items in a folder in a user’s collection.
	// If folderID is not 0, authentication with token is required.
	CollectionItemsByFolder(ctx context.Context, username string, folderID int, pagination *Pagination) (*CollectionItems, error)
	// Retrieve the user’s collection folders which contain a specified release.
	// The releaseID must be non-zero.
	CollectionItemsByRelease(ctx context.Context, username string, releaseID int) (*CollectionItems, error)
	// Retrieve metadata about a folder in a user’s collection.
	Folder(ctx context.Context, username string, folderID int) (*Folder, error)
}

type collectionService struct {
	request requestFunc
	url     string
}

func newCollectionService(req requestFunc, url string) CollectionService {
	return &collectionService{
		request: req,
		url:     url,
	}
}

// Folder serves folder response from discogs.
type Folder struct {
	ID          int    `json:"id"`
	Count       int    `json:"count"`
	Name        string `json:"name"`
	ResourceURL string `json:"resource_url"`
}

func (s *collectionService) Folder(ctx context.Context, username string, folderID int) (*Folder, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}
	var folder *Folder
	err := s.request(ctx, s.url+"/"+username+"/collection/folders/"+strconv.Itoa(folderID), nil, &folder)
	return folder, err
}

// CollectionFolders serves collection response from discogs.
type CollectionFolders struct {
	Folders []Folder `json:"folders"`
}

func (s *collectionService) CollectionFolders(ctx context.Context, username string) (*CollectionFolders, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}
	var collection *CollectionFolders
	err := s.request(ctx, s.url+"/"+username+"/collection/folders", nil, &collection)
	return collection, err
}

// CollectionItemSource ...
type CollectionItemSource struct {
	ID               int              `json:"id"`
	BasicInformation BasicInformation `json:"basic_information"`
	DateAdded        string           `json:"date_added"`
	FolderID         int              `json:"folder_id,omitempty"`
	InstanceID       int              `json:"instance_id"`
	Notes            []Notes          `json:"notes,omitempty"`
	Rating           int              `json:"rating"`
}

// BasicInformation ...
type BasicInformation struct {
	ID          int            `json:"id"`
	Artists     []ArtistSource `json:"artists"`
	CoverImage  string         `json:"cover_image"`
	Formats     []Format       `json:"formats"`
	Labels      []LabelSource  `json:"labels"`
	Genres      []string       `json:"genres"`
	MasterID    int            `json:"master_id"`
	MasterURL   *string        `json:"master_url"`
	ResourceURL string         `json:"resource_url"`
	Styles      []string       `json:"styles"`
	Thumb       string         `json:"thumb"`
	Title       string         `json:"title"`
	Year        int            `json:"year"`
}

// CollectionItems list of items in a user’s collection
type CollectionItems struct {
	Pagination Page                   `json:"pagination"`
	Items      []CollectionItemSource `json:"releases"`
}

// valid sort keys
// https://www.discogs.com/developers#page:user-collection,header:user-collection-collection-items-by-folder
var validItemsByFolderSort = map[string]struct{}{
	"":       struct{}{},
	"label":  struct{}{},
	"artist": struct{}{},
	"title":  struct{}{},
	"catno":  struct{}{},
	"format": struct{}{},
	"rating": struct{}{},
	"added":  struct{}{},
	"year":   struct{}{},
}

func (s *collectionService) CollectionItemsByFolder(ctx context.Context, username string, folderID int, pagination *Pagination) (*CollectionItems, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}
	if pagination != nil {
		if _, ok := validItemsByFolderSort[pagination.Sort]; !ok {
			return nil, ErrInvalidSortKey
		}
	}
	var items *CollectionItems
	err := s.request(ctx, s.url+"/"+username+"/collection/folders/"+strconv.Itoa(folderID)+"/releases", pagination.params(), &items)
	return items, err
}

func (s *collectionService) CollectionItemsByRelease(ctx context.Context, username string, releaseID int) (*CollectionItems, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}
	if releaseID == 0 {
		return nil, ErrInvalidReleaseID
	}
	var items *CollectionItems
	err := s.request(ctx, s.url+"/"+username+"/collection/releases/"+strconv.Itoa(releaseID), nil, &items)
	return items, err
}
