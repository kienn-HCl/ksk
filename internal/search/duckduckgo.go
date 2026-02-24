package search

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var ddgClient *http.Client

func init() {
	jar, _ := cookiejar.New(nil)
	ddgClient = &http.Client{Jar: jar}
}

type DuckDuckGo struct {
	Region string // e.g. "jp", "us", "de"
}

const ddgEndpoint = "https://html.duckduckgo.com/html/"

var ddgHeaders = http.Header{
	"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0"},
	"Content-Type":    {"application/x-www-form-urlencoded"},
	"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
	"Accept-Language": {"en-US,en;q=0.5"},
	"Referer":         {"https://html.duckduckgo.com/"},
}

func (d *DuckDuckGo) Name() string { return "duckduckgo" }

func (d *DuckDuckGo) Search(query string) (*Page, error) {
	form := url.Values{"q": {query}}
	if kl := ddgRegion(d.Region); kl != "" {
		form.Set("kl", kl)
	}
	return ddgDoSearch(form, 1)
}

func (d *DuckDuckGo) NextPage(prev *Page, query string) (*Page, error) {
	if !prev.HasMore || prev.NextParams == nil {
		return nil, fmt.Errorf("no more pages")
	}
	form := url.Values{}
	for k, v := range prev.NextParams {
		form[k] = v
	}
	form.Set("q", query)
	return ddgDoSearch(form, prev.PageNum+1)
}

func (d *DuckDuckGo) PrevPage(query string, pageNum int) (*Page, error) {
	if pageNum <= 1 {
		return d.Search(query)
	}
	page, err := d.Search(query)
	if err != nil {
		return nil, err
	}
	for i := 1; i < pageNum; i++ {
		if !page.HasMore {
			return page, nil
		}
		page, err = d.NextPage(page, query)
		if err != nil {
			return nil, err
		}
	}
	return page, nil
}

func ddgDoSearch(form url.Values, pageNum int) (*Page, error) {
	req, err := http.NewRequest("POST", ddgEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header = ddgHeaders.Clone()

	resp, err := ddgClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing search: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	if doc.Find("form#challenge-form").Length() > 0 {
		return nil, fmt.Errorf("bot detection triggered — wait a moment and try again")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search returned status %d", resp.StatusCode)
	}

	page := &Page{PageNum: pageNum}

	doc.Find(".result.results_links").Each(func(i int, s *goquery.Selection) {
		titleEl := s.Find(".result__title a.result__a")
		title := strings.TrimSpace(titleEl.Text())

		href, _ := titleEl.Attr("href")
		href = ddgCleanURL(href)

		snippet := strings.TrimSpace(s.Find(".result__snippet").Text())

		if title != "" && href != "" {
			page.Results = append(page.Results, Result{
				Title:   title,
				URL:     href,
				Snippet: snippet,
			})
		}
	})

	navForm := doc.Find(".nav-link form")
	if navForm.Length() > 0 {
		page.HasMore = true
		page.NextParams = url.Values{}
		navForm.First().Find("input[type='hidden']").Each(func(i int, s *goquery.Selection) {
			name, _ := s.Attr("name")
			val, _ := s.Attr("value")
			if name != "" {
				page.NextParams.Set(name, val)
			}
		})
	}

	return page, nil
}

func ddgCleanURL(rawURL string) string {
	if strings.HasPrefix(rawURL, "//duckduckgo.com/l/?uddg=") {
		parsed, err := url.Parse("https:" + rawURL)
		if err == nil {
			if uddg := parsed.Query().Get("uddg"); uddg != "" {
				return uddg
			}
		}
	}
	return rawURL
}

// ddgRegion maps a short region code to DuckDuckGo's kl parameter.
var ddgRegionMap = map[string]string{
	"jp": "jp-jp",
	"us": "us-en",
	"uk": "uk-en",
	"de": "de-de",
	"fr": "fr-fr",
	"es": "es-es",
	"it": "it-it",
	"br": "br-pt",
	"ca": "ca-en",
	"au": "au-en",
	"in": "in-en",
	"kr": "kr-kr",
	"cn": "cn-zh",
	"tw": "tw-tzh",
	"ru": "ru-ru",
}

func ddgRegion(region string) string {
	if region == "" {
		return ""
	}
	if kl, ok := ddgRegionMap[region]; ok {
		return kl
	}
	// Fallback: try "xx-xx" pattern
	return region + "-" + region
}
