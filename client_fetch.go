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

	return result, nil
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
	fmt.Fprintf(os.Stderr, "fetch:%s\n", u.String())
	return u, err
}

func parseFetchedContent(r io.Reader) (*FetchResult, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	abst := doc.Find("#novel_ex")
	if abst.Size() == 1 {
		res, err := parseSeriesIndexPage(doc)
		return res, err
	}
	res, err := parseShortContentPage(doc)
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
	res.Pages[0] = FetchPage{}
	res.Pages[0].Preface = parseContentLines(doc, "#novel_p")
	res.Pages[0].Lines = parseContentLines(doc, "#novel_honbun")
	res.Pages[0].Afterword = parseContentLines(doc, "#novel_a")
	raw, err := doc.Find("#novel_color").First().Html()
	if err != nil {
		fmt.Fprintf(os.Stderr, "raw content html error:%s", err)
	}
	res.Pages[0].rawHTML = raw
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
