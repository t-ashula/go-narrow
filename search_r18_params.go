package narrow

import (
	"fmt"
	"net/url"
	"strings"
)

// NewSearchR18Params return new search r18 api parameter object
func NewSearchR18Params() *SearchR18Params {
	params := &SearchR18Params{}
	params.SearchParams = *NewSearchParams()
	return params
}

// ToURL return full URL or nil if params contains invalid condition
func (params *SearchR18Params) ToURL() (*url.URL, error) {
	return params.makeFullURL(params.endPointURL(), params.toQueryFuncs())
}

func (params *SearchR18Params) endPointURL() string {
	return NarouR18APIEndPoint
}

func (params *SearchR18Params) toQueryFuncs() []toQueryFunc {
	return []toQueryFunc{
		params.queryFromStart,
		params.queryFromLimit,
		params.queryFromOutputField,
		params.queryFromWord,
		params.queryFromNotWord,
		params.queryFromSearchField,
		params.queryFromNocGenre,
		params.queryFromNotNocGenre,
		params.queryFromUserID,
		params.queryFromRequiredKeywords,
		params.queryFromLength,
		params.queryFromKaiwaritu,
		params.queryFromSasie,
		params.queryFromReadTime,
		params.queryFromNCode,
		params.queryFromNovelState,
		params.queryFromBuntai,
		params.queryFromStop,
		params.queryFromPickup,
		params.queryFromLastUp,
	}
}

// NocGenres returns `NocGenre` parameter
func (params *SearchR18Params) NocGenres() []NocGenre { return params.nocGenres }

// AddNocGenres add NocGenre settings
func (params *SearchR18Params) AddNocGenres(genres []NocGenre) {
	ngs := make(map[NocGenre]int)
	i := 0
	for _, g := range params.nocGenres {
		ngs[g] = i
		i++
	}
	all := false
	for _, g := range genres {
		if g == NocGenreAll {
			all = true
			break
		}
		if _, has := ngs[g]; !has {
			ngs[g] = i
			i++
		}
	}

	if all {
		params.nocGenres = nil
		return
	}
	params.nocGenres = make([]NocGenre, len(ngs))
	for g, idx := range ngs {
		params.nocGenres[idx] = g
	}
}

// ClearNocGenres clear not word fields setting
func (params *SearchR18Params) ClearNocGenres() { params.nocGenres = nil }

func (params *SearchR18Params) queryFromNocGenre() url.Values {
	vs := make(url.Values)
	l := len(params.nocGenres)
	if l == 0 {
		return vs
	}

	codes := make([]string, l)
	for i, g := range params.nocGenres {
		codes[i] = fmt.Sprintf("%d", g)
	}
	vs.Set(keyNocGenre, strings.Join(codes, "-"))
	return vs
}

// NotNocGenres returns `notnocgenre` parameter
func (params *SearchR18Params) NotNocGenres() []NocGenre { return params.notNocGenres }

// AddNotNocGenres add not biggenres setting
func (params *SearchR18Params) AddNotNocGenres(genres []NocGenre) {
	ngs := make(map[NocGenre]int)
	i := 0
	for _, g := range params.notNocGenres {
		ngs[g] = i
		i++
	}
	all := false
	for _, g := range genres {
		if g == NocGenreAll {
			all = true
			break
		}
		if _, has := ngs[g]; !has {
			ngs[g] = i
			i++
		}
	}
	if all {
		params.notNocGenres = nil
		return
	}

	params.notNocGenres = make([]NocGenre, len(ngs))
	for f, idx := range ngs {
		params.notNocGenres[idx] = f
	}
}

// ClearNotNocGenres clear not word fields setting
func (params *SearchR18Params) ClearNotNocGenres() { params.notNocGenres = nil }

func (params *SearchR18Params) queryFromNotNocGenre() url.Values {
	vs := make(url.Values)
	l := len(params.notNocGenres)
	if l == 0 {
		return vs
	}

	codes := make([]string, l)
	for i, g := range params.notNocGenres {
		codes[i] = fmt.Sprintf("%d", g)
	}
	vs.Set(keyNotNocGenre, strings.Join(codes, "-"))
	return vs
}
