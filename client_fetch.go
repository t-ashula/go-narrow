package narrow

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Fetch download novel contents
func (c *Client) Fetch(ctx context.Context, params *FetchParams) (*FetchResult, error) {
	contentURL, err := params.toContentURL()
	if err != nil {
		return nil, err
	}

	if params.Site != FetchSiteNarou && params.AllowOver18 {
		if c.httpClient.Jar == nil {
			jar, err := cookiejar.New(nil)
			if err != nil {
				return nil, err
			}
			c.httpClient.Jar = jar
		}

		over18 := http.Cookie{Name: "over18", Value: "yes", Expires: time.Now().AddDate(1, 0, 0), Domain: ".syosetu.com"}
		c.httpClient.Jar.SetCookies(contentURL, []*http.Cookie{&over18})
	}

	req, err := http.NewRequest("GET", contentURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("user-agent", userAgent)

	req = req.WithContext(ctx)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result, err := parseFetchedContent(res.Body)
	if err != nil {
		return nil, err
	}
	result.Site = params.Site
	result.NCode = params.NCode
	if result.NovelType == 1 {
		err = nil
		if params.WithContent {
			err = c.fetchAllPageContent(ctx, result, params)
		} else if params.Page > 0 {
			err = c.fetchSinglePageContent(ctx, result, params)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch page content failed. %v", err)
			return result, nil
		}
	}

	return result, nil
}

func (c *Client) fetchAllPageContent(ctx context.Context, result *FetchResult, params *FetchParams) error {
	for i := 1; i <= result.PageCount; i++ {
		page, err := c.fetchPageContent(ctx, params, i)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch page %d content failed, %v", i, err)
			return err
		}
		// Pages[i].PublishDate, LastUpdateDate is already set
		page.PublishDate = result.Pages[i-1].PublishDate
		page.LastUpdateDate = result.Pages[i-1].LastUpdateDate
		result.Pages[i-1] = *page
	}
	return nil
}

func (c *Client) fetchSinglePageContent(ctx context.Context, result *FetchResult, params *FetchParams) error {
	if params.Page > result.PageCount {
		// do nothing
		fmt.Fprintf(os.Stderr, "specified page %d greater than fetched index %d", params.Page, result.PageCount)
		return nil
	}
	pageNo := params.Page
	page, err := c.fetchPageContent(ctx, params, pageNo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch page %d failed %s", pageNo, err)
		return err
	}

	// Pages[i].PublishDate, LastUpdateDate is already set
	page.PublishDate = result.Pages[pageNo-1].PublishDate
	page.LastUpdateDate = result.Pages[pageNo-1].LastUpdateDate
	result.Pages[pageNo-1] = *page

	return nil
}

func (c *Client) fetchPageContent(ctx context.Context, params *FetchParams, pageNo int) (*FetchPage, error) {
	u, err := params.toContentURL()
	if err != nil {
		return nil, err
	}
	u.Path = fmt.Sprintf("%s%d/", u.Path, pageNo)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("user-agent", userAgent)

	req = req.WithContext(ctx)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	page, err := parseContentPage(doc)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// NewFetchParams returns new FetchParams
func NewFetchParams() *FetchParams {
	params := &FetchParams{}
	return params
}

func (params *FetchParams) toContentURL() (*url.URL, error) {
	subDomain := "ncode"
	if params.Site != FetchSiteNarou {
		subDomain = "novel18"
	}

	u, err := url.Parse(fmt.Sprintf("https://%s.syosetu.com/%s/", subDomain, params.NCode))
	return u, err
}

func parseFetchedContent(r io.Reader) (*FetchResult, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	abst := doc.Find("#novel_ex")
	if abst.Size() == 0 {
		res, err := parseShortContentPage(doc)
		return res, err
	}

	res, err := parseSeriesIndexPage(doc)
	return res, err
}

func parseSeriesIndexPage(doc *goquery.Document) (*FetchResult, error) {
	res := &FetchResult{}
	res.Title = strings.TrimSpace(doc.Find("title").First().Text())
	res.WriterName = strings.TrimSpace(doc.Find("div.novel_writername").First().Text())
	res.NovelType = 1
	res.Abstruct = doc.Find("#novel_ex").Text()
	pages := doc.Find("div.index_box > dl.novel_sublist2")
	res.PageCount = pages.Size()
	res.Pages = make([]FetchPage, res.PageCount)
	pages.Each(func(i int, s *goquery.Selection) {
		res.Pages[i] = FetchPage{}
		subTitle := s.Find("dd.subtitle a").First().Text()
		res.Pages[i].SubTitle = strings.TrimSpace(subTitle)

		pubDateStr := s.Find("dt.long_update").First().Text()
		pubDateStr = kaiRe.ReplaceAllString(pubDateStr, "")
		pubDateStr = strings.TrimSpace(pubDateStr)
		pubDate, err := asJST(pubDateStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pubdate parse failed %v\n", err)
		} else {
			res.Pages[i].PublishDate = pubDate
		}

		kai := s.Find("dt.long_update span").First()
		if kai.Size() == 1 {
			if upd, ok := kai.Attr("title"); ok {
				upd = kaikouRe.ReplaceAllString(upd, "")
				upd = strings.TrimSpace(upd)
				updateDate, err := asJST(upd)
				if err != nil {
					fmt.Fprintf(os.Stderr, "update date parse failed %v\n", err)
				} else {
					res.Pages[i].LastUpdateDate = &updateDate
				}
			}
		}
	})
	return res, nil
}

func asJST(str string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation(novelUpdateTimeFormat, str, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

const novelUpdateTimeFormat = "2006/01/02 15:04"

var kaiRe = regexp.MustCompile(`\s*（改）\s*`)
var kaikouRe = regexp.MustCompile(`\s*改稿\s*`)

func parseShortContentPage(doc *goquery.Document) (*FetchResult, error) {
	res := &FetchResult{}
	res.Title = doc.Find("title").First().Text()
	res.WriterName = strings.TrimSpace(doc.Find("div.novel_writername").First().Text())
	res.NovelType = 2
	res.PageCount = 1
	res.Pages = make([]FetchPage, 1)
	page, err := parseContentPage(doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse content page error:%s", err)
	}
	res.Pages[0] = *page
	return res, nil
}

func parseContentLines(doc *goquery.Document, selector string) []ContentLine {
	honbun := doc.Find(selector).First()
	if honbun.Size() == 0 {
		return nil
	}
	lines := honbun.Find("p")
	content := make([]ContentLine, lines.Size())
	lines.Each(func(i int, s *goquery.Selection) {
		raw, err := goquery.OuterHtml(s) // want p#L
		if err != nil {
			fmt.Fprintf(os.Stderr, "raw line error:%s", err)
		}
		content[i] = ContentLine{RawLine: raw}
	})
	return content
}

func parseContentPage(doc *goquery.Document) (*FetchPage, error) {
	page := &FetchPage{}
	page.SubTitle = strings.TrimSpace(doc.Find("div.novel_subtitle").First().Text())
	chapter := doc.Find(".chapter_title")
	if chapter.Size() != 0 {
		title := strings.TrimSpace(chapter.First().Text())
		page.ChapterTitle = &title
	}
	page.Preface = parseContentLines(doc, "#novel_p")
	page.Lines = parseContentLines(doc, "#novel_honbun")
	page.Afterword = parseContentLines(doc, "#novel_a")
	raw, err := doc.Find("#novel_color").First().Html()
	if err != nil {
		fmt.Fprintf(os.Stderr, "raw content html error:%s", err)
	}
	page.rawHTML = raw
	return page, nil
}
