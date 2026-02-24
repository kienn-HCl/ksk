package search

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Brave struct {
	Region string // e.g. "jp", "us", "de"
}

const braveEndpoint = "https://search.brave.com/search"

var braveHeaders = http.Header{
	"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64; rv:133.0) Gecko/20100101 Firefox/133.0"},
	"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
	"Accept-Language": {"en-US,en;q=0.5"},
}

func (b *Brave) Name() string { return "brave" }

func (b *Brave) Search(query string) (*Page, error) {
	return braveDoSearch(query, 0, 1, b.Region)
}

func (b *Brave) NextPage(prev *Page, query string) (*Page, error) {
	if !prev.HasMore {
		return nil, fmt.Errorf("no more pages")
	}
	// Brave's offset parameter is a 0-indexed page number
	return braveDoSearch(query, prev.PageNum, prev.PageNum+1, b.Region)
}

func (b *Brave) PrevPage(query string, pageNum int) (*Page, error) {
	if pageNum <= 1 {
		return b.Search(query)
	}
	return braveDoSearch(query, pageNum-1, pageNum, b.Region)
}

func braveDoSearch(query string, offset, pageNum int, region string) (*Page, error) {
	params := url.Values{
		"q":      {query},
		"source": {"web"},
	}
	if region != "" {
		params.Set("country", braveRegion(region))
	}
	if offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", offset))
	}

	req, err := http.NewRequest("GET", braveEndpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header = braveHeaders.Clone()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == 503 {
		return nil, fmt.Errorf("rate limit triggered (status %d) — try again later", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	page := &Page{PageNum: pageNum}

	// Main results: div.snippet[data-type="web"]
	doc.Find(`div.snippet[data-type="web"]`).Each(func(i int, s *goquery.Selection) {
		braveExtractResult(s, page)
	})

	// Fallback: any div.snippet with data-pos
	if len(page.Results) == 0 {
		doc.Find("div.snippet[data-pos]").Each(func(i int, s *goquery.Selection) {
			braveExtractResult(s, page)
		})
	}

	// Detect next page: Brave uses pagination links with offset parameter
	hasNextLink := false
	doc.Find("a[href]").EachWithBreak(func(i int, a *goquery.Selection) bool {
		href, _ := a.Attr("href")
		if strings.Contains(href, "offset=") {
			hasNextLink = true
			return false
		}
		return true
	})
	if hasNextLink || len(page.Results) >= 10 {
		page.HasMore = true
	}

	return page, nil
}

func braveExtractResult(s *goquery.Selection, page *Page) {
	// Title: div.search-snippet-title or a.title text
	var title string
	titleEl := s.Find("div.search-snippet-title").First()
	if titleEl.Length() > 0 {
		title = strings.TrimSpace(titleEl.Text())
	}
	if title == "" {
		titleEl = s.Find("a.title").First()
		if titleEl.Length() > 0 {
			title = strings.TrimSpace(titleEl.Text())
		}
	}
	if title == "" {
		return
	}

	// URL: first <a> with an external href
	var href string
	s.Find("a[href]").EachWithBreak(func(i int, a *goquery.Selection) bool {
		h, exists := a.Attr("href")
		if exists && strings.HasPrefix(h, "http") && !strings.Contains(h, "brave.com") {
			href = h
			return false
		}
		return true
	})
	if href == "" {
		return
	}

	// Snippet: div.generic-snippet .content, or div.description
	var snippet string
	for _, sel := range []string{
		"div.generic-snippet .content",
		"div.description",
		"p.snippet-description",
	} {
		el := s.Find(sel).First()
		if el.Length() > 0 {
			snippet = strings.TrimSpace(el.Text())
			break
		}
	}

	page.Results = append(page.Results, Result{
		Title:   title,
		URL:     href,
		Snippet: snippet,
	})
}

// braveRegion maps a short region code to Brave's country parameter.
var braveRegionMap = map[string]string{
	"jp": "jp",
	"us": "us",
	"uk": "gb",
	"de": "de",
	"fr": "fr",
	"es": "es",
	"it": "it",
	"br": "br",
	"ca": "ca",
	"au": "au",
	"in": "in",
	"kr": "kr",
	"cn": "cn",
	"tw": "tw",
	"ru": "ru",
}

func braveRegion(region string) string {
	if r, ok := braveRegionMap[region]; ok {
		return r
	}
	return region
}
