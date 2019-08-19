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

func TestSearchR18Params_endPointURL(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchR18Params
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.endPointURL(); got != tt.want {
				t.Errorf("SearchR18Params.endPointURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchR18Params_toQueryFuncs(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchR18Params
		want   []func() url.Values
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.toQueryFuncs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchR18Params.toQueryFuncs() = %v, want %v", got, tt.want)
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
	// TODO: Add test cases.
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
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddNocGenres(tt.args.genres)
		})
	}
}

func TestSearchR18Params_ClearNocGenres(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchR18Params
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.ClearNocGenres()
		})
	}
}

func TestSearchR18Params_queryFromNocGenre(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchR18Params
		want   url.Values
	}{
	// TODO: Add test cases.
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
	// TODO: Add test cases.
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
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.AddNotNocGenres(tt.args.genres)
		})
	}
}

func TestSearchR18Params_ClearNotNocGenres(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchR18Params
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.ClearNotNocGenres()
		})
	}
}

func TestSearchR18Params_queryFromNotNocGenre(t *testing.T) {
	tests := []struct {
		name   string
		params *SearchR18Params
		want   url.Values
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.queryFromNotNocGenre(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchR18Params.queryFromNotNocGenre() = %v, want %v", got, tt.want)
			}
		})
	}
}
