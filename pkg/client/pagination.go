package client

import (
	"net/url"
	"strconv"
)

// By default, the number of objects returned per page is 100. The maximum number of object supported per page.
// It can be adjusted by adding the 'per_page' parameter in the query string.

const ItemsPerPage = 100

// PageOptions is options for list method of paginatable resources.
// It's used to create query string.
type PageOptions struct {
	PerPage int `url:"limit,omitempty"`
	Page    int `url:"page,omitempty"`
}

type ReqOpt func(reqURL *url.URL)

// WithPageLimit : Number of items to return.
func WithPageLimit(pageLimit int) ReqOpt {
	if pageLimit <= 0 || pageLimit > ItemsPerPage {
		pageLimit = ItemsPerPage
	}
	return WithQueryParam("perPage", strconv.Itoa(pageLimit))
}

// WithPage : Number for the page (inclusive). The page number starts with 1.
// If page is 0, first page is assumed.
func WithPage(page int) ReqOpt {
	if page == 0 {
		page = 1
	}
	return WithQueryParam("page", strconv.Itoa(page))
}

func WithQueryParam(key string, value string) ReqOpt {
	return func(reqURL *url.URL) {
		q := reqURL.Query()
		q.Set(key, value)
		reqURL.RawQuery = q.Encode()
	}
}
