package search

import "net/url"

type Result struct {
	Title   string
	URL     string
	Snippet string
}

type Page struct {
	Results    []Result
	NextParams url.Values
	PageNum    int
	HasMore    bool
}

type Backend interface {
	Search(query string) (*Page, error)
	NextPage(prev *Page, query string) (*Page, error)
	PrevPage(query string, pageNum int) (*Page, error)
	Name() string
}
