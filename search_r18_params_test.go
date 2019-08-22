package narrow

import (
	"net/url"
	"reflect"
	"testing"
)

func TestNewSearchR18Params(t *testing.T) {
	tests := []struct {
		name string
		want *SearchR18Params
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSearchR18Params(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSearchR18Params() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchR18Params_ToURL(t *testing.T) {
	tests := []struct {
		name    string
		params  *SearchR18Params
		want    *url.URL
		wantErr bool
	}{
		{"all default params",
			NewSearchR18Params(),
			parseURL("https://api.syosetu.com/novel18api/api/?out=json"), false},
		{"with st(offset)",
			&SearchR18Params{SearchParams: SearchParams{offset: 42}},
			parseURL("https://api.syosetu.com/novel18api/api/?out=json&st=42"), false},
		{"with limit",
			&SearchR18Params{SearchParams: SearchParams{limit: 42}},
			parseURL("https://api.syosetu.com/novel18api/api/?out=json&limit=42"), false},
		{"with st(offset), limit",
			&SearchR18Params{SearchParams: SearchParams{offset: 42, limit: 30}},
			parseURL("https://api.syosetu.com/novel18api/api/?out=json&st=42&limit=30"), false},
		{"with output field All",
			&SearchR18Params{SearchParams: SearchParams{outputFields: []OutputField{OutputFieldAll}}},
			parseURL("https://api.syosetu.com/novel18api/api/?out=json"), false},
		{"with output field Title",
			&SearchR18Params{SearchParams: SearchParams{outputFields: []OutputField{OutputFieldTitle}}},
			parseURL("https://api.syosetu.com/novel18api/api/?out=json&of=t"), false},
		{"with output field Title, NCode",
			&SearchR18Params{SearchParams: SearchParams{outputFields: []OutputField{OutputFieldTitle, OutputFieldNCode}}},
			parseURL("https://api.syosetu.com/novel18api/api/?out=json&of=t-n"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.params.ToURL()
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchR18Params.ToURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// should not change path
			if got.Path != tt.want.Path {
				t.Errorf("SearchR18Params.ToURL().Path = %v, want.Path %v", got.Path, tt.want.Path)
				return
			}

			// query part
			gr, wr := got.RawQuery, tt.want.RawQuery
			if len(gr) != len(wr) {
				t.Errorf("SearchR18Params.ToURL().RawQuery = %v, len = %v, want.RawQuery = %v, len = %v", gr, len(gr), wr, len(wr))
				return
			}

			gq, wq := got.Query(), tt.want.Query()
			if !reflect.DeepEqual(gq, wq) {
				t.Errorf("SearchR18Params.ToURL().Query %v, want.Query %v", gq, wq)
				return
			}
		})
	}
}

func TestSearchR18Params_NocGenres(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchR18Params
		want   []NocGenre
	}{
		{"default", &SearchR18Params{}, nil},
		{"noc genres", &SearchR18Params{nocGenres: []NocGenre{NocGenreNocturne}}, []NocGenre{NocGenreNocturne}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.NocGenres(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchR18Params.NocGenres() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchR18Params_AddNocGenres(t *testing.T) {
	type args struct {
		genres []NocGenre
	}
	tests := []struct {
		name   string
		params *SearchR18Params
		args   args
		want   []NocGenre
	}{
		{"genre + []", &SearchR18Params{nocGenres: []NocGenre{NocGenreMoonlightWomen}}, args{[]NocGenre{}}, []NocGenre{NocGenreMoonlightWomen}},
		{"[] + genre", &SearchR18Params{nocGenres: []NocGenre{}}, args{[]NocGenre{NocGenreNocturne}}, []NocGenre{NocGenreNocturne}},
		{"genres + genres",
			&SearchR18Params{nocGenres: []NocGenre{NocGenreNocturne}}, args{[]NocGenre{NocGenreMidnight}},
			[]NocGenre{NocGenreNocturne, NocGenreMidnight}},
		{"merge genres",
			&SearchR18Params{nocGenres: []NocGenre{NocGenreNocturne, NocGenreMidnight}}, args{[]NocGenre{NocGenreMidnight}},
			[]NocGenre{NocGenreNocturne, NocGenreMidnight}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddNocGenres(tt.args.genres)
			result := tt.params.NocGenres()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddNocGenres(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.genres, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchR18Params_ClearNocGenres(t *testing.T) {
	params := &SearchR18Params{nocGenres: []NocGenre{NocGenreNocturne}}
	params.ClearNocGenres()
	if params.NocGenres() != nil {
		t.Errorf("SearchParams.ClearNocGenres() should change NotNocGenres be nil, but %v", params.NocGenres())
	}
}

func TestSearchR18Params_queryFromNocGenre(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchR18Params
		want   url.Values
	}{
		{"no NocGenres, no query", &SearchR18Params{}, makeValues([][2]string{})},
		{`NocGenres:[MidNight], nocgenre=4`, &SearchR18Params{nocGenres: []NocGenre{NocGenreMidnight}}, makeValues([][2]string{{"nocgenre", "4"}})},
		{`NocGenres:[Nocturne, MidNight], nocgenre=1-4`,
			&SearchR18Params{nocGenres: []NocGenre{NocGenreNocturne, NocGenreMidnight}}, makeValues([][2]string{{"nocgenre", "1-4"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromNocGenre(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchR18Params.queryFromNocGenre() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchR18Params_NotNocGenres(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchR18Params
		want   []NocGenre
	}{
		{"default", &SearchR18Params{}, nil},
		{"not noc genres", &SearchR18Params{notNocGenres: []NocGenre{NocGenreNocturne}}, []NocGenre{NocGenreNocturne}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.NotNocGenres(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchR18Params.NotNocGenres() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchR18Params_AddNotNocGenres(t *testing.T) {
	type args struct {
		genres []NocGenre
	}
	tests := []struct {
		name   string
		params *SearchR18Params
		args   args
		want   []NocGenre
	}{
		{"genre + []", &SearchR18Params{notNocGenres: []NocGenre{NocGenreMoonlightWomen}}, args{[]NocGenre{}}, []NocGenre{NocGenreMoonlightWomen}},
		{"[] + genre", &SearchR18Params{notNocGenres: []NocGenre{}}, args{[]NocGenre{NocGenreNocturne}}, []NocGenre{NocGenreNocturne}},
		{"genres + genres",
			&SearchR18Params{notNocGenres: []NocGenre{NocGenreNocturne}}, args{[]NocGenre{NocGenreMidnight}},
			[]NocGenre{NocGenreNocturne, NocGenreMidnight}},
		{"merge genres",
			&SearchR18Params{notNocGenres: []NocGenre{NocGenreNocturne, NocGenreMidnight}}, args{[]NocGenre{NocGenreMidnight}},
			[]NocGenre{NocGenreNocturne, NocGenreMidnight}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddNotNocGenres(tt.args.genres)
			result := tt.params.NotNocGenres()
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("SearchParams.AddNotNocGenres(%v) result %v(l:%d,c:%d), want %v(l:%d,c:%d)",
					tt.args.genres, result, len(result), cap(result), tt.want, len(tt.want), cap(tt.want))
			}
		})
	}
}

func TestSearchR18Params_ClearNotNocGenres(t *testing.T) {
	params := &SearchR18Params{notNocGenres: []NocGenre{NocGenreNocturne}}
	params.ClearNotNocGenres()
	if params.NotNocGenres() != nil {
		t.Errorf("SearchParams.ClearNotNocGenres() should change NotNocGenres be nil, but %v", params.NotNocGenres())
	}
}

func TestSearchR18Params_queryFromNotNocGenre(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchR18Params
		want   url.Values
	}{
		{"no notNocGenres, no query", &SearchR18Params{}, makeValues([][2]string{})},
		{`notNocGenres:[MidNight], notnocgenre=4`, &SearchR18Params{notNocGenres: []NocGenre{NocGenreMidnight}}, makeValues([][2]string{{"notnocgenre", "4"}})},
		{`notNocGenres:[Nocturne, MidNight], notnocgenre=1-4`,
			&SearchR18Params{notNocGenres: []NocGenre{NocGenreNocturne, NocGenreMidnight}}, makeValues([][2]string{{"notnocgenre", "1-4"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromNotNocGenre(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchR18Params.queryFromNotNocGenre() = %v, want %v", got, tt.want)
			}
		})
	}
}
