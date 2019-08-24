package narrow

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// NewSearchParams return new search parameter object
func NewSearchParams() *SearchParams {
	params := &SearchParams{}
	params.requiredKeywordFlags = make(map[requiredKeyword]bool)
	params.lengths = make(minmaxPair)
	params.kaiwaritus = make(minmaxPair)
	params.readTimes = make(minmaxPair)
	params.sasies = make(minmaxPair)
	return params
}

// ToURL return full URL or nil if params contains invalid condition
func (params *SearchParams) ToURL() (*url.URL, error) {
	return params.makeFullURL(params.endPointURL(), params.toQueryFuncs())
}

type toQueryFunc func() url.Values

func (params *SearchParams) makeFullURL(endPoint string, funcs []toQueryFunc) (*url.URL, error) {
	fullURL, err := url.Parse(endPoint)
	if err != nil {
		return nil, err
	}

	q := fullURL.Query()
	defer func() { fullURL.RawQuery = q.Encode() }()

	// set output format
	q.Add(outputFormatKey, outputFormat)

	if ok, err := params.Valid(); !ok {
		return nil, err
	}

	for _, f := range funcs {
		if m := f(); len(m) != 0 {
			for k, vs := range m {
				q[k] = vs
			}
		}
	}

	return fullURL, nil
}

func (params *SearchParams) endPointURL() string {
	return NarouAPIEndPoint
}

func (params *SearchParams) toQueryFuncs() []toQueryFunc {
	return []toQueryFunc{
		params.queryFromStart,
		params.queryFromLimit,
		params.queryFromOutputField,
		params.queryFromWord,
		params.queryFromNotWord,
		params.queryFromSearchField,
		params.queryFromBigGenre,
		params.queryFromNotBigGenre,
		params.queryFromGenre,
		params.queryFromNotGenre,
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

// Valid returns params is OK or not
func (params *SearchParams) Valid() (bool, error) {
	if params == nil {
		return true, nil
	}
	// TODO: impl.
	return true, nil
}

const minLimit, maxLimit = 1, 500

// Limit return `lim` parameter
func (params *SearchParams) Limit() int { return params.limit }

// SetLimit set `lim` parameter
func (params *SearchParams) SetLimit(limit int) {
	if limit < minLimit || limit > maxLimit {
		return
	}
	params.limit = limit
}

// ClearLimit clear `lim` parameter
func (params *SearchParams) ClearLimit() { params.limit = 0 }

func (params *SearchParams) queryFromLimit() url.Values {
	vs := make(url.Values)
	if params.limit != 0 {
		vs.Set("lim", fmt.Sprintf("%d", params.limit))
	}
	return vs
}

const minOffset, maxOffset = 1, 2000

// Start return `st` (offset) parameter
func (params *SearchParams) Start() int { return params.offset }

// SetStart set `st` (offset) parameter
func (params *SearchParams) SetStart(start int) {
	if start < minOffset || start > maxOffset {
		return
	}
	params.offset = start
}

// ClearStart clear `st` (offset) parameter
func (params *SearchParams) ClearStart() { params.offset = 0 }

func (params *SearchParams) queryFromStart() url.Values {
	vs := make(url.Values)
	if params.offset != 0 {
		vs.Set("st", fmt.Sprintf("%d", params.offset))
	}
	return vs
}

// Order return `order`
func (params *SearchParams) Order() OrderItem { return params.order }

// SetOrder set `order` param
func (params *SearchParams) SetOrder(order OrderItem) { params.order = order }

// ClearOrder clear `order` param
func (params *SearchParams) ClearOrder() { params.order = OrderItemNew }

func (params *SearchParams) queryFromOrder() url.Values {
	vs := make(url.Values)

	if params.order == OrderItemNew {
		return vs
	}
	if name, ok := orderItemNames[params.order]; ok {
		vs.Set("order", name)
	}
	return vs
}

// OutputFields rturn `of` parameters
func (params *SearchParams) OutputFields() []OutputField { return params.outputFields }

// AddOutputFields add OutputFields
func (params *SearchParams) AddOutputFields(fields []OutputField) {
	fs := make(map[OutputField]int)
	i := 0
	for _, f := range params.outputFields {
		fs[f] = i
		i++
	}

	all := false
	for _, f := range fields {
		if f == OutputFieldAll {
			all = true
			break
		}
		if _, has := fs[f]; !has {
			fs[f] = i
			i++
		}
	}
	if all {
		params.outputFields = nil
		return
	}

	params.outputFields = make([]OutputField, len(fs))
	for f, idx := range fs {
		params.outputFields[idx] = f
	}
}

// ClearOutputFields clear output fields setting
func (params *SearchParams) ClearOutputFields() { params.outputFields = nil }

func (params *SearchParams) queryFromOutputField() url.Values {
	vs := make(url.Values)
	l := len(params.outputFields)
	if l == 0 {
		return vs
	}

	all := false
	fields := make([]string, 0, l)
	for _, field := range params.outputFields {
		if field == OutputFieldAll {
			all = true
			break
		}
		name, ok := outputFieldShortNames[field]
		if ok {
			fields = append(fields, name)
		} else {
			// just ignore. TODO: logging?
		}
	}

	if all || len(fields) == 0 {
		// no need 'of' query
		return vs
	}

	vs.Set(keyOutputField, strings.Join(fields, "-"))
	return vs
}

// Words return `word`  parameter
func (params *SearchParams) Words() []string { return params.words }

// AddWords add search words
func (params *SearchParams) AddWords(words []string) {
	ws := make(map[string]int)
	i := 0
	for _, w := range params.words {
		ws[w] = i
		i++
	}
	for _, w := range words {
		if _, has := ws[w]; !has {
			ws[w] = i
			i++
		}
	}

	params.words = make([]string, len(ws))
	for w, idx := range ws {
		params.words[idx] = w
	}
}

// ClearWords clear search word fields setting
func (params *SearchParams) ClearWords() { params.words = nil }

func (params *SearchParams) queryFromWord() url.Values {
	vs := make(url.Values)
	l := len(params.words)
	if l == 0 {
		return vs
	}
	vs.Set(keyWord, strings.Join(params.words, " "))
	return vs
}

// NotWords return `notword`  parameter
func (params *SearchParams) NotWords() []string { return params.notWords }

// AddNotWords add search ignore words
func (params *SearchParams) AddNotWords(words []string) {
	ws := make(map[string]int)
	i := 0
	for _, w := range params.notWords {
		ws[w] = i
		i++
	}
	for _, w := range words {
		if _, has := ws[w]; !has {
			ws[w] = i
			i++
		}
	}
	params.notWords = make([]string, len(ws))
	for w, idx := range ws {
		params.notWords[idx] = w
	}
}

// ClearNotWords clear not word fields setting
func (params *SearchParams) ClearNotWords() { params.notWords = nil }

func (params *SearchParams) queryFromNotWord() url.Values {
	vs := make(url.Values)
	if len(params.notWords) == 0 {
		return vs
	}
	vs.Set(keyNotWord, strings.Join(params.notWords, " "))
	return vs
}

// SearchFields return search fields
func (params *SearchParams) SearchFields() []SearchField { return params.searchFields }

// AddSearchFields add search fields
func (params *SearchParams) AddSearchFields(fields []SearchField) {
	sfs := make(map[SearchField]int)
	i := 0
	for _, f := range params.searchFields {
		sfs[f] = i
		i++
	}
	all := false
	for _, f := range fields {
		if f == SearchFieldAll {
			all = true
			break
		}
		if _, has := sfs[f]; !has {
			sfs[f] = i
			i++
		}
	}
	if all {
		params.searchFields = nil
		return
	}

	params.searchFields = make([]SearchField, len(sfs))
	for f, idx := range sfs {
		params.searchFields[idx] = f
	}
}

// ClearSearchFields clear not word fields setting
func (params *SearchParams) ClearSearchFields() { params.searchFields = nil }

func (params *SearchParams) queryFromSearchField() url.Values {
	vs := make(url.Values)
	l := len(params.searchFields)
	if l == 0 {
		return vs
	}

	names := make(map[string]int)
	for _, field := range params.searchFields {
		name, ok := searchFieldNames[field]
		if ok {
			names[name] = 1
		}
	}
	if len(names) == 0 {
		return vs
	}
	for k := range names {
		vs.Set(k, "1")
	}
	return vs
}

// BigGenres returns `biggenre` parameter
func (params *SearchParams) BigGenres() []BigGenre { return params.bigGenres }

// AddBigGenres add BigGenre settings
func (params *SearchParams) AddBigGenres(genres []BigGenre) {
	bgs := make(map[BigGenre]int)
	i := 0
	for _, g := range params.bigGenres {
		bgs[g] = i
		i++
	}
	all := false
	for _, g := range genres {
		if g == BigGenreAll {
			all = true
			break
		}
		if _, has := bgs[g]; !has {
			bgs[g] = i
			i++
		}
	}

	if all {
		params.bigGenres = nil
		return
	}
	params.bigGenres = make([]BigGenre, len(bgs))
	for f, idx := range bgs {
		params.bigGenres[idx] = f
	}
}

// ClearBigGenres clear not word fields setting
func (params *SearchParams) ClearBigGenres() { params.bigGenres = nil }

func (params *SearchParams) queryFromBigGenre() url.Values {
	vs := make(url.Values)
	l := len(params.bigGenres)
	if l == 0 {
		return vs
	}

	codes := make([]string, l)
	for i, g := range params.bigGenres {
		codes[i] = fmt.Sprintf("%d", g)
	}
	vs.Set(keyBigGenre, strings.Join(codes, "-"))
	return vs
}

// NotBigGenres returns `notbiggenre` parameter
func (params *SearchParams) NotBigGenres() []BigGenre { return params.notBigGenres }

// AddNotBigGenres add not biggenres setting
func (params *SearchParams) AddNotBigGenres(genres []BigGenre) {
	bgs := make(map[BigGenre]int)
	i := 0
	for _, g := range params.notBigGenres {
		bgs[g] = i
		i++
	}
	all := false
	for _, g := range genres {
		if g == BigGenreAll {
			all = true
			break
		}
		if _, has := bgs[g]; !has {
			bgs[g] = i
			i++
		}
	}
	if all {
		params.notBigGenres = nil
		return
	}

	params.notBigGenres = make([]BigGenre, len(bgs))
	for f, idx := range bgs {
		params.notBigGenres[idx] = f
	}
}

// ClearNotBigGenres clear not word fields setting
func (params *SearchParams) ClearNotBigGenres() { params.notBigGenres = nil }

func (params *SearchParams) queryFromNotBigGenre() url.Values {
	vs := make(url.Values)
	l := len(params.notBigGenres)
	if l == 0 {
		return vs
	}

	codes := make([]string, l)
	for i, g := range params.notBigGenres {
		codes[i] = fmt.Sprintf("%d", g)
	}
	vs.Set(keyNotBigGenre, strings.Join(codes, "-"))
	return vs
}

// Genres returns `genre` parameter
func (params *SearchParams) Genres() []Genre { return params.genres }

// AddGenres add genre setting
func (params *SearchParams) AddGenres(genres []Genre) {
	gs := make(map[Genre]int)
	i := 0
	for _, g := range params.genres {
		gs[g] = i
		i++
	}
	all := false
	for _, g := range genres {
		if g == GenreAll {
			all = true
			break
		}
		if _, has := gs[g]; !has {
			gs[g] = i
			i++
		}

	}
	if all {
		params.genres = nil
		return
	}
	params.genres = make([]Genre, len(gs))
	for g, idx := range gs {
		params.genres[idx] = g
	}
}

// ClearGenres clear not word fields setting
func (params *SearchParams) ClearGenres() { params.genres = nil }

func (params *SearchParams) queryFromGenre() url.Values {
	vs := make(url.Values)
	l := len(params.genres)
	if l == 0 {
		return vs
	}
	codes := make([]string, l)
	for i, k := range params.genres {
		codes[i] = fmt.Sprintf("%d", k)
	}
	vs.Set(keyGenre, strings.Join(codes, "-"))
	return vs
}

// NotGenres returns `notgenre` parameter
func (params *SearchParams) NotGenres() []Genre { return params.notGenres }

// AddNotGenres add genre setting
func (params *SearchParams) AddNotGenres(genres []Genre) {
	gs := make(map[Genre]int)
	i := 0
	for _, g := range params.notGenres {
		gs[g] = i
		i++
	}
	all := false
	for _, g := range genres {
		if g == GenreAll {
			all = true
			break
		}
		if _, has := gs[g]; !has {
			gs[g] = i
			i++
		}
	}
	if all {
		params.notGenres = nil
		return
	}
	params.notGenres = make([]Genre, len(gs))
	for f, idx := range gs {
		params.notGenres[idx] = f
	}
}

// ClearNotGenres clear not word fields setting
func (params *SearchParams) ClearNotGenres() { params.notGenres = nil }

func (params *SearchParams) queryFromNotGenre() url.Values {
	vs := make(url.Values)
	l := len(params.notGenres)
	if l == 0 {
		return vs
	}

	codes := make([]string, l)
	for i, k := range params.notGenres {
		codes[i] = fmt.Sprintf("%d", k)
	}
	vs.Set(keyNotGenre, strings.Join(codes, "-"))
	return vs
}

// UserIDs return `userid` param
func (params *SearchParams) UserIDs() []int { return params.userIDs }

// AddUserIDs add search user ids
func (params *SearchParams) AddUserIDs(users []int) {
	ids := make(map[int]int)
	i := 0
	for _, u := range params.userIDs {
		ids[u] = i
		i++
	}
	for _, u := range users {
		if _, has := ids[u]; !has {
			ids[u] = i
			i++
		}
	}

	params.userIDs = make([]int, len(ids))
	for u, idx := range ids {
		params.userIDs[idx] = u
	}
}

// ClearUserIDs clear user id fields setting
func (params *SearchParams) ClearUserIDs() { params.userIDs = nil }

func (params *SearchParams) queryFromUserID() url.Values {
	vs := make(url.Values)
	l := len(params.userIDs)
	if l == 0 {
		return vs
	}
	ids := make([]string, l)
	for i, id := range params.userIDs {
		ids[i] = fmt.Sprintf("%d", id)
	}
	vs.Set(keyUserID, strings.Join(ids, "-"))
	return vs
}

// IsR15 return r15 flag
func (params *SearchParams) IsR15() bool {
	return params.requiredKeywordFlags[requiredKeywordIsR15]
}

// IsNotR15 return not r15 flag
func (params *SearchParams) IsNotR15() bool {
	return params.requiredKeywordFlags[requiredKeywordIsNotR15]
}

// SetIsR15 set r15 flag, reset not-r15 flag
func (params *SearchParams) SetIsR15(r15 bool) {
	params.requiredKeywordFlags[requiredKeywordIsR15] = r15
	if params.requiredKeywordFlags[requiredKeywordIsR15] {
		params.requiredKeywordFlags[requiredKeywordIsNotR15] = false
	}
}

// SetIsNotR15 set not r15 flag, reset r15 flag
func (params *SearchParams) SetIsNotR15(r15 bool) {
	params.requiredKeywordFlags[requiredKeywordIsNotR15] = r15
	if params.requiredKeywordFlags[requiredKeywordIsNotR15] {
		params.requiredKeywordFlags[requiredKeywordIsR15] = false
	}
}

// IsBL return BL flag
func (params *SearchParams) IsBL() bool {
	return params.requiredKeywordFlags[requiredKeywordIsBL]
}

// IsNotBL return not BL flag
func (params *SearchParams) IsNotBL() bool {
	return params.requiredKeywordFlags[requiredKeywordIsNotBL]
}

// SetIsBL set bl flag, reset not-bl flag
func (params *SearchParams) SetIsBL(bl bool) {
	params.requiredKeywordFlags[requiredKeywordIsBL] = bl
	if params.requiredKeywordFlags[requiredKeywordIsBL] {
		params.requiredKeywordFlags[requiredKeywordIsNotBL] = false
	}
}

// SetIsNotBL set not bl flag, reset bl flag
func (params *SearchParams) SetIsNotBL(bl bool) {
	params.requiredKeywordFlags[requiredKeywordIsNotBL] = bl
	if params.requiredKeywordFlags[requiredKeywordIsNotBL] {
		params.requiredKeywordFlags[requiredKeywordIsBL] = false
	}
}

// IsGL return GL flag
func (params *SearchParams) IsGL() bool {
	return params.requiredKeywordFlags[requiredKeywordIsGL]
}

// IsNotGL return not GL flag
func (params *SearchParams) IsNotGL() bool {
	return params.requiredKeywordFlags[requiredKeywordIsNotGL]
}

// SetIsGL set gl flag, reset not-gl flag
func (params *SearchParams) SetIsGL(gl bool) {
	params.requiredKeywordFlags[requiredKeywordIsGL] = gl
	if params.requiredKeywordFlags[requiredKeywordIsGL] {
		params.requiredKeywordFlags[requiredKeywordIsNotGL] = false
	}
}

// SetIsNotGL set not gl flag, reset gl flag
func (params *SearchParams) SetIsNotGL(gl bool) {
	params.requiredKeywordFlags[requiredKeywordIsNotGL] = gl
	if params.requiredKeywordFlags[requiredKeywordIsNotGL] {
		params.requiredKeywordFlags[requiredKeywordIsGL] = false
	}
}

// IsZankoku return Zankoku flag
func (params *SearchParams) IsZankoku() bool {
	return params.requiredKeywordFlags[requiredKeywordIsZankoku]
}

// IsNotZankoku return not Zankoku flag
func (params *SearchParams) IsNotZankoku() bool {
	return params.requiredKeywordFlags[requiredKeywordIsNotZankoku]
}

// SetIsZankoku set zankoku flag, reset not-zankoku flag
func (params *SearchParams) SetIsZankoku(zankoku bool) {
	params.requiredKeywordFlags[requiredKeywordIsZankoku] = zankoku
	if params.requiredKeywordFlags[requiredKeywordIsZankoku] {
		params.requiredKeywordFlags[requiredKeywordIsNotZankoku] = false
	}
}

// SetIsNotZankoku set not zankoku flag, reset zankoku flag
func (params *SearchParams) SetIsNotZankoku(zankoku bool) {
	params.requiredKeywordFlags[requiredKeywordIsNotZankoku] = zankoku
	if params.requiredKeywordFlags[requiredKeywordIsNotZankoku] {
		params.requiredKeywordFlags[requiredKeywordIsZankoku] = false
	}
}

// IsTensei return Tensei flag
func (params *SearchParams) IsTensei() bool {
	return params.requiredKeywordFlags[requiredKeywordIsTensei]
}

// IsNotTensei return not Tensei flag
func (params *SearchParams) IsNotTensei() bool {
	return params.requiredKeywordFlags[requiredKeywordIsNotTensei]
}

// SetIsTensei set tensei flag, reset not-tensei flag
func (params *SearchParams) SetIsTensei(tensei bool) {
	params.requiredKeywordFlags[requiredKeywordIsTensei] = tensei
	if params.requiredKeywordFlags[requiredKeywordIsTensei] {
		params.requiredKeywordFlags[requiredKeywordIsNotTensei] = false
	}
}

// SetIsNotTensei set not tensei flag, reset tensei flag
func (params *SearchParams) SetIsNotTensei(tensei bool) {
	params.requiredKeywordFlags[requiredKeywordIsNotTensei] = tensei
	if params.requiredKeywordFlags[requiredKeywordIsNotTensei] {
		params.requiredKeywordFlags[requiredKeywordIsTensei] = false
	}
}

// IsTenni return Tenni flag
func (params *SearchParams) IsTenni() bool {
	return params.requiredKeywordFlags[requiredKeywordIsTenni]
}

// IsNotTenni return not Tenni flag
func (params *SearchParams) IsNotTenni() bool {
	return params.requiredKeywordFlags[requiredKeywordIsNotTenni]
}

// SetIsTenni set tenni flag, reset not-tenni flag
func (params *SearchParams) SetIsTenni(tenni bool) {
	params.requiredKeywordFlags[requiredKeywordIsTenni] = tenni
	if params.requiredKeywordFlags[requiredKeywordIsTenni] {
		params.requiredKeywordFlags[requiredKeywordIsNotTenni] = false
	}
}

// SetIsNotTenni set not tenni flag, reset tenni flag
func (params *SearchParams) SetIsNotTenni(tenni bool) {
	params.requiredKeywordFlags[requiredKeywordIsNotTenni] = tenni
	if params.requiredKeywordFlags[requiredKeywordIsNotTenni] {
		params.requiredKeywordFlags[requiredKeywordIsTenni] = false
	}
}

// IsTT return Tenni or Tensei flag
func (params *SearchParams) IsTT() bool {
	return params.requiredKeywordFlags[requiredKeywordIsTT]
}

// SetIsTT set tt flag,
func (params *SearchParams) SetIsTT(tt bool) {
	params.requiredKeywordFlags[requiredKeywordIsTT] = tt
}

func (params *SearchParams) queryFromRequiredKeywords() url.Values {
	vs := make(url.Values)
	if params.IsR15() {
		vs.Set("isr15", "1")
	}
	if params.IsNotR15() {
		vs.Set("notr15", "1")
	}
	if params.IsBL() {
		vs.Set("isbl", "1")
	}
	if params.IsNotBL() {
		vs.Set("notbl", "1")
	}
	if params.IsGL() {
		vs.Set("isgl", "1")
	}
	if params.IsNotGL() {
		vs.Set("notgl", "1")
	}
	if params.IsZankoku() {
		vs.Set("iszankoku", "1")
	}
	if params.IsNotZankoku() {
		vs.Set("notzankoku", "1")
	}
	if params.IsTensei() {
		vs.Set("istesei", "1")
	}
	if params.IsNotTensei() {
		vs.Set("nottensei", "1")
	}
	if params.IsTenni() {
		vs.Set("istenni", "1")
	}
	if params.IsNotTenni() {
		vs.Set("nottenni", "1")
	}
	if params.IsTT() {
		vs.Set("istt", "1")
	}
	return vs
}

// Length return `length` parameter if set, or nil
func (params *SearchParams) Length() *int { return params.lengths.val(mmpSingle) }

// SetLength set length parameter
func (params *SearchParams) SetLength(len int) { params.lengths[mmpSingle] = len }

// ClearLength unset length parameter
func (params *SearchParams) ClearLength() { delete(params.lengths, mmpSingle) }

// MinLength return `minlen` parameter if set, or nil
func (params *SearchParams) MinLength() *int { return params.lengths.val(mmpMin) }

// SetMinLength set minLength parameter
func (params *SearchParams) SetMinLength(minLen int) { params.lengths[mmpMin] = minLen }

// ClearMinLength unset minLength parameter
func (params *SearchParams) ClearMinLength() { delete(params.lengths, mmpMin) }

// MaxLength return `maxlen` parameter if set, or nil
func (params *SearchParams) MaxLength() *int { return params.lengths.val(mmpMax) }

// SetMaxLength set maxlength parameter
func (params *SearchParams) SetMaxLength(maxLen int) { params.lengths[mmpMax] = maxLen }

// ClearMaxLength unset maxLength parameter
func (params *SearchParams) ClearMaxLength() { delete(params.lengths, mmpMax) }

func (params *SearchParams) queryFromLength() url.Values {
	vs := make(url.Values)
	if single, ok := params.lengths[mmpSingle]; ok {
		vs.Set(keyLength, fmt.Sprintf("%d", single))
		return vs
	}
	minLen, hasMin := params.lengths[mmpMin]
	maxLen, hasMax := params.lengths[mmpMax]
	if hasMin && hasMax {
		vs.Set(keyLength, fmt.Sprintf("%d-%d", minLen, maxLen))
		return vs
	}
	if hasMin && !hasMax {
		vs.Set(keyLength, fmt.Sprintf("%d-", minLen))
		return vs
	}
	if !hasMin && hasMax {
		vs.Set(keyLength, fmt.Sprintf("-%d", maxLen))
		return vs
	}
	return vs
}

// Kaiwaritu returns `kaiwaritu` parameter
func (params *SearchParams) Kaiwaritu() *int { return params.kaiwaritus.val(mmpSingle) }

// SetKaiwaritu set kaiwaritu parameter
func (params *SearchParams) SetKaiwaritu(kaiwa int) { params.kaiwaritus[mmpSingle] = kaiwa }

// ClearKaiwaritu unset kaiwaritu parameter
func (params *SearchParams) ClearKaiwaritu() { delete(params.kaiwaritus, mmpSingle) }

// MinKaiwaritu returns `kaiwaritu` min parameter
func (params *SearchParams) MinKaiwaritu() *int { return params.kaiwaritus.val(mmpMin) }

// SetMinKaiwaritu set minKaiwaritu parameter
func (params *SearchParams) SetMinKaiwaritu(minKaiwa int) { params.kaiwaritus[mmpMin] = minKaiwa }

// ClearMinKaiwaritu unset minKaiwaritu parameter
func (params *SearchParams) ClearMinKaiwaritu() { delete(params.kaiwaritus, mmpMin) }

// MaxKaiwaritu returns `kaiwaritu` max parameter
func (params *SearchParams) MaxKaiwaritu() *int { return params.kaiwaritus.val(mmpMax) }

// SetMaxKaiwaritu set maxkaiwaritu parameter
func (params *SearchParams) SetMaxKaiwaritu(maxKaiwa int) { params.kaiwaritus[mmpMax] = maxKaiwa }

// ClearMaxKaiwaritu unset maxKaiwaritu parameter
func (params *SearchParams) ClearMaxKaiwaritu() { delete(params.kaiwaritus, mmpMax) }

func (params *SearchParams) queryFromKaiwaritu() url.Values {
	vs := make(url.Values)
	if single, ok := params.kaiwaritus[mmpSingle]; ok {
		vs.Set("kaiwaritu", fmt.Sprintf("%d", single))
		return vs
	}
	minKaiwa, hasMin := params.kaiwaritus[mmpMin]
	maxKaiwa, hasMax := params.kaiwaritus[mmpMax]
	if hasMin && hasMax {
		vs.Set("kaiwaritu", fmt.Sprintf("%d-%d", minKaiwa, maxKaiwa))
		return vs
	}
	if hasMin && !hasMax {
		vs.Set("kaiwaritu", fmt.Sprintf("%d-", minKaiwa))
		return vs
	}
	if !hasMin && hasMax {
		vs.Set("kaiwaritu", fmt.Sprintf("-%d", maxKaiwa))
		return vs
	}
	return vs
}

// Sasie returns `sasie` parameter
func (params *SearchParams) Sasie() *int { return params.sasies.val(mmpSingle) }

// SetSasie set sasie parameter
func (params *SearchParams) SetSasie(sasie int) { params.sasies[mmpSingle] = sasie }

// ClearSasie unset sasie parameter
func (params *SearchParams) ClearSasie() { delete(params.sasies, mmpSingle) }

// MinSasie returns `sasie` min parameter
func (params *SearchParams) MinSasie() *int { return params.sasies.val(mmpMin) }

// SetMinSasie set minSasie parameter
func (params *SearchParams) SetMinSasie(minSasie int) { params.sasies[mmpMin] = minSasie }

// ClearMinSasie unset minSasie parameter
func (params *SearchParams) ClearMinSasie() { delete(params.sasies, mmpMin) }

// MaxSasie returns `sasie` max parameter
func (params *SearchParams) MaxSasie() *int { return params.sasies.val(mmpMax) }

// SetMaxSasie set maxsasie parameter
func (params *SearchParams) SetMaxSasie(maxSasie int) { params.sasies[mmpMax] = maxSasie }

// ClearMaxSasie unset maxSasie parameter
func (params *SearchParams) ClearMaxSasie() { delete(params.sasies, mmpMax) }

func (params *SearchParams) queryFromSasie() url.Values {
	vs := make(url.Values)
	if sasie, ok := params.sasies[mmpSingle]; ok {
		vs.Set(keySasie, fmt.Sprintf("%d", sasie))
		return vs
	}
	minSasie, hasMin := params.sasies[mmpMin]
	maxSasie, hasMax := params.sasies[mmpMax]
	if hasMin && hasMax {
		vs.Set(keySasie, fmt.Sprintf("%d-%d", minSasie, maxSasie))
		return vs
	}
	if hasMin && !hasMax {
		vs.Set(keySasie, fmt.Sprintf("%d-", minSasie))
		return vs
	}
	if !hasMin && hasMax {
		vs.Set(keySasie, fmt.Sprintf("-%d", maxSasie))
		return vs
	}
	return vs
}

// ReadTime returns `time` max parameter
func (params *SearchParams) ReadTime() *int { return params.readTimes.val(mmpSingle) }

// SetReadTime set readTime parameter
func (params *SearchParams) SetReadTime(rt int) { params.readTimes[mmpSingle] = rt }

// ClearReadTime unset `time` parameter
func (params *SearchParams) ClearReadTime() { delete(params.readTimes, mmpSingle) }

// MinReadTime returns `time` min parameter
func (params *SearchParams) MinReadTime() *int { return params.readTimes.val(mmpMin) }

// SetMinReadTime set `time` min parameter
func (params *SearchParams) SetMinReadTime(minRt int) { params.readTimes[mmpMin] = minRt }

// ClearMinReadTime unset `time` min parameter
func (params *SearchParams) ClearMinReadTime() { delete(params.readTimes, mmpMin) }

// MaxReadTime returns `time` max parameter
func (params *SearchParams) MaxReadTime() *int { return params.readTimes.val(mmpMax) }

// SetMaxReadTime set `time` max parameter
func (params *SearchParams) SetMaxReadTime(maxRt int) { params.readTimes[mmpMax] = maxRt }

// ClearMaxReadTime unset `time` max parameter
func (params *SearchParams) ClearMaxReadTime() { delete(params.readTimes, mmpMax) }

func (params *SearchParams) queryFromReadTime() url.Values {
	vs := make(url.Values)
	if readTime, ok := params.readTimes[mmpSingle]; ok {
		vs.Set(keyReadTime, fmt.Sprintf("%d", readTime))
		return vs
	}
	minReadTime, hasMin := params.readTimes[mmpMin]
	maxReadTime, hasMax := params.readTimes[mmpMax]
	if hasMin && hasMax {
		vs.Set(keyReadTime, fmt.Sprintf("%d-%d", minReadTime, maxReadTime))
		return vs
	}
	if hasMin && !hasMax {
		vs.Set(keyReadTime, fmt.Sprintf("%d-", minReadTime))
		return vs
	}
	if !hasMin && hasMax {
		vs.Set(keyReadTime, fmt.Sprintf("-%d", maxReadTime))
		return vs
	}
	return vs
}

// NCodes returns `ncode` param
func (params *SearchParams) NCodes() []string { return params.ncodes }

// AddNCodes add search ncodes
func (params *SearchParams) AddNCodes(ncodes []string) {
	ws := make(map[string]int)
	for _, w := range ncodes {
		ws[w] = 1
	}
	for _, w := range params.ncodes {
		ws[w] = 1
	}
	i := 0
	params.ncodes = make([]string, len(ws))
	for w := range ws {
		params.ncodes[i] = w
		i++
	}
}

// ClearNCodes clear search NCode fields setting
func (params *SearchParams) ClearNCodes() { params.ncodes = nil }

func (params *SearchParams) queryFromNCode() url.Values {
	vs := make(url.Values)
	if len(params.ncodes) == 0 {
		return vs
	}
	vs.Set(keyNCode, strings.Join(params.ncodes, "-"))
	return vs
}

// NovelState return `type` parameter
func (params *SearchParams) NovelState() NovelState { return params.state }

// SetNovelState set `type` parameter
func (params *SearchParams) SetNovelState(state NovelState) { params.state = state }

// ClearNovelState clear `type` parameter
func (params *SearchParams) ClearNovelState() { params.state = NovelStateAll }

func (params *SearchParams) queryFromNovelState() url.Values {
	vs := make(url.Values)

	if params.state == NovelStateAll {
		return vs
	}
	name, ok := novelStateShortNames[params.state]
	if !ok {
		return vs
	}
	vs.Set(keyNovelState, name)
	return vs
}

// Buntais return `buntai` params
func (params *SearchParams) Buntais() []Buntai { return params.buntais }

// AddBuntais add buntai setting
func (params *SearchParams) AddBuntais(buntais []Buntai) {
	bs := make(map[Buntai]int)
	i := 0
	for _, b := range params.buntais {
		bs[b] = i
		i++
	}
	all := false
	for _, b := range buntais {
		if b == BuntaiAll {
			all = true
			break
		}
		if _, has := bs[b]; !has {
			bs[b] = i
			i++
		}
	}
	if all {
		params.buntais = nil
		return
	}
	params.buntais = make([]Buntai, len(bs))
	for b, idx := range bs {
		params.buntais[idx] = b
	}
}

// ClearBuntais clear not word fields setting
func (params *SearchParams) ClearBuntais() { params.buntais = nil }

func (params *SearchParams) queryFromBuntai() url.Values {
	vs := make(url.Values)
	l := len(params.buntais)
	if l == 0 {
		return vs
	}

	codes := make([]string, l)
	for i, b := range params.buntais {
		codes[i] = fmt.Sprintf("%d", b)
	}
	vs.Set(keyBuntai, strings.Join(codes, "-"))
	return vs
}

// StopState returns `stop` param
func (params *SearchParams) StopState() StopState { return params.stopState }

// SetStopState set `stop` parameter
func (params *SearchParams) SetStopState(state StopState) { params.stopState = state }

// ClearStopState clear `stop` parameter
func (params *SearchParams) ClearStopState() { params.stopState = StopStateAll }

func (params *SearchParams) queryFromStop() url.Values {
	vs := make(url.Values)
	if params.stopState == StopStateAll {
		return vs
	}
	vs.Set(keyStop, fmt.Sprintf("%d", params.stopState))
	return vs
}

// PickupState return `ispickup` parameter
func (params *SearchParams) PickupState() PickupState { return params.pickupState }

// SetPickupState set `ispickup` param
func (params *SearchParams) SetPickupState(pickup PickupState) { params.pickupState = pickup }

// ClearPickupState clear `ispickup` param
func (params *SearchParams) ClearPickupState() { params.pickupState = PickupStateNone }

func (params *SearchParams) queryFromPickup() url.Values {
	vs := make(url.Values)

	switch params.pickupState {
	case PickupStateNone:
	case PickupStateNot:
		vs.Set(keyIsPickup, "0")
	case PickupStatePickup:
		vs.Set(keyIsPickup, "1")
	}
	return vs
}

// LastUpType returns `lastup` parameter,
func (params *SearchParams) LastUpType() LastUpType { return params.lastUp }

// LastUpTimeStamps return `lastup` timestamp Start/End
func (params *SearchParams) LastUpTimeStamps() [2]time.Time {
	if params.lastUp == LastUpTypeTimeStamp {
		return [2]time.Time{params.lastUpStart, params.lastUpEnd}
	}
	return [2]time.Time{}
}

// SetLastUp set `lastup` params
func (params *SearchParams) SetLastUp(ltype LastUpType) { params.lastUp = ltype }

// SetLastUpTerm set `lastup` timestamp
func (params *SearchParams) SetLastUpTerm(start, end time.Time) {
	params.lastUp = LastUpTypeTimeStamp
	params.lastUpStart = start
	params.lastUpEnd = end
}

// ClearLastUp clear `lastup` param
func (params *SearchParams) ClearLastUp() {
	params.lastUp = LastUpTypeNone
	params.lastUpStart = time.Time{}
	params.lastUpEnd = time.Time{}
}

func (params *SearchParams) queryFromLastUp() url.Values {
	vs := make(url.Values)

	if params.lastUp == LastUpTypeNone {
		return vs
	}

	if params.lastUp == LastUpTypeTimeStamp {
		vs.Set(keyLastUp, fmt.Sprintf("%d-%d", params.lastUpStart.Unix(), params.lastUpEnd.Unix()))
		return vs
	}

	if name, ok := lastUpTypeNames[params.lastUp]; ok {
		vs.Set(keyLastUp, name)
		return vs
	}
	return vs
}

// NarouAPIEndPoint contains Narou novel api endpoit
const NarouAPIEndPoint = "https://api.syosetu.com/novelapi/api/"

// NarouR18APIEndPoint is Narou R18 novel api endpoint
const NarouR18APIEndPoint = "https://api.syosetu.com/novel18api/api/"

const outputFormatKey = "out"
const outputFormat = "json"

const (
	keyOutputField = "of"
	keyWord        = "word"
	keyNotWord     = "notword"
	keyBigGenre    = "biggenre"
	keyNotBigGenre = "notbiggenre"
	keyGenre       = "genre"
	keyNotGenre    = "notgenre"
	keyUserID      = "userid"
	keyLength      = "length"
	keyReadTime    = "time"
	keyNCode       = "ncode"
	keySasie       = "sasie"
	keyNovelState  = "type"
	keyBuntai      = "buntai"
	keyStop        = "stop"
	keyIsPickup    = "ispickup"
	keyLastUp      = "lastup"

	keyNocGenre    = "nocgenre"
	keyNotNocGenre = "notnocgenre"
)

var outputFieldShortNames = map[OutputField]string{
	OutputFieldAll:            "",
	OutputFieldTitle:          "t",
	OutputFieldNCode:          "n",
	OutputFieldUserID:         "u",
	OutputFieldWriter:         "w",
	OutputFieldStory:          "s",
	OutputFieldBigGenre:       "bg",
	OutputFieldGenre:          "g",
	OutputFieldKeyword:        "k",
	OutputFieldGeneralFirstUp: "gf",
	OutputFieldGeneralLastUp:  "gl",
	OutputFieldNovelType:      "nt",
	OutputFieldEnd:            "e",
	OutputFieldGeneralAllNo:   "ga",
	OutputFieldLength:         "l",
	OutputFieldTime:           "ti",
	OutputFieldIsStop:         "i",
	OutputFieldIsR15:          "ir",
	OutputFieldIsBL:           "ibl",
	OutputFieldIsGL:           "igl",
	OutputFieldIsZankoku:      "izk",
	OutputFieldIsTensei:       "its",
	OutputFieldIsTenni:        "iti",
	OutputFieldPcOrK:          "p",
	OutputFieldGlobalPoint:    "gp",
	OutputFieldFavNovelCount:  "f",
	OutputFieldReviewCount:    "r",
	OutputFieldAllPoint:       "a",
	OutputFieldAllHyokaCount:  "ah",
	OutputFieldSasieCount:     "sa",
	OutputFieldKaiwaritu:      "ka",
	OutputFieldNovelUpdatedAt: "nu",
	OutputFieldUpdatedAt:      "ua",

	OutputFieldNocGenre: "ng",

	OutputFieldDailyPoint:   "dp",
	OutputFieldWeeklyPoint:  "wp",
	OutputFieldMonthlyPoint: "mp",
	OutputFieldQuaterPoint:  "qp",
	OutputFieldYearlyPoint:  "yp",

	OutputFieldImpressionCount: "imp",
}

var orderItemNames = map[OrderItem]string{
	OrderItemNew:             "new",
	OrderItemFavNovelCount:   "favnovelcnt",
	OrderItemReviewCount:     "reviewcnt",
	OrderItemHyoka:           "hyoka",
	OrderItemHyokaAsc:        "hyokaasc",
	OrderItemImpressionCount: "impressioncnt",
	OrderItemHyokaCount:      "hyokacnt",
	OrderItemHyokaCountAsc:   "hyokacntasc",
	OrderItemWeekly:          "weekly",
	OrderItemLengthDesc:      "lengthdesc",
	OrderItemLengthAsc:       "lengthasc",
	OrderItemNCodeDesc:       "ncodedesc",
	OrderItemOld:             "old",
}

var searchFieldNames = map[SearchField]string{
	SearchFieldTitle:   "title",
	SearchFieldStory:   "ex",
	SearchFieldKeyword: "keyword",
	SearchFieldWriter:  "wname",
}

var novelStateShortNames = map[NovelState]string{
	NovelStateAll:                 "",
	NovelStateShortStory:          "t",
	NovelStateRensaiRunning:       "r",
	NovelStateRensaiEnded:         "er",
	NovelStateRensaiAll:           "re",
	NovelStateShortAndRensaiEnded: "ter",
}

var lastUpTypeNames = map[LastUpType]string{
	LastUpTypeThisWeek:  "thisweek",
	LastUpTypeLastWeek:  "lastweek",
	LastUpTypeSevenDay:  "sevenday",
	LastUpTypeThisMonth: "thismonth",
	LastUpTypeLastMonth: "lastmonth",
}

type requiredKeyword int

const (
	requiredKeywordNone requiredKeyword = iota
	requiredKeywordIsR15
	requiredKeywordIsNotR15
	requiredKeywordIsBL
	requiredKeywordIsNotBL
	requiredKeywordIsGL
	requiredKeywordIsNotGL
	requiredKeywordIsZankoku
	requiredKeywordIsNotZankoku
	requiredKeywordIsTensei
	requiredKeywordIsNotTensei
	requiredKeywordIsTenni
	requiredKeywordIsNotTenni
	requiredKeywordIsTT
)

type mmpKey int

const (
	mmpSingle mmpKey = iota
	mmpMin
	mmpMax
)

type minmaxPair map[mmpKey]int

func (mmp minmaxPair) val(k mmpKey) *int {
	if v, ok := mmp[k]; ok {
		return &v
	}
	return nil
}
