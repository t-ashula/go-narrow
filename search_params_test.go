package narrow

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func parseURL(s string) *url.URL {
	u, _ := url.Parse(s)
	return u
}

func makeValues(kvs [][2]string) url.Values {
	vs := make(url.Values)
	for _, kv := range kvs {
		vs.Add(kv[0], kv[1])
	}
	return vs
}

func TestSearchParams_ToURL(t *testing.T) {
	tests := []struct {
		name    string
		params  *SearchParams
		want    *url.URL
		wantErr bool
	}{
		{"all default params",
			&SearchParams{},
			parseURL("https://api.syosetu.com/novelapi/api/?out=json"), false},
		{"with st(offset)",
			&SearchParams{offset: 42},
			parseURL("https://api.syosetu.com/novelapi/api/?out=json&st=42"), false},
		{"with limit",
			&SearchParams{limit: 42},
			parseURL("https://api.syosetu.com/novelapi/api/?out=json&lim=42"), false},
		{"with st(offset), limit",
			&SearchParams{offset: 42, limit: 30},
			parseURL("https://api.syosetu.com/novelapi/api/?out=json&st=42&lim=30"), false},
		{"with output field All",
			&SearchParams{outputFields: []OutputField{OutputFieldAll}},
			parseURL("https://api.syosetu.com/novelapi/api/?out=json"), false},
		{"with output field Title",
			&SearchParams{outputFields: []OutputField{OutputFieldTitle}},
			parseURL("https://api.syosetu.com/novelapi/api/?out=json&of=t"), false},
		{"with output field Title, NCode",
			&SearchParams{outputFields: []OutputField{OutputFieldTitle, OutputFieldNCode}},
			parseURL("https://api.syosetu.com/novelapi/api/?out=json&of=t-n"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.params.ToURL()
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchParams.ToURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// should not change path
			if got.Path != tt.want.Path {
				t.Errorf("SearchParams.ToURL().Path = %v, want.Path %v", got.Path, tt.want.Path)
				return
			}

			// query part
			gr, wr := got.RawQuery, tt.want.RawQuery
			if len(gr) != len(wr) {
				t.Errorf("SearchParams.ToURL().RawQuery = %v, len = %v, want.RawQuery = %v, len = %v", gr, len(gr), wr, len(wr))
				return
			}

			gq, wq := got.Query(), tt.want.Query()
			if !reflect.DeepEqual(gq, wq) {
				t.Errorf("SearchParams.ToURL().Query %v, want.Query %v", gq, wq)
				return
			}
		})
	}
}

func TestSearchParams_Valid(t *testing.T) {
	tests := []struct {
		name    string
		params  *SearchParams
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.params.Valid()
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchParams.Valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SearchParams.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_Limit(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   int
	}{
		{"default Limit returns 0", &SearchParams{}, 0},
		{"Limit returns .limit", &SearchParams{limit: 42}, 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.Limit(); got != tt.want {
				t.Errorf("SearchParams.Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetLimit(t *testing.T) {
	type args struct {
		limit int
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   int
	}{
		{"set 42 ok", &SearchParams{}, args{42}, 42},
		{"should not set < 1", &SearchParams{limit: 42}, args{0}, 42},
		{"should not set > 500", &SearchParams{limit: 42}, args{1000}, 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetLimit(tt.args.limit)
			if tt.params.limit != tt.want {
				t.Errorf("SearchParams.SetLimit(%v) should be change limit %v, but %v", tt.args.limit, tt.want, tt.params.limit)
			}
		})
	}
}

func TestSearchParams_ClearLimit(t *testing.T) {
	params := &SearchParams{}
	params.SetLimit(42)

	params.ClearLimit()
	if params.Limit() != 0 {
		t.Errorf("SearchParams.ClearLimit() should change limit be  0, but %v", params.Limit())
	}
}

func TestSearchParams_queryFromLimit(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no limit should be no query", &SearchParams{}, makeValues([][2]string{})},
		{"limit:X to lim=X", &SearchParams{limit: 42}, makeValues([][2]string{{"lim", "42"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromLimit(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_Start(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   int
	}{
		{"no start", &SearchParams{}, 0},
		{"get start", &SearchParams{offset: 42}, 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.Start(); got != tt.want {
				t.Errorf("SearchParams.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetStart(t *testing.T) {
	type args struct {
		start int
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   int
	}{
		{"set valid value", &SearchParams{}, args{42}, 42},
		{"should not set < 1", &SearchParams{offset: 42}, args{0}, 42},
		{"should set 1", &SearchParams{offset: 42}, args{1}, 1},
		{"should set 2000", &SearchParams{offset: 42}, args{2000}, 2000},
		{"should not set 2001", &SearchParams{offset: 42}, args{2001}, 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetStart(tt.args.start)
			if tt.params.Start() != tt.want {
				t.Errorf("SearchParams.SetStart(%v) should be change offset %v, but %v", tt.args.start, tt.want, tt.params.Start())
			}
		})
	}
}

func TestSearchParams_ClearStart(t *testing.T) {
	params := &SearchParams{}
	params.SetStart(42)
	params.ClearStart()
	if params.Start() != 0 {
		t.Errorf("SearchParams.ClearStart() should change offset be 0, but %v", params.Start())
	}
}

func TestSearchParams_queryFromStart(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no start no query", &SearchParams{}, makeValues([][2]string{})},
		{"Start:X to st=X", &SearchParams{offset: 42}, makeValues([][2]string{{"st", "42"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromStart(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_Order(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   OrderItem
	}{
		{"default order is New", &SearchParams{}, OrderItemNew},
		{"get orderItem", &SearchParams{order: OrderItemHyoka}, OrderItemHyoka},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.Order(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.Order() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetOrder(t *testing.T) {
	type args struct {
		order OrderItem
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   OrderItem
	}{
		{"set orderitem should change .order", &SearchParams{}, args{order: OrderItemWeekly}, OrderItemWeekly},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetOrder(tt.args.order)
			if tt.params.Order() != tt.want {
				t.Errorf("SearchParams.SetOrder(%v) should be change .order %v, but %v", tt.args.order, tt.want, tt.params.order)
			}
		})
	}
}

func TestSearchParams_ClearOrder(t *testing.T) {
	params := &SearchParams{}
	params.SetOrder(OrderItemFavNovelCount)
	params.ClearOrder()
	if params.Order() != OrderItemNew {
		t.Errorf("SearchParams.ClearOrder() should change offset be 0, but %v", params.Order())
	}
}

func TestSearchParams_queryFromOrder(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no order, no query", &SearchParams{}, makeValues([][2]string{})},
		{"order ItemNew, no query", &SearchParams{order: OrderItemNew}, makeValues([][2]string{})},
		{"order:ItemOld, order=old", &SearchParams{order: OrderItemOld}, makeValues([][2]string{{"order", "old"}})},
		{"order:ItemImpressionCount, order=impressioncnt", &SearchParams{order: OrderItemImpressionCount}, makeValues([][2]string{{"order", "impressioncnt"}})},
		{"order:Unknown, no query", &SearchParams{order: 2000}, makeValues([][2]string{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromOrder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.orderToQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_OutputFields(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []OutputField
	}{
		{"no outpitfield, empty array", &SearchParams{}, nil},
		{"outputfield", &SearchParams{outputFields: []OutputField{OutputFieldTitle}}, []OutputField{OutputFieldTitle}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.OutputFields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.OutputFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddOutputFields(t *testing.T) {
	type args struct {
		ofs []OutputField
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []OutputField
	}{
		{"empty array should not change", &SearchParams{}, args{[]OutputField{}}, []OutputField{}},
		{"args include All should clear fields",
			&SearchParams{outputFields: []OutputField{OutputFieldTitle}},
			args{[]OutputField{OutputFieldAll, OutputFieldNCode}}, nil},
		{"add output fields",
			&SearchParams{outputFields: []OutputField{OutputFieldTitle}},
			args{[]OutputField{OutputFieldNCode}}, []OutputField{OutputFieldTitle, OutputFieldNCode}},
		{"merge output fields",
			&SearchParams{outputFields: []OutputField{OutputFieldTitle}},
			args{[]OutputField{OutputFieldNCode, OutputFieldTitle}}, []OutputField{OutputFieldTitle, OutputFieldNCode}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddOutputFields(tt.args.ofs)
			result := tt.params.OutputFields()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddOutputFields(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.ofs, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchParams_ClearOutputFields(t *testing.T) {
	params := &SearchParams{}
	params.AddOutputFields([]OutputField{OutputFieldEnd})
	params.ClearOutputFields()
	if params.OutputFields() != nil {
		t.Errorf("SearchParams.ClearOutputFields() should change outputfields be nil, but %v", params.OutputFields())
	}
}

func TestSearchParams_queryFromOutputField(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no field, no query", &SearchParams{}, makeValues([][2]string{})},
		{"outputFieldTitle should of=t",
			&SearchParams{outputFields: []OutputField{OutputFieldTitle}},
			makeValues([][2]string{{"of", "t"}})},
		{"outputFieldNCode should of=n",
			&SearchParams{outputFields: []OutputField{OutputFieldNCode}},
			makeValues([][2]string{{"of", "n"}})},
		{"outputFieldNCode and Title should of=n-t",
			&SearchParams{outputFields: []OutputField{OutputFieldNCode, OutputFieldTitle}},
			makeValues([][2]string{{"of", "n-t"}})},
		// TODO: more output field
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromOutputField(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromOutputField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_Words(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []string
	}{
		{"no words", &SearchParams{}, nil},
		{"words", &SearchParams{words: []string{"abc", "def"}}, []string{"abc", "def"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.Words(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.Words() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddWords(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []string
	}{
		{"words + []", &SearchParams{words: []string{"some", "words"}}, args{[]string{}}, []string{"some", "words"}},
		{"[] + words", &SearchParams{words: []string{}}, args{[]string{"some", "words"}}, []string{"some", "words"}},
		{"words + words", &SearchParams{words: []string{"some", "words"}}, args{[]string{"more", "phrase"}}, []string{"some", "words", "more", "phrase"}},
		{"merge words", &SearchParams{words: []string{"some", "words"}}, args{[]string{"more", "words"}}, []string{"some", "words", "more"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddWords(tt.args.words)
			result := tt.params.Words()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddWords(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.words, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchParams_ClearWords(t *testing.T) {
	params := &SearchParams{}
	words := []string{"first", "word"}
	params.AddWords(words)
	params.ClearWords()
	if params.Words() != nil {
		t.Errorf("SearchParams.ClearWords() should change words be nil, but %v", params.Words())
	}
	after := []string{"new", "word"}
	params.AddWords(after)
	if got := params.Words(); !reflect.DeepEqual(got, after) {
		t.Errorf("AddWords after SearchParams.ClearWords() should change words be %v, but %v", after, got)
	}
}

func TestSearchParams_queryFromWord(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no words, no query", &SearchParams{}, makeValues([][2]string{})},
		{`word:["word1"], word=word`, &SearchParams{words: []string{"word1"}}, makeValues([][2]string{{"word", "word1"}})},
		{`word:["word1", "word2"], word=word%20word2`,
			&SearchParams{words: []string{"word1", "word2"}},
			makeValues([][2]string{{"word", "word1 word2"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromWord(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_NotWords(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []string
	}{
		{"no words", &SearchParams{}, nil},
		{"words", &SearchParams{notWords: []string{"abc", "def"}}, []string{"abc", "def"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.NotWords(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.NotWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddNotWords(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []string
	}{
		{"words + []", &SearchParams{notWords: []string{"some", "words"}}, args{[]string{}}, []string{"some", "words"}},
		{"[] + words", &SearchParams{notWords: []string{}}, args{[]string{"some", "words"}}, []string{"some", "words"}},
		{"words + words", &SearchParams{notWords: []string{"some", "words"}}, args{[]string{"more", "phrase"}}, []string{"some", "words", "more", "phrase"}},
		{"merge words", &SearchParams{notWords: []string{"some", "words"}}, args{[]string{"more", "words"}}, []string{"some", "words", "more"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddNotWords(tt.args.words)
			result := tt.params.NotWords()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddWords(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.words, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchParams_ClearNotWords(t *testing.T) {
	params := &SearchParams{}
	words := []string{"first", "word"}
	params.AddNotWords(words)
	params.ClearNotWords()
	if params.NotWords() != nil {
		t.Errorf("SearchParams.ClearNotWords() should change words be nil, but %v", params.NotWords())
	}
	after := []string{"new", "word"}
	params.AddNotWords(after)
	if got := params.NotWords(); !reflect.DeepEqual(got, after) {
		t.Errorf("AddWords after SearchParams.ClearNotWords() should change words be %v, but %v", after, got)
	}
}

func TestSearchParams_queryFromNotWord(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no words, no query", &SearchParams{}, makeValues([][2]string{})},
		{`notword:["word1"], notword=word`, &SearchParams{notWords: []string{"word1"}}, makeValues([][2]string{{"notword", "word1"}})},
		{`notword:["word1", "word2"], notword=word%20word2`,
			&SearchParams{notWords: []string{"word1", "word2"}},
			makeValues([][2]string{{"notword", "word1 word2"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromNotWord(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromNotWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SearchFields(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []SearchField
	}{
		{"default", &SearchParams{}, nil},
		{"search fields", &SearchParams{searchFields: []SearchField{SearchFieldTitle}}, []SearchField{SearchFieldTitle}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.SearchFields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.SearchFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddSearchFields(t *testing.T) {
	type args struct {
		fields []SearchField
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []SearchField
	}{
		{"field + []", &SearchParams{searchFields: []SearchField{SearchFieldKeyword}}, args{[]SearchField{}}, []SearchField{SearchFieldKeyword}},
		{"[] + field", &SearchParams{searchFields: []SearchField{}}, args{[]SearchField{SearchFieldStory}}, []SearchField{SearchFieldStory}},
		{"fields + fields",
			&SearchParams{searchFields: []SearchField{SearchFieldKeyword}}, args{[]SearchField{SearchFieldStory}},
			[]SearchField{SearchFieldKeyword, SearchFieldStory}},
		{"merge fields",
			&SearchParams{searchFields: []SearchField{SearchFieldKeyword, SearchFieldStory}}, args{[]SearchField{SearchFieldStory}},
			[]SearchField{SearchFieldKeyword, SearchFieldStory}},
		{"args include All should clear fields",
			&SearchParams{searchFields: []SearchField{SearchFieldKeyword}},
			args{[]SearchField{SearchFieldAll, SearchFieldStory}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddSearchFields(tt.args.fields)
			result := tt.params.SearchFields()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddSearchFields(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.fields, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchParams_ClearSearchFields(t *testing.T) {
	params := &SearchParams{searchFields: []SearchField{SearchFieldKeyword}}
	params.ClearSearchFields()
	if params.SearchFields() != nil {
		t.Errorf("SearchParams.ClearSearchFields() should change SearchField be nil, but %v", params.SearchFields())
	}
}

func TestSearchParams_queryFromSearchField(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no fields, no query", &SearchParams{}, makeValues([][2]string{})},
		{`fields:[Title], title=1`, &SearchParams{searchFields: []SearchField{SearchFieldTitle}}, makeValues([][2]string{{"title", "1"}})},
		{`fields:[Keyword], keyword=1`, &SearchParams{searchFields: []SearchField{SearchFieldKeyword}}, makeValues([][2]string{{"keyword", "1"}})},
		{`fields:[Story], ex=1`, &SearchParams{searchFields: []SearchField{SearchFieldStory}}, makeValues([][2]string{{"ex", "1"}})},
		{`fields:[Writer], wname=1`, &SearchParams{searchFields: []SearchField{SearchFieldWriter}}, makeValues([][2]string{{"wname", "1"}})},
		{`fields:[Keyword, Title], keyword=1&title=1`,
			&SearchParams{searchFields: []SearchField{SearchFieldKeyword, SearchFieldTitle}},
			makeValues([][2]string{{"title", "1"}, {"keyword", "1"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromSearchField(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromSearchField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_BigGenres(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []BigGenre
	}{
		{"default", &SearchParams{}, nil},
		{"big genres", &SearchParams{bigGenres: []BigGenre{BigGenreFantasy}}, []BigGenre{BigGenreFantasy}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.BigGenres(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.BigGenres() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddBigGenres(t *testing.T) {
	type args struct {
		genres []BigGenre
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []BigGenre
	}{
		{"genre + []", &SearchParams{bigGenres: []BigGenre{BigGenreFantasy}}, args{[]BigGenre{}}, []BigGenre{BigGenreFantasy}},
		{"[] + genre", &SearchParams{bigGenres: []BigGenre{}}, args{[]BigGenre{BigGenreBungei}}, []BigGenre{BigGenreBungei}},
		{"genres + genres",
			&SearchParams{bigGenres: []BigGenre{BigGenreFantasy}}, args{[]BigGenre{BigGenreBungei}},
			[]BigGenre{BigGenreFantasy, BigGenreBungei}},
		{"merge genres",
			&SearchParams{bigGenres: []BigGenre{BigGenreFantasy, BigGenreBungei}}, args{[]BigGenre{BigGenreBungei}},
			[]BigGenre{BigGenreFantasy, BigGenreBungei}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddBigGenres(tt.args.genres)
			result := tt.params.BigGenres()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddBigGenres(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.genres, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchParams_ClearBigGenres(t *testing.T) {
	params := &SearchParams{bigGenres: []BigGenre{BigGenreFantasy}}
	params.ClearBigGenres()
	if params.BigGenres() != nil {
		t.Errorf("SearchParams.BigGenres() should change BigGenres be nil, but %v", params.BigGenres())
	}
}

func TestSearchParams_queryFromBigGenre(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no bigGenres, no query", &SearchParams{}, makeValues([][2]string{})},
		{`bigGenres:[Fantasy], biggenre=2`, &SearchParams{bigGenres: []BigGenre{BigGenreFantasy}}, makeValues([][2]string{{"biggenre", "2"}})},
		{`bigGenres:[Other], biggenre=99`, &SearchParams{bigGenres: []BigGenre{BigGenreOther}}, makeValues([][2]string{{"biggenre", "99"}})},
		{`bigGenres:[Fantasy, Bungei], biggenre=2-3`,
			&SearchParams{bigGenres: []BigGenre{BigGenreFantasy, BigGenreBungei}}, makeValues([][2]string{{"biggenre", "2-3"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromBigGenre(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromBigGenre() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_NotBigGenres(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []BigGenre
	}{
		{"default", &SearchParams{}, nil},
		{"not big genres", &SearchParams{notBigGenres: []BigGenre{BigGenreFantasy}}, []BigGenre{BigGenreFantasy}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.NotBigGenres(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.NotBigGenres() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddNotBigGenres(t *testing.T) {
	type args struct {
		genres []BigGenre
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []BigGenre
	}{
		{"genre + []", &SearchParams{notBigGenres: []BigGenre{BigGenreFantasy}}, args{[]BigGenre{}}, []BigGenre{BigGenreFantasy}},
		{"[] + genre", &SearchParams{notBigGenres: []BigGenre{}}, args{[]BigGenre{BigGenreBungei}}, []BigGenre{BigGenreBungei}},
		{"genres + genres",
			&SearchParams{notBigGenres: []BigGenre{BigGenreFantasy}}, args{[]BigGenre{BigGenreBungei}},
			[]BigGenre{BigGenreFantasy, BigGenreBungei}},
		{"merge genres",
			&SearchParams{notBigGenres: []BigGenre{BigGenreFantasy, BigGenreBungei}}, args{[]BigGenre{BigGenreBungei}},
			[]BigGenre{BigGenreFantasy, BigGenreBungei}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddNotBigGenres(tt.args.genres)
			result := tt.params.NotBigGenres()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddNotBigGenres(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.genres, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchParams_ClearNotBigGenres(t *testing.T) {
	params := &SearchParams{notBigGenres: []BigGenre{BigGenreFantasy}}
	params.ClearNotBigGenres()
	if params.NotBigGenres() != nil {
		t.Errorf("SearchParams.ClearNotBigGenres() should change NotBigGenres be nil, but %v", params.NotBigGenres())
	}
}

func TestSearchParams_queryFromNotBigGenre(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no notBigGenres, no query", &SearchParams{}, makeValues([][2]string{})},
		{`notBigGenres:[Fantasy], notbiggenre=2`, &SearchParams{notBigGenres: []BigGenre{BigGenreFantasy}}, makeValues([][2]string{{"notbiggenre", "2"}})},
		{`notBigGenres:[Other], notbiggenre=99`, &SearchParams{notBigGenres: []BigGenre{BigGenreOther}}, makeValues([][2]string{{"notbiggenre", "99"}})},
		{`notBigGenres:[Fantasy, Bungei], notbiggenre=2-3`,
			&SearchParams{notBigGenres: []BigGenre{BigGenreFantasy, BigGenreBungei}}, makeValues([][2]string{{"notbiggenre", "2-3"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromNotBigGenre(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromNotBigGenre() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_Genres(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []Genre
	}{
		{"default", &SearchParams{}, nil},
		{"genres", &SearchParams{genres: []Genre{GenreBungeiAction}}, []Genre{GenreBungeiAction}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.Genres(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.Genres() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddGenres(t *testing.T) {
	type args struct {
		genres []Genre
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []Genre
	}{
		{"genre + []", &SearchParams{genres: []Genre{GenreFantasyHighFantasy}}, args{[]Genre{}}, []Genre{GenreFantasyHighFantasy}},
		{"[] + genre", &SearchParams{genres: []Genre{}}, args{[]Genre{GenreBungeiAction}}, []Genre{GenreBungeiAction}},
		{"genres + genres",
			&SearchParams{genres: []Genre{GenreFantasyHighFantasy}}, args{[]Genre{GenreBungeiAction}},
			[]Genre{GenreFantasyHighFantasy, GenreBungeiAction}},
		{"merge genres",
			&SearchParams{genres: []Genre{GenreFantasyHighFantasy, GenreBungeiAction}}, args{[]Genre{GenreBungeiAction}},
			[]Genre{GenreFantasyHighFantasy, GenreBungeiAction}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddGenres(tt.args.genres)
			result := tt.params.Genres()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddGenres(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.genres, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchParams_ClearGenres(t *testing.T) {
	params := &SearchParams{genres: []Genre{GenreFantasyHighFantasy}}
	params.ClearGenres()
	if params.Genres() != nil {
		t.Errorf("SearchParams.ClearGenres() should change Genres be nil, but %v", params.Genres())
	}
}

func TestSearchParams_queryFromGenre(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no Genres, no query", &SearchParams{}, makeValues([][2]string{})},
		{`Genres:[Fantasy], genre=201`, &SearchParams{genres: []Genre{GenreFantasyHighFantasy}}, makeValues([][2]string{{"genre", "201"}})},
		{`Genres:[OtherOther], genre=9999`, &SearchParams{genres: []Genre{GenreOtherOther}}, makeValues([][2]string{{"genre", "9999"}})},
		{`Genres:[Fantasy, Bungei], genre=201-301`,
			&SearchParams{genres: []Genre{GenreFantasyHighFantasy, GenreBungeiJunbungaku}}, makeValues([][2]string{{"genre", "201-301"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromGenre(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromGenre() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_NotGenres(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []Genre
	}{
		{"default", &SearchParams{}, nil},
		{"not genres", &SearchParams{notGenres: []Genre{GenreBungeiAction}}, []Genre{GenreBungeiAction}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.NotGenres(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.NotGenres() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddNotGenres(t *testing.T) {
	type args struct {
		genres []Genre
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []Genre
	}{
		{"genre + []", &SearchParams{notGenres: []Genre{GenreFantasyHighFantasy}}, args{[]Genre{}}, []Genre{GenreFantasyHighFantasy}},
		{"[] + genre", &SearchParams{notGenres: []Genre{}}, args{[]Genre{GenreBungeiAction}}, []Genre{GenreBungeiAction}},
		{"genres + genres",
			&SearchParams{notGenres: []Genre{GenreFantasyHighFantasy}}, args{[]Genre{GenreBungeiAction}},
			[]Genre{GenreFantasyHighFantasy, GenreBungeiAction}},
		{"merge genres",
			&SearchParams{notGenres: []Genre{GenreFantasyHighFantasy, GenreBungeiAction}}, args{[]Genre{GenreBungeiAction}},
			[]Genre{GenreFantasyHighFantasy, GenreBungeiAction}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddNotGenres(tt.args.genres)
			result := tt.params.NotGenres()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddNotGenres(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.genres, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchParams_ClearNotGenres(t *testing.T) {
	params := &SearchParams{notGenres: []Genre{GenreFantasyHighFantasy}}
	params.ClearNotGenres()
	if params.NotGenres() != nil {
		t.Errorf("SearchParams.ClearNotGenres() should change NotGenres be nil, but %v", params.NotGenres())
	}
}

func TestSearchParams_queryFromNotGenre(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no notGenres, no query", &SearchParams{}, makeValues([][2]string{})},
		{`NotGenres:[Fantasy], genre=201`, &SearchParams{notGenres: []Genre{GenreFantasyHighFantasy}}, makeValues([][2]string{{"notgenre", "201"}})},
		{`NotGenres:[OtherOther], genre=9999`, &SearchParams{notGenres: []Genre{GenreOtherOther}}, makeValues([][2]string{{"notgenre", "9999"}})},
		{`NotGenres:[Fantasy, Bungei], genre=201-301`,
			&SearchParams{notGenres: []Genre{GenreFantasyHighFantasy, GenreBungeiJunbungaku}}, makeValues([][2]string{{"notgenre", "201-301"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromNotGenre(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromNotGenre() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_UserIDs(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []int
	}{
		{"default", &SearchParams{}, nil},
		{"not genres", &SearchParams{userIDs: []int{123}}, []int{123}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.UserIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.UserIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddUserIDs(t *testing.T) {
	type args struct {
		users []int
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []int
	}{
		{"userIDs + []", &SearchParams{userIDs: []int{123, 456}}, args{[]int{}}, []int{123, 456}},
		{"[] + userIDs", &SearchParams{userIDs: []int{}}, args{[]int{123, 456}}, []int{123, 456}},
		{"userIDs + userIDs", &SearchParams{userIDs: []int{123, 456}}, args{[]int{1011, 789}}, []int{123, 456, 1011, 789}},
		{"merge userIDs", &SearchParams{userIDs: []int{123, 456}}, args{[]int{789, 456}}, []int{123, 456, 789}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddUserIDs(tt.args.users)
			result := tt.params.UserIDs()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddWords(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.users, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchParams_ClearUserIDs(t *testing.T) {
	params := &SearchParams{userIDs: []int{123}}
	params.ClearUserIDs()
	if params.UserIDs() != nil {
		t.Errorf("SearchParams.ClearUserIDs() should change UserIDs be nil, but %v", params.UserIDs())
	}
}

func TestSearchParams_queryFromUserID(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no userids, no query", &SearchParams{}, makeValues([][2]string{})},
		{`userIDs:[123], userid=123`, &SearchParams{userIDs: []int{123}}, makeValues([][2]string{{"userid", "123"}})},
		{`userIDs:[123, 456], userid=123-456`, &SearchParams{userIDs: []int{123, 456}}, makeValues([][2]string{{"userid", "123-456"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromUserID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_IsR15(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default", &SearchParams{}, false},
		{"isR15", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsR15: true}}, true},
		{"not isR15", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsR15: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsR15(); got != tt.want {
				t.Errorf("SearchParams.IsR15() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_IsNotR15(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"NotR15", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotR15: true}}, true},
		{"not notR15", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotR15: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsNotR15(); got != tt.want {
				t.Errorf("SearchParams.IsNotR15() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetIsR15(t *testing.T) {
	type args struct {
		r15 bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsR15(tt.args.r15)
			result := tt.params.IsR15()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsR15(%v) result %v, want %v", tt.args.r15, result, tt.want)
			}
		})
	}
	t.Run("reset IsNotR15", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsNotR15(true)
		params.SetIsR15(true)
		if params.IsNotR15() != false {
			t.Errorf("SearchParams.SetIsR15(true) set IsNotR15 false, but %v", params.IsNotR15())
		}
	})
}

func TestSearchParams_SetIsNotR15(t *testing.T) {
	type args struct {
		r15 bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsNotR15(tt.args.r15)
			result := tt.params.IsNotR15()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsR15(%v) result %v, want %v", tt.args.r15, result, tt.want)
			}
		})
	}
	t.Run("reset IsNotR15", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsR15(true)
		params.SetIsNotR15(true)
		if params.IsR15() != false {
			t.Errorf("SearchParams.SetIsNotR15(true) set IsR15 false, but %v", params.IsR15())
		}
	})
}

func TestSearchParams_IsBL(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"Bl", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsBL: true}}, true},
		{"not bl", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsBL: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsBL(); got != tt.want {
				t.Errorf("SearchParams.IsBL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_IsNotBL(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"NotBl", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotBL: true}}, true},
		{"not NotBL", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotBL: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsNotBL(); got != tt.want {
				t.Errorf("SearchParams.IsNotBL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetIsBL(t *testing.T) {
	type args struct {
		bl bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsBL(tt.args.bl)
			result := tt.params.IsBL()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsBL(%v) result %v, want %v", tt.args.bl, result, tt.want)
			}
		})
	}
	t.Run("reset IsNotBL", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsNotBL(true)
		params.SetIsBL(true)
		if params.IsNotBL() != false {
			t.Errorf("SearchParams.SetIsBL(true) set IsNotBL false, but %v", params.IsNotBL())
		}
	})
}

func TestSearchParams_SetIsNotBL(t *testing.T) {
	type args struct {
		bl bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsNotBL(tt.args.bl)
			result := tt.params.IsNotBL()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsNotBL(%v) result %v, want %v", tt.args.bl, result, tt.want)
			}
		})
	}
	t.Run("reset IsBL", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsBL(true)
		params.SetIsNotBL(true)
		if params.IsBL() != false {
			t.Errorf("SearchParams.SetIsNotBL(true) set IsBL false, but %v", params.IsBL())
		}
	})
}

func TestSearchParams_IsGL(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"Gl", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsGL: true}}, true},
		{"not gl", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsGL: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsGL(); got != tt.want {
				t.Errorf("SearchParams.IsGL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_IsNotGL(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"NotGl", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotGL: true}}, true},
		{"not NotGL", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotGL: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsNotGL(); got != tt.want {
				t.Errorf("SearchParams.IsNotGL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetIsGL(t *testing.T) {
	type args struct {
		gl bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsGL(tt.args.gl)
			result := tt.params.IsGL()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsGL(%v) result %v, want %v", tt.args.gl, result, tt.want)
			}
		})
	}
	t.Run("reset IsNotGL", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsNotGL(true)
		params.SetIsGL(true)
		if params.IsNotGL() != false {
			t.Errorf("SearchParams.SetIsGL(true) set IsNotGL false, but %v", params.IsNotGL())
		}
	})
}

func TestSearchParams_SetIsNotGL(t *testing.T) {
	type args struct {
		gl bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsNotGL(tt.args.gl)
			result := tt.params.IsNotGL()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsNotGL(%v) result %v, want %v", tt.args.gl, result, tt.want)
			}
		})
	}
	t.Run("reset IsGL", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsGL(true)
		params.SetIsNotGL(true)
		if params.IsGL() != false {
			t.Errorf("SearchParams.SetIsNotGL(true) set IsGL false, but %v", params.IsGL())
		}
	})
}

func TestSearchParams_IsZankoku(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"Zankoku", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsZankoku: true}}, true},
		{"not zankoku", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsZankoku: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsZankoku(); got != tt.want {
				t.Errorf("SearchParams.IsZankoku() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_IsNotZankoku(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"not Zankoku", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotZankoku: true}}, true},
		{"not NotZankoku", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotZankoku: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsNotZankoku(); got != tt.want {
				t.Errorf("SearchParams.IsNotZankoku() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetIsZankoku(t *testing.T) {
	type args struct {
		zankoku bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsZankoku(tt.args.zankoku)
			result := tt.params.IsZankoku()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsZANKOKU(%v) result %v, want %v", tt.args.zankoku, result, tt.want)
			}
		})
	}
	t.Run("reset IsNotZankoku", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsNotZankoku(true)
		params.SetIsZankoku(true)
		if params.IsNotZankoku() != false {
			t.Errorf("SearchParams.SetIsZANKOKU(true) set IsNotZANKOKU false, but %v", params.IsNotZankoku())
		}
	})
}

func TestSearchParams_SetIsNotZankoku(t *testing.T) {
	type args struct {
		zankoku bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsNotZankoku(tt.args.zankoku)
			result := tt.params.IsNotZankoku()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsNotZankoku(%v) result %v, want %v", tt.args.zankoku, result, tt.want)
			}
		})
	}
	t.Run("reset IsZankoku", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsZankoku(true)
		params.SetIsNotZankoku(true)
		if params.IsZankoku() != false {
			t.Errorf("SearchParams.SetIsNotZanoku(true) set IsZankoku false, but %v", params.IsZankoku())
		}
	})
}

func TestSearchParams_IsTensei(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"Tensei", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsTensei: true}}, true},
		{"not tensei", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsTensei: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsTensei(); got != tt.want {
				t.Errorf("SearchParams.IsTensei() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_IsNotTensei(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"NotTensei", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotTensei: true}}, true},
		{"not NotTensei", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotTensei: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsNotTensei(); got != tt.want {
				t.Errorf("SearchParams.IsNotTensei() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetIsTensei(t *testing.T) {
	type args struct {
		tensei bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsTensei(tt.args.tensei)
			result := tt.params.IsTensei()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsTensei(%v) result %v, want %v", tt.args.tensei, result, tt.want)
			}
		})
	}
	t.Run("reset IsNotTensei", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsNotTensei(true)
		params.SetIsTensei(true)
		if params.IsNotTensei() != false {
			t.Errorf("SearchParams.SetIsTensei(true) set IsNotTensei false, but %v", params.IsNotTensei())
		}
	})
}

func TestSearchParams_SetIsNotTensei(t *testing.T) {
	type args struct {
		tensei bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsNotTensei(tt.args.tensei)
			result := tt.params.IsNotTensei()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsNotTensei(%v) result %v, want %v", tt.args.tensei, result, tt.want)
			}
		})
	}
	t.Run("reset IsTensei", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsTensei(true)
		params.SetIsNotTensei(true)
		if params.IsTensei() != false {
			t.Errorf("SearchParams.SetIsNotZanoku(true) set IsTensei false, but %v", params.IsTensei())
		}
	})
}

func TestSearchParams_IsTenni(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"Tenni", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsTenni: true}}, true},
		{"not tenni", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsTenni: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsTenni(); got != tt.want {
				t.Errorf("SearchParams.IsTenni() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestSearchParams_IsNotTenni(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"NotTenni", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotTenni: true}}, true},
		{"not NotTenni", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotTenni: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsNotTenni(); got != tt.want {
				t.Errorf("SearchParams.IsNotTenni() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetIsTenni(t *testing.T) {
	type args struct {
		tenni bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsTenni(tt.args.tenni)
			result := tt.params.IsTenni()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsTenni(%v) result %v, want %v", tt.args.tenni, result, tt.want)
			}
		})
	}
	t.Run("reset IsNotTenni", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsNotTenni(true)
		params.SetIsTenni(true)
		if params.IsNotTenni() != false {
			t.Errorf("SearchParams.SetIsTenni(true) set IsNotTenni false, but %v", params.IsNotTenni())
		}
	})
}

func TestSearchParams_SetIsNotTenni(t *testing.T) {
	type args struct {
		tenni bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsNotTenni(tt.args.tenni)
			result := tt.params.IsNotTenni()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsNotTenni(%v) result %v, want %v", tt.args.tenni, result, tt.want)
			}
		})
	}
	t.Run("reset IsTenni", func(t *testing.T) {
		params := NewSearchParams()
		params.SetIsTenni(true)
		params.SetIsNotTenni(true)
		if params.IsTenni() != false {
			t.Errorf("SearchParams.SetIsNotZanoku(true) set IsTenni false, but %v", params.IsTenni())
		}
	})
}

func TestSearchParams_IsTT(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   bool
	}{
		{"default is false", &SearchParams{}, false},
		{"TT", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsTT: true}}, true},
		{"not TT", &SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsTT: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsTT(); got != tt.want {
				t.Errorf("SearchParams.IsTT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetIsTT(t *testing.T) {
	type args struct {
		tt bool
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   bool
	}{
		{"set true", NewSearchParams(), args{true}, true},
		{"set false", NewSearchParams(), args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetIsTT(tt.args.tt)
			result := tt.params.IsTT()
			if result != tt.want {
				t.Errorf("SearchParams.SetIsTT(%v) result %v, want %v", tt.args.tt, result, tt.want)
			}
		})
	}
}

func TestSearchParams_queryFromRequiredKeywords(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"default", NewSearchParams(), makeValues([][2]string{})},
		{"isr15",
			&SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsR15: true}},
			makeValues([][2]string{{"isr15", "1"}})},
		{"is not r15",
			&SearchParams{requiredKeywordFlags: map[requiredKeyword]bool{requiredKeywordIsNotR15: true}},
			makeValues([][2]string{{"notr15", "1"}})},
		// TODO: more flags
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromRequiredKeywords(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromRequiredKeywords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_Length(t *testing.T) {
	len := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"lengths.single", &SearchParams{lengths: minmaxPair{mmpSingle: 42}}, &len},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.Length(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetLength(t *testing.T) {
	type args struct {
		len int
	}
	len := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid length", NewSearchParams(), args{42}, &len},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetLength(tt.args.len)
			result := tt.params.Length()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetLength(%v) result %v, want %v", tt.args.len, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearLength(t *testing.T) {
	params := NewSearchParams()
	params.SetLength(123)
	params.ClearLength()
	if params.Length() != nil {
		t.Errorf("SearchParams.ClearLength() should change Length be nil, but %v", params.Length())
	}
}

func TestSearchParams_MinLength(t *testing.T) {
	len := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"lengths.single", &SearchParams{lengths: minmaxPair{mmpMin: 42}}, &len},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.MinLength(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.MinLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetMinLength(t *testing.T) {
	type args struct {
		minLen int
	}
	len := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid length", NewSearchParams(), args{42}, &len},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetMinLength(tt.args.minLen)
			result := tt.params.MinLength()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetMinLength(%v) result %v, want %v", tt.args.minLen, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearMinLength(t *testing.T) {
	params := NewSearchParams()
	params.SetMinLength(123)
	params.ClearMinLength()
	if params.MinLength() != nil {
		t.Errorf("SearchParams.ClearMinLength() should change Length be nil, but %v", params.MinLength())
	}
}

func TestSearchParams_MaxLength(t *testing.T) {
	len := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"lengths.single", &SearchParams{lengths: minmaxPair{mmpMax: 42}}, &len},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.MaxLength(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.MaxLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetMaxLength(t *testing.T) {
	type args struct {
		maxLen int
	}
	len := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid length", NewSearchParams(), args{42}, &len},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetMaxLength(tt.args.maxLen)
			result := tt.params.MaxLength()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetMaxLength(%v) result %v, want %v", tt.args.maxLen, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearMaxLength(t *testing.T) {
	params := NewSearchParams()
	params.SetMaxLength(123)
	params.ClearMaxLength()
	if params.MaxLength() != nil {
		t.Errorf("SearchParams.ClearMaxLength() should change Length be nil, but %v", params.MaxLength())
	}
}

func TestSearchParams_queryFromLength(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no length, no query", NewSearchParams(), makeValues([][2]string{})},
		{"length only, length=len", &SearchParams{lengths: minmaxPair{mmpSingle: 42}}, makeValues([][2]string{{"length", "42"}})},
		{"min only, length=min-", &SearchParams{lengths: minmaxPair{mmpMin: 123}}, makeValues([][2]string{{"length", "123-"}})},
		{"max only, length=-max", &SearchParams{lengths: minmaxPair{mmpMax: 456}}, makeValues([][2]string{{"length", "-456"}})},
		{"min/max both, length=min-max", &SearchParams{lengths: minmaxPair{mmpMin: 123, mmpMax: 456}}, makeValues([][2]string{{"length", "123-456"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromLength(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_Kaiwaritu(t *testing.T) {
	kaiwa := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"kaiwaritu.single", &SearchParams{kaiwaritus: minmaxPair{mmpSingle: 42}}, &kaiwa},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.Kaiwaritu(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.Kaiwaritu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetKaiwaritu(t *testing.T) {
	type args struct {
		kaiwa int
	}
	kaiwa := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid kaiwaritu", NewSearchParams(), args{42}, &kaiwa},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetKaiwaritu(tt.args.kaiwa)
			result := tt.params.Kaiwaritu()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetKaiwaritu(%v) result %v, want %v", tt.args.kaiwa, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearKaiwaritu(t *testing.T) {
	params := NewSearchParams()
	params.SetKaiwaritu(42)
	params.ClearKaiwaritu()
	if params.Kaiwaritu() != nil {
		t.Errorf("SearchParams.ClearKaiwaritu() should change Kaiwaritu be nil, but %v", params.Kaiwaritu())
	}
}

func TestSearchParams_MinKaiwaritu(t *testing.T) {
	kaiwa := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"kaiwaritus.single", &SearchParams{kaiwaritus: minmaxPair{mmpMin: 42}}, &kaiwa},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.MinKaiwaritu(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.MinKaiwaritu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetMinKaiwaritu(t *testing.T) {
	type args struct {
		minKaiwa int
	}
	ritu := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid kaiwaritu", NewSearchParams(), args{42}, &ritu},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetMinKaiwaritu(tt.args.minKaiwa)
			result := tt.params.MinKaiwaritu()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetMinKaiwaritu(%v) result %v, want %v", tt.args.minKaiwa, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearMinKaiwaritu(t *testing.T) {
	params := NewSearchParams()
	params.SetMinKaiwaritu(123)
	params.ClearMinKaiwaritu()
	if params.MinKaiwaritu() != nil {
		t.Errorf("SearchParams.ClearMinKaiwaritu() should change MinKaiwaritu be nil, but %v", params.MinKaiwaritu())
	}
}

func TestSearchParams_MaxKaiwaritu(t *testing.T) {
	kaiwa := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"kaiwaritus.single", &SearchParams{kaiwaritus: minmaxPair{mmpMax: 42}}, &kaiwa},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.MaxKaiwaritu(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.MaxKaiwaritu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetMaxKaiwaritu(t *testing.T) {
	type args struct {
		maxKaiwa int
	}
	kaiwa := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid kaiwaritu", NewSearchParams(), args{42}, &kaiwa},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetMaxKaiwaritu(tt.args.maxKaiwa)
			result := tt.params.MaxKaiwaritu()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetMaxKaiwaritu(%v) result %v, want %v", tt.args.maxKaiwa, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearMaxKaiwaritu(t *testing.T) {
	params := NewSearchParams()
	params.SetMaxKaiwaritu(123)
	params.ClearMaxKaiwaritu()
	if params.MaxKaiwaritu() != nil {
		t.Errorf("SearchParams.ClearMaxKaiwaritu() should change MaxKaiwaritu be nil, but %v", params.MaxKaiwaritu())
	}
}

func TestSearchParams_queryFromKaiwaritu(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no kaiwaritu, no query", NewSearchParams(), makeValues([][2]string{})},
		{"kaiwaritu only, kaiwaritu=len", &SearchParams{kaiwaritus: minmaxPair{mmpSingle: 42}}, makeValues([][2]string{{"kaiwaritu", "42"}})},
		{"min only, kaiwaritu=min-", &SearchParams{kaiwaritus: minmaxPair{mmpMin: 123}}, makeValues([][2]string{{"kaiwaritu", "123-"}})},
		{"max only, kaiwaritu=-max", &SearchParams{kaiwaritus: minmaxPair{mmpMax: 456}}, makeValues([][2]string{{"kaiwaritu", "-456"}})},
		{"min/max both, kaiwaritu=min-max", &SearchParams{kaiwaritus: minmaxPair{mmpMin: 123, mmpMax: 456}}, makeValues([][2]string{{"kaiwaritu", "123-456"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromKaiwaritu(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromKaiwaritu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_Sasie(t *testing.T) {
	sasie := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"sasie.single", &SearchParams{sasies: minmaxPair{mmpSingle: 42}}, &sasie},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.Sasie(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.Sasie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetSasie(t *testing.T) {
	type args struct {
		sasie int
	}
	sasie := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid kaiwaritu", NewSearchParams(), args{42}, &sasie},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetSasie(tt.args.sasie)
		})
	}
}

func TestSearchParams_ClearSasie(t *testing.T) {
	params := NewSearchParams()
	params.SetSasie(42)
	params.ClearSasie()
	if params.Sasie() != nil {
		t.Errorf("SearchParams.ClearSasie() should change Sasie be nil, but %v", params.Sasie())
	}
}

func TestSearchParams_MinSasie(t *testing.T) {
	sasie := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"sasie.single", &SearchParams{sasies: minmaxPair{mmpMin: 42}}, &sasie},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.MinSasie(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.MinSasie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetMinSasie(t *testing.T) {
	type args struct {
		minSasie int
	}
	ritu := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid sasie", NewSearchParams(), args{42}, &ritu},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetMinSasie(tt.args.minSasie)
			result := tt.params.MinSasie()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetMinSasie(%v) result %v, want %v", tt.args.minSasie, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearMinSasie(t *testing.T) {
	params := NewSearchParams()
	params.SetMinSasie(123)
	params.ClearMinSasie()
	if params.MinSasie() != nil {
		t.Errorf("SearchParams.ClearMinSasie() should change MinSasie be nil, but %v", params.MinSasie())
	}
}

func TestSearchParams_MaxSasie(t *testing.T) {
	sasie := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"sasie.single", &SearchParams{sasies: minmaxPair{mmpMax: 42}}, &sasie},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.MaxSasie(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.MaxSasie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetMaxSasie(t *testing.T) {
	type args struct {
		maxSasie int
	}
	sasie := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid sasie", NewSearchParams(), args{42}, &sasie},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetMaxSasie(tt.args.maxSasie)
			result := tt.params.MaxSasie()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetMaxSasie(%v) result %v, want %v", tt.args.maxSasie, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearMaxSasie(t *testing.T) {
	params := NewSearchParams()
	params.SetMaxSasie(123)
	params.ClearMaxSasie()
	if params.MaxSasie() != nil {
		t.Errorf("SearchParams.ClearMaxSasie() should change MaxSasie be nil, but %v", params.MaxSasie())
	}
}

func TestSearchParams_queryFromSasie(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no sasie, no query", NewSearchParams(), makeValues([][2]string{})},
		{"sasie only, sasie=len", &SearchParams{sasies: minmaxPair{mmpSingle: 42}}, makeValues([][2]string{{"sasie", "42"}})},
		{"min only, sasie=min-", &SearchParams{sasies: minmaxPair{mmpMin: 12}}, makeValues([][2]string{{"sasie", "12-"}})},
		{"max only, sasie=-max", &SearchParams{sasies: minmaxPair{mmpMax: 34}}, makeValues([][2]string{{"sasie", "-34"}})},
		{"min/max both, sasie=min-max", &SearchParams{sasies: minmaxPair{mmpMin: 12, mmpMax: 34}}, makeValues([][2]string{{"sasie", "12-34"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromSasie(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromSasie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_ReadTime(t *testing.T) {
	readTime := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"time.single", &SearchParams{readTimes: minmaxPair{mmpSingle: 42}}, &readTime},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.ReadTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.ReadTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetReadTime(t *testing.T) {
	type args struct {
		rt int
	}
	readTime := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid read time", NewSearchParams(), args{42}, &readTime},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetReadTime(tt.args.rt)
			result := tt.params.ReadTime()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetReadTime(%v) result %v, want %v", tt.args.rt, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearReadTime(t *testing.T) {
	params := NewSearchParams()
	params.SetReadTime(42)
	params.ClearReadTime()
	if params.ReadTime() != nil {
		t.Errorf("SearchParams.ClearReadTime() should change ReadTime be nil, but %v", params.ReadTime())
	}
}

func TestSearchParams_MinReadTime(t *testing.T) {
	readTime := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"readtime.single", &SearchParams{readTimes: minmaxPair{mmpMin: 42}}, &readTime},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.MinReadTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.MinReadTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetMinReadTime(t *testing.T) {
	type args struct {
		minRt int
	}
	readTime := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid time", NewSearchParams(), args{42}, &readTime},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetMinReadTime(tt.args.minRt)
			result := tt.params.MinReadTime()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetMaxSasie(%v) result %v, want %v", tt.args.minRt, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearMinReadTime(t *testing.T) {
	params := NewSearchParams()
	params.SetMinReadTime(123)
	params.ClearMinReadTime()
	if params.MinReadTime() != nil {
		t.Errorf("SearchParams.ClearMinReadTime() should change MinReadTime be nil, but %v", params.MinReadTime())
	}
}

func TestSearchParams_MaxReadTime(t *testing.T) {
	readTime := 42
	tests := []struct {
		name   string
		params *SearchParams
		want   *int
	}{
		{"default", &SearchParams{}, nil},
		// TODO: too internal
		{"readtime.single", &SearchParams{readTimes: minmaxPair{mmpMax: 42}}, &readTime},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.MaxReadTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.MaxReadTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetMaxReadTime(t *testing.T) {
	type args struct {
		maxRt int
	}
	readTime := 42
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   *int
	}{
		{"set valid sasie", NewSearchParams(), args{42}, &readTime},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetMaxReadTime(tt.args.maxRt)
			result := tt.params.MaxReadTime()
			if *result != *tt.want {
				t.Errorf("SearchParams.SetMaxReadTime(%v) result %v, want %v", tt.args.maxRt, *result, *tt.want)
			}
		})
	}
}

func TestSearchParams_ClearMaxReadTime(t *testing.T) {
	params := NewSearchParams()
	params.SetMaxReadTime(123)
	params.ClearMaxReadTime()
	if params.MaxReadTime() != nil {
		t.Errorf("SearchParams.ClearMaxReadTime() should change MaxReadTime be nil, but %v", params.MaxReadTime())
	}
}

func TestSearchParams_queryFromReadTime(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no readTime, no query", NewSearchParams(), makeValues([][2]string{})},
		{"readTime only, time=len", &SearchParams{readTimes: minmaxPair{mmpSingle: 42}}, makeValues([][2]string{{"time", "42"}})},
		{"min only, time=min-", &SearchParams{readTimes: minmaxPair{mmpMin: 12}}, makeValues([][2]string{{"time", "12-"}})},
		{"max only, time=-max", &SearchParams{readTimes: minmaxPair{mmpMax: 34}}, makeValues([][2]string{{"time", "-34"}})},
		{"min/max both, time=min-max", &SearchParams{readTimes: minmaxPair{mmpMin: 12, mmpMax: 34}}, makeValues([][2]string{{"time", "12-34"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromReadTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromReadTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_NCodes(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []string
	}{
		{"no ncodes", &SearchParams{}, nil},
		{"nocdes", &SearchParams{ncodes: []string{"abc", "def"}}, []string{"abc", "def"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.NCodes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.NCodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddNCodes(t *testing.T) {
	type args struct {
		ncodes []string
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []string
	}{
		{"ncodes + []", &SearchParams{ncodes: []string{"some", "ncodes"}}, args{[]string{}}, []string{"some", "ncodes"}},
		{"[] + ncodes", &SearchParams{ncodes: []string{}}, args{[]string{"some", "ncodes"}}, []string{"some", "ncodes"}},
		{"ncodes + ncodes", &SearchParams{ncodes: []string{"some", "ncodes"}}, args{[]string{"more", "phrase"}}, []string{"some", "ncodes", "more", "phrase"}},
		{"merge ncodes", &SearchParams{ncodes: []string{"some", "ncodes"}}, args{[]string{"more", "ncodes"}}, []string{"some", "ncodes", "more"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddNCodes(tt.args.ncodes)
		})
	}
}

func TestSearchParams_ClearNCodes(t *testing.T) {
	params := &SearchParams{}
	ncodes := []string{"first", "ncode"}
	params.AddNCodes(ncodes)
	params.ClearNCodes()
	if params.NCodes() != nil {
		t.Errorf("SearchParams.ClearNcodes() should change ncodes be nil, but %v", params.NCodes())
	}
}

func TestSearchParams_queryFromNCode(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no ncodes, no query", &SearchParams{}, makeValues([][2]string{})},
		{`ncode:["ncode1"], ncode=ncode`, &SearchParams{ncodes: []string{"ncode1"}}, makeValues([][2]string{{"ncode", "ncode1"}})},
		{`ncode:["ncode1", "ncode2"], ncode=ncode-ncode2`,
			&SearchParams{ncodes: []string{"ncode1", "ncode2"}},
			makeValues([][2]string{{"ncode", "ncode1-ncode2"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromNCode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromNCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_NovelState(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   NovelState
	}{
		{"default order is All", &SearchParams{}, NovelStateAll},
		{"get novelstate", &SearchParams{state: NovelStateShortStory}, NovelStateShortStory},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.NovelState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.NovelState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetNovelState(t *testing.T) {
	type args struct {
		state NovelState
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   NovelState
	}{
		{"set NovelState should change .state", &SearchParams{}, args{NovelStateRensaiAll}, NovelStateRensaiAll},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetNovelState(tt.args.state)
			result := tt.params.NovelState()
			if result != tt.want {
				t.Errorf("SearchParams.SetNovelState(%v) should be change .state %v, but %v", tt.args.state, tt.want, result)
			}
		})
	}
}

func TestSearchParams_ClearNovelState(t *testing.T) {
	params := &SearchParams{}
	params.SetNovelState(NovelStateShortStory)
	params.ClearNovelState()
	if params.NovelState() != NovelStateAll {
		t.Errorf("SearchParams.ClearNovelState() should change offset be 0, but %v", params.NovelState())
	}
}

func TestSearchParams_queryFromNovelState(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no state, no query", &SearchParams{}, makeValues([][2]string{})},
		{"state:All, no query", &SearchParams{state: NovelStateAll}, makeValues([][2]string{})},
		{"state:Rensai, type=re", &SearchParams{state: NovelStateRensaiAll}, makeValues([][2]string{{"type", "re"}})},
		{"state:Unknown, no query", &SearchParams{state: 2000}, makeValues([][2]string{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromNovelState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromNovelState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_Buntais(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   []Buntai
	}{
		{"no buntai, empty array", &SearchParams{}, nil},
		{"buntais", &SearchParams{buntais: []Buntai{BuntaiNoIndentManyEmptyLines}}, []Buntai{BuntaiNoIndentManyEmptyLines}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.Buntais(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.Buntais() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_AddBuntais(t *testing.T) {
	type args struct {
		buntais []Buntai
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   []Buntai
	}{
		{"empty array should not change", &SearchParams{}, args{[]Buntai{}}, []Buntai{}},
		{"args include All should clear fields",
			&SearchParams{buntais: []Buntai{BuntaiIndentAverageEmptyLines}},
			args{[]Buntai{BuntaiAll, BuntaiIndentManyEmptyLines}}, nil},
		{"add output fields",
			&SearchParams{buntais: []Buntai{BuntaiIndentAverageEmptyLines}},
			args{[]Buntai{BuntaiIndentManyEmptyLines}}, []Buntai{BuntaiIndentAverageEmptyLines, BuntaiIndentManyEmptyLines}},
		{"merge output fields",
			&SearchParams{buntais: []Buntai{BuntaiIndentAverageEmptyLines}},
			args{[]Buntai{BuntaiNoIndentAveraegEmpytLines, BuntaiIndentAverageEmptyLines}},
			[]Buntai{BuntaiIndentAverageEmptyLines, BuntaiNoIndentAveraegEmpytLines}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddBuntais(tt.args.buntais)
			result := tt.params.Buntais()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddBuntais(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.buntais, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchParams_ClearBuntais(t *testing.T) {
	params := &SearchParams{}
	params.AddBuntais([]Buntai{BuntaiIndentAverageEmptyLines})
	params.ClearBuntais()
	if params.Buntais() != nil {
		t.Errorf("SearchParams.ClearBuntais() should change buntais be nil, but %v", params.Buntais())
	}
}

func TestSearchParams_queryFromBuntai(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no buntai, no query", &SearchParams{}, makeValues([][2]string{})},
		{"BuntaiIndentAverageEmptyLines should buntai=6",
			&SearchParams{buntais: []Buntai{BuntaiIndentAverageEmptyLines}},
			makeValues([][2]string{{"buntai", "6"}})},
		{"NoIndentAveraegEmpytLines and NoIndentManyEmptyLines should buntai=2-1",
			&SearchParams{buntais: []Buntai{BuntaiNoIndentAveraegEmpytLines, BuntaiNoIndentManyEmptyLines}},
			makeValues([][2]string{{"buntai", "2-1"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromBuntai(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromBuntai() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_StopState(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   StopState
	}{
		{"default stop is all", &SearchParams{}, StopStateAll},
		{"get stopstate", &SearchParams{stopState: StopStateExclude}, StopStateExclude},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.StopState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.StopState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetStopState(t *testing.T) {
	type args struct {
		state StopState
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   StopState
	}{
		{"set orderitem should change .order", &SearchParams{}, args{StopStateOnly}, StopStateOnly},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetStopState(tt.args.state)
			result := tt.params.StopState()
			if result != tt.want {
				t.Errorf("SearchParams.SetStopState(%v) should be change .order %v, but %v", tt.args.state, tt.want, result)
			}
		})
	}
}

func TestSearchParams_ClearStopState(t *testing.T) {
	params := &SearchParams{}
	params.SetStopState(StopStateExclude)
	params.ClearStopState()
	if params.StopState() != StopStateAll {
		t.Errorf("SearchParams.ClearStopState() should change StopState be ALL, but %v", params.StopState())
	}
}

func TestSearchParams_queryFromStop(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no stopState, no query", &SearchParams{}, makeValues([][2]string{})},
		{"stopState All, no query", &SearchParams{stopState: StopStateAll}, makeValues([][2]string{})},
		{"stopState:Exclude, stopState=1", &SearchParams{stopState: StopStateExclude}, makeValues([][2]string{{"stop", "1"}})},
		{"stopState:Only, stopState=2", &SearchParams{stopState: StopStateOnly}, makeValues([][2]string{{"stop", "2"}})},
		// {"stopState:Unknown, no query", &SearchParams{stopState: 2000}, makeValues([][2]string{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromStop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromStop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_PickupState(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   PickupState
	}{
		{"default pickup is none", &SearchParams{}, PickupStateNone},
		{"get pickup", &SearchParams{pickupState: PickupStatePickup}, PickupStatePickup},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.PickupState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.PickupState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetPickupState(t *testing.T) {
	type args struct {
		pickup PickupState
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   PickupState
	}{
		{"set pickup should change .pickup", &SearchParams{}, args{PickupStatePickup}, PickupStatePickup},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetPickupState(tt.args.pickup)
			result := tt.params.PickupState()
			if result != tt.want {
				t.Errorf("SearchParams.SetPickupState(%v) should be change .pickup %v, but %v", tt.args.pickup, tt.want, result)
			}
		})
	}
}

func TestSearchParams_ClearPickupState(t *testing.T) {
	params := &SearchParams{}
	params.SetPickupState(PickupStatePickup)
	params.ClearPickupState()
	if params.PickupState() != PickupStateNone {
		t.Errorf("SearchParams.ClearPickupState() should change pickup be None, but %v", params.PickupState())
	}
}

func TestSearchParams_queryFromPickup(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no order, no query", &SearchParams{}, makeValues([][2]string{})},
		{"pickupState None, no query", &SearchParams{pickupState: PickupStateNone}, makeValues([][2]string{})},
		{"pickupState:Pickup, ispickup=1", &SearchParams{pickupState: PickupStatePickup}, makeValues([][2]string{{"ispickup", "1"}})},
		{"pickupState:Not, ispickup=0", &SearchParams{pickupState: PickupStateNot}, makeValues([][2]string{{"ispickup", "0"}})},
		{"pickupState:Unknown, no query", &SearchParams{pickupState: 2000}, makeValues([][2]string{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromPickup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromPickup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_LastUpType(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchParams
		want   LastUpType
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.LastUpType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.LastUpType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_LastUpTimeStamps(t *testing.T) {
	t1 := time.Date(2019, 2, 3, 4, 5, 6, 7, time.UTC)
	t2 := time.Date(2020, 2, 3, 4, 5, 6, 7, time.UTC)
	tests := []struct {
		name   string
		params *SearchParams
		want   [2]time.Time
	}{
		{"default", &SearchParams{}, [2]time.Time{}},
		{"not timestampType, not timestamps", &SearchParams{lastUp: LastUpTypeLastMonth, lastUpStart: t1, lastUpEnd: t2}, [2]time.Time{}},
		{"timestampType, not timestamps", &SearchParams{lastUp: LastUpTypeTimeStamp, lastUpStart: t1, lastUpEnd: t2}, [2]time.Time{t1, t2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.LastUpTimeStamps(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.LastUpTimeStamps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchParams_SetLastUp(t *testing.T) {
	type args struct {
		ltype LastUpType
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   LastUpType
	}{
		{"set orderitem should change .order", &SearchParams{}, args{ltype: LastUpTypeLastMonth}, LastUpTypeLastMonth},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetLastUp(tt.args.ltype)
			result := tt.params.LastUpType()
			if result != tt.want {
				t.Errorf("SearchParams.SetLastUP(%v) should be change .lastup %v, but %v", tt.args.ltype, tt.want, result)
			}
		})
	}
}

func TestSearchParams_SetLastUpTerm(t *testing.T) {
	type args struct {
		start time.Time
		end   time.Time
	}
	type wants struct {
		start time.Time
		end   time.Time
		ltype LastUpType
	}
	tests := []struct {
		name   string
		params *SearchParams
		args   args
		want   wants
	}{
		{"set orderitem should change .order", &SearchParams{},
			args{start: time.Date(2019, 1, 2, 3, 4, 5, 6, time.UTC), end: time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)},
			wants{start: time.Date(2019, 1, 2, 3, 4, 5, 6, time.UTC), end: time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC), ltype: LastUpTypeTimeStamp},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.SetLastUpTerm(tt.args.start, tt.args.end)
			resultType := tt.params.LastUpType()
			resultTerm := tt.params.LastUpTimeStamps()
			if resultTerm[0] != tt.want.start || resultTerm[1] != tt.want.end || resultType != tt.want.ltype {
				t.Errorf("SearchParams.SetLastUpTerm(%v, %v) should be change .ltype %v .term %v, %v, but %v, %v, %v",
					tt.args.start, tt.args.end, LastUpTypeTimeStamp, tt.want.start, tt.want.end, resultType, resultTerm[0], resultTerm[1])
			}
		})
	}
}

func TestSearchParams_ClearLastUp(t *testing.T) {
	params := &SearchParams{lastUp: LastUpTypeLastMonth}
	params.ClearLastUp()
	if params.LastUpType() != LastUpTypeNone {
		t.Errorf("SearchParams.ClearLastUp() should change LastUp be nil, but %v", params.LastUpType())
	}
}

func TestSearchParams_queryFromLastUp(t *testing.T) {
	t1 := time.Date(2019, 1, 2, 3, 4, 5, 6, time.UTC)
	t2 := time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)
	tests := []struct {
		name   string
		params *SearchParams
		want   url.Values
	}{
		{"no lastup, no query", &SearchParams{}, makeValues([][2]string{})},
		{`lastup:ThisWeek, lastup=thisweek`, &SearchParams{lastUp: LastUpTypeThisWeek}, makeValues([][2]string{{"lastup", "thisweek"}})},
		{`lastup:TimeStamp, lastup=T1-T2`,
			&SearchParams{lastUp: LastUpTypeTimeStamp, lastUpStart: t1, lastUpEnd: t2},
			makeValues([][2]string{{"lastup", fmt.Sprintf("%d-%d", t1.Unix(), t2.Unix())}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromLastUp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchParams.queryFromLastUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_minmaxPair_val(t *testing.T) {
	type args struct {
		k mmpKey
	}
	tests := []struct {
		name string
		mmp  minmaxPair
		args args
		want *int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mmp.val(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("minmaxPair.val() = %v, want %v", got, tt.want)
			}
		})
	}
}
